package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func CreateBridgeScript() string {
	functions := ""
	dir, _ := os.ReadDir("./functions/")
	for _, i := range dir {
		if strings.HasSuffix(i.Name(), ".py") {
			name := i.Name()[:len(i.Name())-3]
			functions += " const " + name + " = async (...args) => {return await _bridge('" + i.Name() + "', ...args)};"
		}
	}
	functions = strings.ReplaceAll(functions, "\t", "")
	functions = strings.ReplaceAll(functions, "   ", "")
	functions = strings.ReplaceAll(functions, "   ", "")
	functions = strings.ReplaceAll(functions, "  ", "")
	functions = strings.ReplaceAll(functions, "  ", "")
	functions = strings.ReplaceAll(functions, "\n", "")
	return functions
}

func (self Context) ProcessBridgeProxy(function string, args string, reqTime string) (string, error) {
	if strings.Contains(function, "/") {
		return "{\"error\": \"Invalid function name\"}", nil
	}
	fmt.Printf("%s(...) -> ", function)
	pwd, _ := os.Getwd()
	filePath := fmt.Sprintf("%s/functions/%s", pwd, function)
	resp := self.execute(filePath, reqTime, []byte(args))
	respM := make(map[string]interface{})
	_ = json.Unmarshal([]byte(resp.body), &respM)
	if respM["error"] != nil {
		fmt.Println("error")
	} else if respM["type"] != nil {
		fmt.Println(respM["type"])
	} else {
		fmt.Println("OK")
	}
	return resp.body, nil
}
