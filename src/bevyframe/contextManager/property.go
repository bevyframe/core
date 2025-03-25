package contextManager

import "strings"

func (context ContextManager) getVar(addr string, varName string) (string, []byte, bool) {
	for key, value := range context.Context {
		if key == strings.Join([]string{addr, varName}, "/") {
			return value.VarType, value.VarData, true
		}
	}
	return "null", []byte(""), false
}
