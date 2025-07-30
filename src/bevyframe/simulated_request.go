package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// RequestData represents the expected JSON structure from stdin
type RequestData struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type ResponseData struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

func SimulatedRequest() {
	jsonData, err := io.ReadAll(os.Stdin)
	if err != nil {
		response := ResponseData{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
			Body: "Error reading input: " + err.Error(),
		}
		responseJSON, _ := json.Marshal(response)
		fmt.Println(string(responseJSON))
		return
	}
	var reqData RequestData
	err = json.Unmarshal(jsonData, &reqData)
	if err != nil {
		response := ResponseData{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "plain/text",
			},
			Body: "Error parsing JSON: " + err.Error(),
		}
		responseJSON, _ := json.Marshal(response)
		fmt.Println(string(responseJSON))
		return
	}
	manifest, err := loadManifest()
	if err != nil {
		response := ResponseData{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
			Body: "Failed to load manifest: " + err.Error(),
		}
		responseJSON, _ := json.Marshal(response)
		fmt.Println(string(responseJSON))
		return
	}
	frame := newServer(*manifest)

	var cred = map[string]string{
		"email": fmt.Sprintf("Guest@%s", frame.manifest.Accounts.DefaultNetwork),
		"token": "",
	}

	if cookie, ok := reqData.Headers["Cookie"]; ok {
		cookieParts := strings.Split(cookie, ";")
		for _, part := range cookieParts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "s=") {
				sessionValue, _ := strings.CutPrefix(part, "s=")
				if sessionCred, err := frame.getSession(sessionValue); err == nil {
					cred = sessionCred
				}
				break
			}
		}
	}

	reqTime := time.Now().UTC().Format("01/02/2006 03:04:05 PM")

	path := reqData.Path
	query := make(map[string]string)
	if strings.Contains(reqData.Path, "?") {
		pathParts := strings.SplitN(reqData.Path, "?", 2)
		path = pathParts[0]
		queryParts := strings.Split(pathParts[1], "&")
		for _, part := range queryParts {
			if strings.Contains(part, "=") {
				kv := strings.SplitN(part, "=", 2)
				query[kv[0]] = kv[1]
			} else {
				query[part] = ""
			}
		}
	}

	context := Context{
		path:    path,
		app:     frame,
		email:   cred["email"],
		token:   cred["token"],
		ip:      reqData.Headers["X-Forwarded-For"],
		method:  reqData.Method,
		headers: reqData.Headers,
		query:   query,
	}

	pwd, err := os.Getwd()
	if err != nil {
		response := ResponseData{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
			Body: "Failed to get working directory: " + err.Error(),
		}
		responseJSON, _ := json.Marshal(response)
		fmt.Println(string(responseJSON))
		return
	}

	filePath := fmt.Sprintf("%s/pages%s", pwd, path)

	for key, value := range context.app.manifest.App.Routing {
		if key == path || key+"/" == path {
			filePath = fmt.Sprintf("%s/pages%s", pwd, value)
			break
		} else {
			variables, err := matchRouting(key, path)
			if err == nil {
				filePath = fmt.Sprintf("%s/pages%s", pwd, value)
				for k, v := range variables {
					context.query[k] = v
				}
				break
			}
		}
	}

	fileStat, err := os.Stat(filePath)
	if err == nil && fileStat.IsDir() {
		filePath, err = findFilePath(filePath)
		if err != nil {
			response := ResponseData{
				StatusCode: 404,
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
				Body: "File not found",
			}
			responseJSON, _ := json.Marshal(response)
			fmt.Println(string(responseJSON))
			return
		}
	}

	var bodyBytes []byte
	if reqData.Body != "" {
		bodyBytes = []byte(reqData.Body)
	}

	resp := context.execute(filePath, reqTime, bodyBytes)

	response := ResponseData{
		StatusCode: resp.statusCode,
		Headers:    resp.headers,
		Body:       resp.body,
	}

	responseJSON, _ := json.Marshal(response)
	fmt.Println(string(responseJSON))
}
