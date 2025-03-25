package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

func main() {
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	var page Page
	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&page)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("<!DOCTYPE html>")
	fmt.Print("<html>")
	fmt.Print("<head>")
	fmt.Print("<meta charset=\"UTF-8\">")
	fmt.Print("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0, maximum-scale=1, user-scalable=0\">")
	fmt.Printf("<title>%s</title>", page.Title)
	fmt.Printf("<meta name=\"description\" content=\"%s\">", page.Description)
	fmt.Printf("<meta name=\"author\" content=\"%s\">", page.Author)
	fmt.Printf("<link rel=\"icon\" href=\"%s\">", page.Icon)
	fmt.Print("<script src=\"{bevyframe}/bridge.js\"></script>")
	fmt.Print("<script src=\"{bevyframe}/widgets.js\"></script>")
	fmt.Print("<style src=\"{bevyframe}/style.css\"></style>")
	fmt.Printf("<meta name=\"og:title\" content=\"%s\">", page.OpenGraph.Title)
	fmt.Printf("<meta name=\"og:description\" content=\"%s\">", page.OpenGraph.Description)
	fmt.Printf("<meta name=\"og:image\" content=\"%s\">", page.OpenGraph.Image)
	fmt.Printf("<meta name=\"og:url\" content=\"%s\">", page.OpenGraph.URL)
	fmt.Printf("<meta name=\"og:type\" content=\"%s\">", page.OpenGraph.Type)
	fmt.Print("</head>")
	fmt.Printf("<body class=\"body_%s\">", page.Color)
	fmt.Print("<nav class=\"Navbar\" id=\"navbar\">")
	for _, item := range page.Navbar.Items {
		fmt.Printf("<a href=\"%s\" class=\"%s\">", item.Link, item.Status)
		fmt.Print("<button>")
		fmt.Printf("<span class=\"material-symbols-rounded\" alt=\"%s\">%s</span>", item.Alt, item.Icon)
		fmt.Print("</button>")
		fmt.Print("</a>")
	}
	fmt.Print("</nav>")
	var rootStyle string
	if page.Root.LeftMargin != "" {
		rootStyle += fmt.Sprintf("margin-left:%s;", page.Root.LeftMargin)
	}
	if page.Root.TopMargin != "" {
		rootStyle += fmt.Sprintf("margin-top:%s;", page.Root.TopMargin)
	}
	if page.Root.RightMargin != "" {
		rootStyle += fmt.Sprintf("margin-right:%s;", page.Root.RightMargin)
	}
	if page.Root.BottomMargin != "" {
		rootStyle += fmt.Sprintf("margin-bottom:%s;", page.Root.BottomMargin)
	}
	if page.Root.TextAlign != "" {
		rootStyle += fmt.Sprintf("text-align:%s;", page.Root.TextAlign)
	}
	if page.Root.FontSize != "" {
		rootStyle += fmt.Sprintf("font-size:%s;", page.Root.FontSize)
	}
	if page.Root.VerticalAlign != "" {
		rootStyle += fmt.Sprintf("vertical-align:%s;", page.Root.VerticalAlign)
	}
	rootStyle = strings.ReplaceAll(rootStyle, "\"", "\\\"")
	loginRequired := ""
	if page.Root.LoginRequired == "true" {
		loginRequired = "login-required"
	}
	fmt.Printf("<div id=\"root\" style=\""+rootStyle+"\" %s>", loginRequired)
	fmt.Print(renderWidgets(page.Root.Content))
	fmt.Print("</div>")
	fmt.Print("</body>")
	fmt.Print("</html>")
	fmt.Println()
}
