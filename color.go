package main

import (
	hj "github.com/nikolaydubina/htmljson"
	"strings"
)

const (
	jqvSArrBrackets string = "color: indigo; font-weight: 600;"    // Array brackets
	jqvSMapBrackets        = "color: salmon; font-weight: 600;"    // Map brackets
	jqvSNum                = "color: fireBrick; font-weight: 400;" // int & float
	jqvSBoolTrue           = "color: green;"                       // Bool true
	jqvSBoolFalse          = "color: red;"                         // Bool false
	jqvSMapKey             = "color: teal; font-weight: 600;"      // Map key
	jqvSBlack              = "color: black"                        // Comma & Colon
	jvqSNull               = "color: indianRed;"                   // nil
	jqvSString             = "color: black; font-weight: 600;"     // Strings
)

func jvNullHTML(k string) string {
	return `<span style="` + jvqSNull + `">` + jvSpacer() + `null</span>`
}

func jvBoolHTML(k string, v bool) string {
	x := "false"
	if v {
		x = "true"
		return `<span style="` + jqvSBoolTrue + `">` + jvSpacer() + x + `</span>"`
	}
	return `<span style="` + jqvSBoolFalse + `">` + jvSpacer() + x + `</span>"`
}

func jvStringHTML(k string, v string) string {
	return `<span style="` + jqvSString + `">` + jvSpacer() + `"` + v + `"</span>`
}
func jvNumberHTML(k string, v float64, s string) string {
	return `<span style="` + jqvSNum + `">` + jvSpacer() + s + `</span>`
}

var jvArrayHTML = hj.ArrayMarshaler{
	OpenBracket:  `<span style="` + jqvSArrBrackets + `">[</span>`,
	CloseBracket: `<span style="` + jqvSArrBrackets + `">]</span>`,
	Comma:        `<span style="color: black">,</span>`,
}

var jvMapHTML = hj.MapMarshaler{
	OpenBracket:  `<span style="` + jqvSMapBrackets + `">{</span>`,
	CloseBracket: `<span style="` + jqvSMapBrackets + `">}</span>`,
	Comma:        `<span style="` + jqvSBlack + `">,</span>`,
	Colon:        `<span style="` + jqvSBlack + `">:</span>`,
	Key: func(key string, v string) string {
		return `<span style="` + jqvSMapKey + `">"` + v + `"</span>`
	},
}

type jvRowHTML struct {
	Padding int
}

func (s jvRowHTML) Marshal(v string, depth int) string {
	indent := strings.Repeat("&nbsp;", depth*2)
	return `<div style"margin-bottom: 0px;  margin-top:0px;"><span>` + indent + `</span>` + v + `</div>`
}

var jvMarsh = hj.Marshaler{
	Null: jvNullHTML, Bool: jvBoolHTML, String: jvStringHTML,
	Number: jvNumberHTML, Array: jvArrayHTML, Map: jvMapHTML,
	Row: jvRowHTML{Padding: 1}.Marshal,
}

func jvSpacer() string {
	return strings.Repeat("&nbsp;", 2)
}
