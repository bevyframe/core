package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func (app Frame) errorHandler(traceback string) []byte {
	traceback = strings.Replace(traceback, "Response.Type: Error\n\n", "", 1)
	traceback = strings.ReplaceAll(traceback, "<", "&lt;")
	traceback = strings.ReplaceAll(traceback, ">", "&gt;")
	traceback = strings.ReplaceAll(traceback, "\n", "</br>")
	traceback = strings.ReplaceAll(traceback, " ", "&nbsp;")
	traceback = strings.ReplaceAll(traceback, "\"", "\\\"")
	prop, _ := json.Marshal(map[string]string{
		"style": "font-family: monospace; padding: 20px;",
		"class": "the_box",
	})
	var widgets string
	if strings.HasPrefix(traceback, "Traceback&nbsp;(most&nbsp;recent&nbsp;call&nbsp;last):</br>") {
		traceback = strings.TrimPrefix(traceback, "Traceback&nbsp;(most&nbsp;recent&nbsp;call&nbsp;last):</br>")
		pwd, _ := os.Getwd()
		traceback = strings.ReplaceAll(traceback, fmt.Sprintf("File&nbsp;\\\"%s/pages/", pwd), "Path&nbsp;\\\"")
		lines := strings.Split(traceback, "</br>")
		lastLine := lines[len(lines)-1]
		traceback = strings.TrimSuffix(traceback, "</br>"+lastLine)
		widgets = "[[\"h1\", {}, [\"Traceback (most recent call last):\"]],"
		widgets += fmt.Sprintf("[\"div\", %s, [\"%s\"]],", prop, traceback)
		widgets += fmt.Sprintf("[\"h3\", {}, [\"%s\"]]]", lastLine)
	} else {
		widgets = fmt.Sprintf("[[\"div\", %s, [\"%s\"]]]", prop, traceback)
	}
	p := Page{
		charset:     "UTF-8",
		viewport:    "width=device-width, initial-scale=1.0, maximum-scale=1, user-scalable=0",
		description: "",
		author:      "",
		icon:        "/favicon.ico",
		title:       "BevyFrame Error Handler",
		data:        map[string]interface{}{},
		themeColor:  "blank",
		style:       app.style,
		openGraph: OpenGraph{
			title:       "",
			description: "",
			image:       "",
			url:         "",
			type_:       "",
		},
		widgets: widgets,
		app:     app,
	}
	return []byte(p.renderPage())
}
