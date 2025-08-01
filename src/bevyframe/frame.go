package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Frame struct {
	manifest Manifest
	name     string
	secret   []byte
	style    string
}

type TheProtocolsProxy struct {
	Network  string `json:"network"`
	Endpoint string `json:"endpoint"`
	Body     string `json:"body"`
}

func newServer(manifest Manifest) Frame {
	var secretHex string
	secretHex = os.Getenv("SECRET")
	if secretHex == "" {
		secretBytes, err := os.ReadFile("./.secret")
		if err != nil {
			fmt.Println("Error reading ./.secret file:", err)
			os.Exit(1)
		}
		secretHex = string(secretBytes)
	}
	bytesObject, err := hex.DecodeString(secretHex)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	styleType := strings.Split(manifest.App.Style, ":")[0]
	styleName := strings.Split(manifest.App.Style, ":")[1]
	var style []byte
	if styleType == "python" {
		bevystylePy := strings.ReplaceAll(os.Getenv("BEVYFRAME_PYTHON_SDK"), "/bevyframe_py", "/bevystyle_py")
		styleCmd := exec.Command(bevystylePy, styleName)
		style, err = styleCmd.Output()
		if err != nil {
			fmt.Println("Error running style command:", err)
			os.Exit(1)
		}
	} else if styleType == "https" {
		resp, err := http.Get(fmt.Sprintf("https:%s", styleName))
		if err != nil {
			fmt.Println("Error fetching style from URL:", err)
			os.Exit(1)
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				fmt.Println("Error closing resp.Body:", err)
			}
		}(resp.Body)
		style, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading style from response body:", err)
			os.Exit(1)
		}
	}
	return Frame{
		manifest: manifest,
		name:     manifest.App.Package,
		secret:   bytesObject,
		style:    string(style),
	}
}

func (self Frame) runServer(debug bool) {
	http.HandleFunc("/.well-known/bevyframe/pwa.webmanifest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		out, err := self.processPWA()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, err = w.Write(out)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/.well-known/bevyframe/widgets.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		http.ServeFile(w, r, FindInstallation()+"/scripts/Widgets.js")
	})
	http.HandleFunc("/.well-known/bevyframe/bridge.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		http.ServeFile(w, r, FindInstallation()+"/scripts/bridge.js")
	})
	http.HandleFunc("/.well-known/bevyframe/renderWidget.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		http.ServeFile(w, r, FindInstallation()+"/scripts/renderWidget.js")
	})
	http.HandleFunc("/.well-known/bevyframe/buildContext.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		http.ServeFile(w, r, FindInstallation()+"/scripts/buildContext.js")
	})
	http.HandleFunc("/.well-known/theprotocols", func(w http.ResponseWriter, r *http.Request) {
		var data TheProtocolsProxy
		var rawData []byte
		rawData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(rawData, &data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			fmt.Println("Data2")
			return
		}
		url := fmt.Sprintf("https://%s/protocols/%s", data.Network, data.Endpoint)
		bodyReader := strings.NewReader(data.Body)
		resp, err := http.Post(url, "application/json", bodyReader)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(respBody)
	})
	http.HandleFunc("/sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		out, err := self.getServiceWorker()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, err = w.Write(out)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf(".%s", r.URL.Path))
	})
	http.HandleFunc("/node_modules/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf(".%s", r.URL.Path))
	})
	http.HandleFunc("/.well-known/bevyframe/proxy", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]interface{})
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return
		}
		function := data["func"].(string)
		args := data["args"].(string)
		path := data["path"].(string)
		s := data["cookie"].(string)
		reqTime := time.Now().UTC().Format("01/02/2006 03:04:05 PM")
		s, _ = strings.CutPrefix(s, "s=")
		s = strings.SplitN(s, ";", 1)[0]
		cred, err := self.getSession(s)
		id := r.RemoteAddr
		if err == nil || (cred["email"] != "" && strings.Split(cred["email"], "@")[0] != "Guest") {
			id = cred["email"]
		}
		fmt.Printf("func: %s [%s] ", id, reqTime)
		context := Context{
			path:    path,
			app:     self,
			email:   cred["email"],
			token:   cred["token"],
			ip:      r.RemoteAddr,
			method:  "FUNCTION",
			headers: headersToMap(r.Header),
			query:   map[string]string{},
		}
		resp, err := context.ProcessBridgeProxy(function, args, reqTime)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(resp))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s := ""
		for _, value := range r.Cookies() {
			if value.Name == "s" {
				s = value.Value
				break
			}
		}
		cred, err := self.getSession(s)
		id := ""
		if err == nil || strings.Split(cred["email"], "@")[0] != "Guest" {
			id = cred["email"]
		} else {
			id = r.RemoteAddr
			cred = map[string]string{
				"email": fmt.Sprintf("Guest@%s", self.manifest.Accounts.DefaultNetwork),
				"token": "",
			}
		}
		context := Context{
			path:    r.URL.Path,
			app:     self,
			email:   cred["email"],
			token:   cred["token"],
			ip:      r.RemoteAddr,
			method:  r.Method,
			headers: headersToMap(r.Header),
			query:   map[string]string{},
		}
		reqTime := time.Now().UTC().Format("01/02/2006 03:04:05 PM")
		fmt.Printf("(   ) %s [%s] %s %s   ", id, reqTime, context.method, context.path)
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		filePath := fmt.Sprintf("%s/pages%s", pwd, context.path)
		for key, value := range context.app.manifest.App.Routing {
			if key == context.path || key+"/" == context.path {
				filePath = fmt.Sprintf("%s/pages%s", pwd, value)
			} else {
				variables, err := matchRouting(key, context.path)
				if err == nil {
					filePath = fmt.Sprintf("%s/pages%s", pwd, value)
					for k, v := range variables {
						context.query[k] = v
					}
				}
			}
		}
		if r.URL.RawQuery != "" {
			queryParams := r.URL.Query()
			for key, values := range queryParams {
				if len(values) > 0 {
					context.query[key] = values[0]
				}
			}
		}
		fileStat, err1 := os.Stat(filePath)
		resp := Response{
			statusCode: 404,
			headers:    map[string]string{},
			body:       "NotFound",
		}
		if err1 == nil {
			if fileStat.IsDir() {
				filePath, err = findFilePath(filePath)
			}
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		resp = context.execute(filePath, reqTime, body)
		if _, ok := resp.headers["Digest"]; !ok {
			hash := sha256.Sum256([]byte(resp.body))
			digest := base64.StdEncoding.EncodeToString(hash[:])
			resp.headers["Digest"] = fmt.Sprintf("SHA-256=%s", digest)
		}
		for key, value := range resp.headers {
			w.Header().Set(key, value)
		}
		w.WriteHeader(resp.statusCode)
		fmt.Printf("\r(%d)\n", resp.statusCode)
		_, err = w.Write([]byte(resp.body))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
	fmt.Println("\nBevyFrame 0.6 ⍺")
	fmt.Printf(" * Serving BevyFrame app '%s'\n", self.manifest.App.Package)
	fmt.Print(" * Mode: ")
	if debug {
		fmt.Println("debug")
	} else {
		fmt.Println("production")
	}
	url := fmt.Sprintf("http://localhost:%d/", self.manifest.Development.Port)
	fmt.Println(" * Running on", url)
	fmt.Println()
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Print("\r  \nServer stopped.\n\n")
		os.Exit(0)
	}()
	err := http.ListenAndServe(fmt.Sprintf(":%d", self.manifest.Development.Port), nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}
}

func headersToMap(header http.Header) map[string]string {
	headers := map[string]string{}
	for key, value := range header {
		headers[key] = value[0]
	}
	return headers
}
