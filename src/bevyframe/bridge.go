package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func CreateBridgeScript() string {
	functions := ""
	dir, err := os.ReadDir("./functions/")
	if err != nil {
		fmt.Println(err)
		return "{\"error\": \"Application is broken\"}"
	}
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
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	filePath := fmt.Sprintf("%s/functions/%s", pwd, function)
	resp := self.execute(filePath, reqTime, []byte(args))
	respM := make(map[string]interface{})
	err = json.Unmarshal([]byte(resp.body), &respM)
	if err != nil {
		fmt.Println(err)
	}
	if respM["error"] != nil {
		fmt.Println("error")
	} else if respM["type"] != nil {
		fmt.Println(respM["type"])
	} else {
		fmt.Println("OK")
	}
	return resp.body, nil
}
