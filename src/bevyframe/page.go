package main

import (
	"fmt"
	"os"
	"strings"
)

type OpenGraph struct {
	title       string
	description string
	image       string
	url         string
	type_       string
}

type Page struct {
	charset     string
	viewport    string
	description string
	author      string
	icon        string
	title       string
	data        map[string]interface{}
	themeColor  string
	style       string
	openGraph   OpenGraph
	widgets     string
	app         Frame
}

func (context Context) loadPage(str string) Page {
	ret := Page{
		charset:     "UTF-8",
		viewport:    "width=device-width, initial-scale=1.0, maximum-scale=1, user-scalable=0",
		description: "BevyFrame Test App",
		author:      "",
		icon:        "/favicon.ico",
		title:       "Login - BevyFrame Test App",
		data:        map[string]interface{}{},
		themeColor:  "blue",
		style:       "",
		openGraph: OpenGraph{
			title:       "WebApp",
			description: "BevyFrame App",
			image:       "/Banner.png",
			url:         "",
			type_:       "website",
		},
		widgets: "",
		app:     context.app,
	}
	inHeaders := true
	for _, line := range strings.Split(str, "\n") {
		if inHeaders {
			if line == "" {
				inHeaders = false
			} else {
				parts := strings.SplitN(line, ":", 2)
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				switch key {
				case "Response.Charset":
					ret.charset = value
				case "Response.Viewport":
					ret.viewport = value
				case "Response.Description":
					ret.description = value
				case "Response.Author":
					ret.author = value
				case "Response.Icon":
					ret.icon = value
				case "Response.Title":
					ret.title = value
				case "Response.Data":
					ret.data = map[string]interface{}{}
				case "Response.ThemeColor":
					ret.themeColor = value
				case "Response.Style":
					ret.style = value
				case "Response.OpenGraph.title":
					ret.openGraph.title = value
				case "Response.OpenGraph.description":
					ret.openGraph.description = value
				case "Response.OpenGraph.image":
					ret.openGraph.image = value
				case "Response.OpenGraph.url":
					ret.openGraph.url = value
				case "Response.OpenGraph.type":
					ret.openGraph.type_ = value
				}
			}
		} else {
			ret.widgets += line
		}
	}
	return ret
}

func (p Page) createHTML(script string) string {
	filePath := FindInstallation() + "/scripts/renderWidget.js"
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	renderWidgetJS := string(content)
	renderWidgetJS = strings.ReplaceAll(renderWidgetJS, "`---body---`", p.widgets)
	ret := "<!DOCTYPE html>\n<html>\n\t<head>"
	ret += "\t\t<meta charset=\"" + p.charset + "\">\n"
	ret += "\t\t<meta name=\"viewport\" content=\"" + p.viewport + "\">\n"
	ret += "\t\t<meta name=\"description\" content=\"" + p.description + "\">\n"
	ret += "\t\t<meta name=\"author\" content=\"" + p.author + "\">\n"
	ret += "\t\t<link rel=\"manifest\" href=\"/.well-known/bevyframe/pwa.webmanifest\" />\n"
	ret += "\t\t<link rel=\"icon\" href=\"" + p.icon + "\">\n"
	ret += "\t\t<title>" + p.title + "</title>\n"
	ret += "\t\t<meta name=\"og:title\" content=\"" + p.openGraph.title + "\">\n"
	ret += "\t\t<meta name=\"og:description\" content=\"" + p.openGraph.description + "\">\n"
	ret += "\t\t<meta name=\"og:image\" content=\"" + p.openGraph.image + "\">\n"
	ret += "\t\t<meta name=\"og:url\" content=\"" + p.openGraph.url + "\">\n"
	ret += "\t\t<meta name=\"og:type\" content=\"" + p.openGraph.type_ + "\">\n"
	ret += "\t\t<script src=\"/.well-known/bevyframe/bridge.js\"></script>\n"
	ret += "\t\t<script src=\"/.well-known/bevyframe/buildContext.js\"></script>\n"
	ret += "\t\t<script src=\"/.well-known/bevyframe/widgets.js\"></script>\n"
	ret += "\t\t<script>\n"
	ret += renderWidgetJS + "\n"
	ret += script
	ret += "</script>\n"
	ret += "\t\t<style>" + p.app.style + "</style>\n"
	ret += "\t</head>\n"
	ret += "\t<body class=\"body_" + p.themeColor + "\" onload=\"buildDocument()\"></body></html>"
	return ret
}

func (p Page) renderPage() string {
	script := "if (typeof navigator.serviceWorker !== 'undefined') navigator.serviceWorker.register('sw.js');\n"
	script += CreateBridgeScript() + "\n\t\t"
	script += "const buildDocument = () => {renderAll()};"
	ret := p.createHTML(script)
	return ret
}

func (self Context) renderJS(script string) string {
	p := Page{
		charset:     "UTF-8",
		viewport:    "width=device-width, initial-scale=1.0, maximum-scale=1, user-scalable=0",
		description: self.app.manifest.Publishing.Description,
		author:      "",
		icon:        self.app.manifest.App.Icon,
		title:       self.app.manifest.App.Name,
		data:        map[string]interface{}{},
		themeColor:  "blank",
		style:       self.app.style,
		openGraph: OpenGraph{
			title:       self.app.manifest.App.Name,
			description: self.app.manifest.Publishing.Description,
			image:       self.app.manifest.App.Icon,
			url:         "",
			type_:       "pwa",
		},
		widgets: "[]",
		app:     self.app,
	}
	return p.createHTML(script)
}
