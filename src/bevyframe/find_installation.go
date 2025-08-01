package main

import (
	"os"
	"strings"
)

func FindInstallation() string {
	if len(os.Args) < 1 {
		return "/opt/bevyframe"
	}
	execPath, err := os.Executable()
	if err != nil {
		return "/opt/bevyframe"
	}
	installPath := strings.TrimSuffix(execPath, "/bin/bevyframe")

	return installPath
}
