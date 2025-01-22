package main

import (
	"context"
	"encoding/json"
	"flag"
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

var (
	loadfileDialog *qt.QFileDialog
	input          *qt.QPlainTextEdit
	filter         *qt.QLineEdit
	output         *qt.QTextEdit
	colorize       bool
)

func main() {

	flag.BoolVar(&colorize, "colors", false, "Colorize the JSON")
	flag.Parse()

	app := qt.NewQApplication(flag.Args())

	if len(flag.Args()) > 0 {
		filterValue = flag.Args()[0]
	}
	if len(flag.Args()) > 1 {
		b, err := os.ReadFile(flag.Args()[1])
		if err != nil {
			log.Fatal("failed to read " + flag.Args()[1] + " : " + err.Error())
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
		}
	})

	loadfileButton := qt.NewQPushButton5("Load", window.QWidget)
	loadfileButton.SetMaximumWidth(50)
	loadfileButton.OnClicked1(func(_ bool) {
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
	filterLabel := qt.NewQLabel3("Filter:")
	filterLabel.SetMinimumWidth(50)
	filterLabel.SetMaximumWidth(50) // Same as the button
	filterLabel.SetAlignment(qt.AlignRight)
	filterLayout := qt.NewQHBoxLayout2()

	output = qt.NewQTextEdit(window.QWidget)
	output.SetSizeAdjustPolicy(qt.QAbstractScrollArea__AdjustToContents)
	output.SetMinimumHeight(300)

	if colorize {
		/* Static background color for the output widget
		in case the OS theme is dark
		*/
		app.SetStyleSheet("QTextEdit {background-color: rgb(255, 255, 255) }")
	}

	widget := qt.NewQWidget(window.QWidget)

	filterWidget := qt.NewQWidget(widget)
	filterWidget.SetLayout(filterLayout.Layout())
	filterWidget.Layout().AddWidget(filterLabel.QWidget)
	filterWidget.Layout().AddWidget(filter.QWidget)

	widget.SetLayout(qt.NewQVBoxLayout2().QLayout)
	widget.Layout().AddWidget(inputSection)
	widget.Layout().AddWidget(filterWidget)
	widget.Layout().AddWidget(output.QWidget)

	window.Show()
	window.SetCentralWidget(widget)

	refresh()

	qt.QApplication_Exec()
}

func refresh() {
	inputValue = input.ToPlainText()
	if colorize {
		output.SetHtml(runJQ(context.Background(), inputValue, filterValue))
	} else {
		output.SetText(runJQ(context.Background(), inputValue, filterValue))
	}
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
	var anyResults []any

	if colorize {
		anyResults = make([]any, 0)
	}

	for {
		v, exists := iter.Next()
		if !exists {
			break
		}

		if err, ok := v.(error); ok {
			return err.Error()
		}
		var s []byte

		if colorize {
			anyResults = append(anyResults, v)
		} else {
			s, _ = json.MarshalIndent(v, "", "  ")
			results = append(results, string(s))
		}

	}

	if colorize {
		html := string(jvMarsh.Marshal(anyResults[0]))
		return html
	} else {
		txt := strings.Join(results, "\n")
		return txt
	}
}
