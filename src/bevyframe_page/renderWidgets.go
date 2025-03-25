package main

import (
	"encoding/xml"
	"fmt"
	"strings"
)

func bevyToHTML(element string) string {
	switch element {
	case "Title":
		return "h1"
	case "Box":
		return "div class=\"the_box\""
	case "Container":
		return "div"
	case "Line":
		return "p"
	case "Textbox":
		return "input class=\"textbox\""
	case "Button":
		return "button class=\"button\""
	case "SmallButton":
		return "button class=\"button small\""
	case "MiniButton":
		return "button class=\"button mini\""
	case "Form":
		return "form"
	default:
		return "?????"
	}
}

func renderWidgets(content string) string {
	var result strings.Builder
	decoder := xml.NewDecoder(strings.NewReader(content))

	for {
		token, err := decoder.Token()
		if err != nil || token == nil {
			break
		}

		switch t := token.(type) {
		case xml.StartElement:
			result.WriteString("<")
			result.WriteString(bevyToHTML(t.Name.Local))
			style := ""

			for _, attr := range t.Attr {
				if attr.Name.Local != "the_box" && attr.Name.Local != "style" {
					name := strings.ToLower(attr.Name.Local)
					switch name {
					case "leftmargin":
						style += fmt.Sprintf("margin-left: %s;", attr.Value)
					case "topmargin":
						style += fmt.Sprintf("margin-top: %s;", attr.Value)
					case "rightmargin":
						style += fmt.Sprintf("margin-right: %s;", attr.Value)
					case "bottommargin":
						style += fmt.Sprintf("margin-bottom: %s;", attr.Value)
					case "textalign":
						style += fmt.Sprintf("text-align: %s;", attr.Value)
					case "fontsize":
						style += fmt.Sprintf("font-size: %s;", attr.Value)
					case "verticalalign":
						style += fmt.Sprintf("vertical-align: %s;", attr.Value)
					case "width":
						style += fmt.Sprintf("width: %s;", attr.Value)
					case "height":
						style += fmt.Sprintf("height: %s;", attr.Value)
					default:
						result.WriteString(fmt.Sprintf(" %s=\"%s\"", name, attr.Value))
					}
				}
			}
			result.WriteString(fmt.Sprintf(" style=\"%s\"", style))
			result.WriteString(">")

		case xml.EndElement:
			result.WriteString("</")
			result.WriteString(strings.SplitN(bevyToHTML(t.Name.Local), " ", 2)[0])
			result.WriteString(">")

		case xml.CharData:
			result.Write(t)

		case xml.Comment:
			result.WriteString("<!--")
			result.Write(t)
			result.WriteString("-->")
		}
	}

	return result.String()
}
