package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/itchyny/gojq"
	qt "github.com/mappu/miqt/qt6"
)

var filterValue = "."
var inputValue = `[{
  "fruit": "mango"
}, {
  "fruit": "banana"
}]`

var loadfileDialog *qt.QFileDialog
var input *qt.QPlainTextEdit
var filter *qt.QLineEdit
var output *qt.QPlainTextEdit

func main() {
	qt.NewQApplication(os.Args)

	if len(os.Args) > 1 {
		filterValue = os.Args[1]
	}
	if len(os.Args) > 2 {
		b, err := os.ReadFile(os.Args[2])
		if err != nil {
			log.Fatal("failed to read " + os.Args[2] + " : " + err.Error())
		}
		inputValue = string(b)
	} else {
		if m, _ := os.Stdin.Stat(); m.Mode()&os.ModeCharDevice != os.ModeCharDevice {
			b, err := io.ReadAll(os.Stdin)
			if err == nil {
				inputValue = string(b)
			}
		}
	}

	window := qt.NewQMainWindow2()
	window.SetMinimumSize2(400, 500)
	window.SetWindowTitle("jqview")

	loadfileDialog = qt.NewQFileDialog4(window.QWidget, "Select a JSON file")
	loadfileDialog.OnFileSelected(func(filepath string) {
		b, err := os.ReadFile(filepath)
		if err != nil {
			log.Print("failed to read " + filepath + " : " + err.Error())
		} else {
			input.SetPlainText(string(b))
			refresh()
		}
	})

	loadfileButton := qt.NewQPushButton5("Load", window.QWidget)
	loadfileButton.SetMaximumWidth(50)
	loadfileButton.OnClicked(func() {
		loadfileDialog.Open()
	})

	input = qt.NewQPlainTextEdit(window.QWidget)
	input.SetPlaceholderText("JSON input")
	input.SetPlainText(inputValue)
	input.OnTextChanged(refresh)

	inputSection := qt.NewQWidget(window.QWidget)
	inputSection.SetLayout(qt.NewQHBoxLayout2().QLayout)
	inputSection.SetMaximumHeight(150)
	inputSection.Layout().AddWidget(loadfileButton.QWidget)
	inputSection.Layout().AddWidget(input.QWidget)

	filter = qt.NewQLineEdit(window.QWidget)
	filter.SetPlaceholderText("jq filter")
	filter.SetText(filterValue)
	filter.OnTextChanged(func(value string) {
		filterValue = value
		refresh()
	})

	output = qt.NewQPlainTextEdit(window.QWidget)
	output.SetSizeAdjustPolicy(qt.QAbstractScrollArea__AdjustToContents)
	output.SetMinimumHeight(300)

	widget := qt.NewQWidget(window.QWidget)
	widget.SetLayout(qt.NewQVBoxLayout2().QLayout)
	widget.Layout().AddWidget(inputSection)
	widget.Layout().AddWidget(filter.QWidget)
	widget.Layout().AddWidget(output.QWidget)

	window.Show()
	window.SetCentralWidget(widget)

	refresh()

	qt.QApplication_Exec()
}

func refresh() {
	inputValue = input.ToPlainText()
	output.SetPlainText(runJQ(context.Background(), inputValue, filterValue))
}

func runJQ(
	ctx context.Context,
	input string,
	filter string,
) string {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	var object interface{}
	err := json.Unmarshal([]byte(input), &object)
	if err != nil {
		return err.Error()
	}

	query, err := gojq.Parse(filter)
	if err != nil {
		return err.Error()
	}

	iter := query.RunWithContext(ctx, object)

	var results []string
	for {
		v, exists := iter.Next()
		if !exists {
			break
		}

		if err, ok := v.(error); ok {
			return err.Error()
		}

		s, _ := json.MarshalIndent(v, "", "  ")
		results = append(results, string(s))
	}

	return strings.Join(results, "\n")
}
