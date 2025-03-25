package main

import (
	"fmt"
	"regexp"
	"strings"
)

func matchRouting(paramStr string, path string) (map[string]string, error) {
	regexStr := strings.ReplaceAll(paramStr, "*", ".*")
	variables := regexp.MustCompile(`<(.*?)>`).FindAllStringSubmatch(regexStr, -1)
	for _, v := range variables {
		regexStr = strings.ReplaceAll(regexStr, v[0], `(?P<`+v[1]+`>.*?)`)
	}
	regexStr = strings.ReplaceAll(regexStr, "/", `\/`)
	regex := regexp.MustCompile(`^` + regexStr + `$`)
	match := regex.FindStringSubmatch(path)
	if match != nil {
		variableValues := make(map[string]string)
		for i, name := range regex.SubexpNames() {
			if i != 0 && name != "" {
				variableValues[name] = match[i]
			}
		}
		return variableValues, nil
	}
	return map[string]string{}, fmt.Errorf("%s is not a valid routing expression", paramStr)
}
