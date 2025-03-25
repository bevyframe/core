package contextManager

import (
	"strings"
)

type Variable struct {
	VarType string
	VarData []byte
}

func GetVar(context map[string]Variable, addr string, varName string) (Variable, bool) {
	gName := strings.Join([]string{addr, varName}, "/")
	var0, ok := context[gName]
	return var0, ok
}

func SetVar(context *map[string]Variable, addr string, varName string, varType string, varData []byte) bool {
	gName := strings.Join([]string{addr, varName}, "/")
	(*context)[gName] = Variable{
		VarType: varType,
		VarData: varData,
	}
	return true
}
