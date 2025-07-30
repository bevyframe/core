package main

import "encoding/xml"

type Page struct {
	XMLName       xml.Name  `xml:"Page"`
	Title         string    `xml:"title,attr"`
	Description   string    `xml:"description,attr"`
	Color         string    `xml:"color,attr"`
	Icon          string    `xml:"icon,attr"`
	Author        string    `xml:"author,attr"`
	OpenGraph     OpenGraph `xml:"OpenGraph"`
	Navbar        Navbar    `xml:"Navbar"`
	Root          Root      `xml:"Root"`
	LoginRequired string    `xml:"loginRequired,attr"`
}

type OpenGraph struct {
	XMLName     xml.Name `xml:"OpenGraph"`
	Title       string   `xml:"title,attr"`
	Type        string   `xml:"type,attr"`
	Image       string   `xml:"image,attr"`
	URL         string   `xml:"url,attr"`
	Description string   `xml:"description,attr"`
}

type Navbar struct {
	XMLName xml.Name  `xml:"Navbar"`
	Items   []NavItem `xml:"NavItem"`
}

type NavItem struct {
	XMLName xml.Name `xml:"NavItem"`
	Icon    string   `xml:"icon,attr"`
	Link    string   `xml:"link,attr"`
	Alt     string   `xml:"alt,attr"`
	Status  string   `xml:"status,attr"`
}

type Root struct {
	XMLName       xml.Name `xml:"Root"`
	LeftMargin    string   `xml:"leftMargin,attr"`
	TopMargin     string   `xml:"topMargin,attr"`
	RightMargin   string   `xml:"rightMargin,attr"`
	BottomMargin  string   `xml:"bottomMargin,attr"`
	TextAlign     string   `xml:"textAlign,attr"`
	FontSize      string   `xml:"fontSize,attr"`
	VerticalAlign string   `xml:"verticalAlign,attr"`
	LoginRequired string   `xml:"loginRequired,attr"`
	Content       string   `xml:",innerxml"`
}
