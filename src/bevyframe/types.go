package main

type Context struct {
	path    string
	app     Frame
	email   string
	token   string
	ip      string
	method  string
	headers map[string]string
	query   map[string]string
}

type Response struct {
	statusCode int
	headers    map[string]string
	body       string
}
