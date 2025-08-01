package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func findFilePath(filePath string) (string, error) {
	files, err := os.ReadDir(filePath)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if file.Name() == "index.html" {
			return fmt.Sprintf("%s/index.html", filePath), nil
		} else if file.Name() == "index.bevy" {
			return fmt.Sprintf("%s/index.bevy", filePath), nil
		} else if file.Name() == "index" {
			return fmt.Sprintf("%s/index", filePath), nil
		} else if file.Name() == "__init__.py" {
			return fmt.Sprintf("%s/__init__.py", filePath), nil
		} else if file.Name() == "index.js" {
			return fmt.Sprintf("%s/index.js", filePath), nil
		}
	}
	return "", fmt.Errorf("file not found")
}

func (context Context) execute(filePath string, reqTime string, bd []byte) Response {
	username, network := "Guest", context.app.manifest.Accounts.DefaultNetwork
	if strings.Contains(context.email, "@") {
		username = strings.Split(context.email, "@")[0]
		network = strings.Split(context.email, "@")[1]
	}
	r := ""
	r += fmt.Sprintf("Package: %s\n", context.app.name)
	r += fmt.Sprintf("Cred.Email: %s@%s\n", username, network)
	r += fmt.Sprintf("Cred.Username: %s\n", username)
	r += fmt.Sprintf("Cred.Network: %s\n", network)
	r += fmt.Sprintf("Cred.Token: %s\n", context.token)
	r += fmt.Sprintf("Path: %s\n", context.path)
	r += fmt.Sprintf("IP: %s\n", context.ip)
	r += fmt.Sprintf("Method: %s\n", context.method)
	r += fmt.Sprintf("Permissions: %s\n", strings.Join(context.app.manifest.Accounts.Permissions, ","))
	r += fmt.Sprintf("LoginView: %s\n", context.app.manifest.App.LoginView)
	if _, ok := context.headers["Date"]; !ok {
		r += fmt.Sprintf("Header.Date: %s\n", reqTime)
	}
	for key, value := range context.headers {
		r += fmt.Sprintf("Header.%s: %s\n", key, value)
	}
	for key, value := range context.query {
		r += fmt.Sprintf("Query.%s: %s\n", key, value)
	}
	b := []byte(r)
	if bd != nil {
		b = bytes.Join([][]byte{b, []byte("\n\n"), bd, []byte("\n")}, nil)
	}
	filetype, err := exec.Command("file", filePath).Output()
	if err != nil {
		return Response{
			statusCode: 500,
			headers: map[string]string{
				"Content-Type": "text/html",
			},
			body: "<h1>Internal Server Error</h1>`file` command failed",
		}
	}
	cmd := exec.Command("cat")
	if strings.HasSuffix(filePath, ".bevy") {
		cmdP := exec.Command(FindInstallation()+"/bin/bevyframe_page", filePath)
		out, err := cmdP.Output()
		if err != nil {
			return Response{
				statusCode: 500,
				headers: map[string]string{
					"Content-Type": "text/html",
				},
				body: "<h1>Internal Server Error</h1>bevyframe_page failed",
			}
		}
		cmd = exec.Command(os.Getenv("BEVYFRAME_HTML_SDK"), string(out))
	} else if strings.Contains(string(filetype), "Python script") {
		cmd = exec.Command(os.Getenv("BEVYFRAME_PYTHON_SDK"), filePath)
	} else if strings.Contains(string(filetype), "HTML document text") {
		cmd = exec.Command(os.Getenv("BEVYFRAME_HTML_SDK"), filePath)
	} else if strings.Contains(string(filetype), "Java source text") {
		sdkInfo := strings.Split(os.Getenv("BEVYFRAME_JAVA_SDK"), " ")
		cmd = exec.Command(sdkInfo[0], "-classpath", sdkInfo[1], filePath)
	} else {
		extS := strings.Split(filePath, ".")
		ext := extS[len(extS)-1]
		if ext == "js" {
			r = "const stdin = \"" + strings.ReplaceAll(r, "\"", "\\\"") + "\""
			r = strings.ReplaceAll(r, "\n", "\\n") + "\n"
			r += CreateBridgeScript() + "\n"
			script, err := os.ReadFile(filePath)
			if err != nil {
				return Response{
					statusCode: 500,
					headers:    map[string]string{},
					body:       "<h1>Internal Server Error</h1>Failed to read script file",
				}
			}
			renderJSb, err := os.ReadFile(FindInstallation() + "/scripts/renderJS.js")
			if err != nil {
				return Response{
					statusCode: 500,
					headers:    map[string]string{},
					body:       "<h1>Internal Server Error</h1>Failed to render JS",
				}
			}
			renderJS := string(renderJSb)
			renderJS = strings.ReplaceAll(renderJS, "/* PAGE SCRIPT HERE */", string(script))
			r += renderJS
			return Response{
				statusCode: 200,
				headers: map[string]string{
					"Content-Type": "text/html",
				},
				body: context.renderJS(r),
			}
		} else {
			cmd = exec.Command(filePath)
		}
	}
	cmd.Stdin = bytes.NewReader(b)
	out, err := cmd.Output()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return Response{
				statusCode: 404,
				headers: map[string]string{
					"Content-Type": "text/html",
				},
				body: "<h1>Not Found</h1>The requested file does not exist.",
			}
		}
		return Response{
			statusCode: 404,
			headers: map[string]string{
				"Content-Type": "text/html",
			},
			body: "<h1>Not Found</h1>Script execution failed.",
		}
	}
	out = bytes.ReplaceAll(out, []byte("\r\n"), []byte("\n"))
	if strings.Contains(string(out), "\n\n<!DOCTYPE html>") {
		out = []byte(strings.ReplaceAll(string(out), "src=\"{bevyframe}/style.css\">", ">"+context.app.style))
		out = []byte(strings.ReplaceAll(
			string(out),
			"src=\"{bevyframe}/widgets.js\">", "src=\"/.well-known/bevyframe/widgets.js\">",
		))
		out = []byte(strings.ReplaceAll(
			string(out),
			"src=\"{bevyframe}/style.css\">", ">"+context.app.style,
		))
		out = []byte(strings.ReplaceAll(
			string(out),
			"src=\"{bevyframe}/bridge.js\">", "src=\"/.well-known/bevyframe/bridge.js\"></script><script>"+
				CreateBridgeScript(),
		))
		out = []byte(strings.ReplaceAll(
			string(out),
			"src=\"{bevyframe}/renderWidget.js\">", "src=\"/.well-known/bevyframe/renderWidget.js\"></script>",
		))
	}

	lines := bytes.Split(out, []byte("\n"))
	if len(lines) == 0 {
		return Response{
			statusCode: 500,
			headers:    map[string]string{},
			body:       "<h1>Internal Server Error</h1>empty output",
		}
	}
	parts := bytes.SplitN(lines[0], []byte(" "), 3)
	if len(parts) < 2 {
		return Response{
			statusCode: 500,
			headers:    map[string]string{},
			body:       fmt.Sprint("Internal Server Error: ", string(out)),
		}
	}
	statusCode, err := strconv.Atoi(string(parts[1]))
	if err != nil {
		return Response{
			statusCode: 500,
			headers: map[string]string{
				"Content-Type": "text/html",
			},
			body: fmt.Sprint("<h1>Internal Server Error</h1> Executable returned unrecognized API"),
		}
	}
	headers := make(map[string]string)

	for _, line := range lines[1:] {
		if len(strings.ReplaceAll(string(line), " ", "")) == 0 {
			break
		}
		t := bytes.SplitN(line, []byte(": "), 2)
		if len(t) < 2 {
			continue
		}
		headers[string(t[0])] = string(t[1])
	}

	var body []byte
	if len(lines) > len(headers) {
		body = bytes.Join(lines[len(headers)+1:], []byte("\n"))
		body = bytes.Trim(body, "\r\n")
		body = bytes.Trim(body, "\r")
		body = bytes.Trim(body, "\n")
	} else {
		body = []byte("")
	}

	headers["Server"] = "BevyFrame"

	if _, ok := headers["Content-Length"]; !ok {
		delete(headers, "Content-Length")
	}

	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "text/plain"
	}

	if _, ok := headers["Date"]; !ok {
		headers["Date"] = reqTime
	}

	if _, ok := headers["Connection"]; !ok {
		headers["Connection"] = "close"
	}

	if _, ok := headers["Content-Encoding"]; !ok {
		headers["Content-Encoding"] = "identity"
	}

	if _, ok := headers["Transfer-Encoding"]; !ok {
		headers["Transfer-Encoding"] = "identity"
	}

	if _, ok := headers["Vary"]; !ok {
		headers["Vary"] = "Accept-Encoding"
	}

	if _, ok := headers["X-Content-Type-Options"]; !ok {
		headers["X-Content-Type-Options"] = "nosniff"
	}

	if _, ok := headers["X-Frame-Options"]; !ok {
		headers["X-Frame-Options"] = "DENY"
	}

	if _, ok := headers["X-XSS-Protection"]; !ok {
		headers["X-XSS-Protection"] = "1; mode=block"
	}

	if _, ok := headers["Strict-Transport-Security"]; !ok {
		headers["Strict-Transport-Security"] = "max-age=31536000"
	}

	if _, ok := headers["Referrer-Policy"]; !ok {
		headers["Referrer-Policy"] = "no-referrer"
	}

	if _, ok := headers["Feature-Policy"]; !ok {
		headers["Feature-Policy"] = "accelerometer 'none'; camera 'none'; geolocation 'none'; gyroscope 'none'; magnetometer 'none'; microphone 'none'; payment 'none'; usb 'none'"
	}

	if _, ok := headers["Cache-Control"]; !ok {
		headers["Cache-Control"] = "no-cache, no-store, must-revalidate"
	}

	if _, ok := headers["Pragma"]; !ok {
		headers["Pragma"] = "no-cache"
	}

	if _, ok := headers["Expires"]; !ok {
		headers["Expires"] = "0"
	}

	email, ok1 := headers["Cred-Email"]
	token, ok2 := headers["Cred-Token"]
	if ok1 && ok2 {
		cookie, err := context.app.getSessionToken(email, token)
		if err == nil {
			headers["Set-Cookie"] = fmt.Sprintf("s=%s; Path=/", cookie)
		}
		delete(headers, "Cred-Email")
		delete(headers, "Cred-Token")
	}

	if headers["Content-Type"] == "application/bevyframe" {
		bodyStr := string(body)
		if strings.HasPrefix(bodyStr, "Response.Type: Page") {
			page := context.loadPage(bodyStr)
			body = []byte(page.renderPage())
			headers["Content-Type"] = "text/html"
		} else if strings.HasPrefix(bodyStr, "Response.Type: Redirect") {
			headers["Location"] = strings.Split(bodyStr, "\n")[1]
			statusCode = 303
			headers["Content-Type"] = "text/plain"
		} else if strings.HasPrefix(bodyStr, "Response.Type: Error") {
			statusCode = 500
			if !strings.HasPrefix(context.headers["User-Agent"], "Mozilla/") {
				body = []byte(bodyStr)
			} else {
				body = context.app.errorHandler(bodyStr)
			}
			headers["Content-Type"] = "text/html"
		} else {
			statusCode = 200
			headers["Content-Type"] = "text/plain"
		}
	}

	return Response{
		statusCode: statusCode,
		headers:    headers,
		body:       string(body),
	}
}
