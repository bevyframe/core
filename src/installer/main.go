package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
)

//go:embed bevyframe.tar.gz
var embeddedTarball []byte

func main() {
	if os.Geteuid() != 0 {
		fmt.Println("Please run this installer as root.")
		return
	}

	// delete /opt/bevyframe if it exists, then re-create it anyway
	err := os.RemoveAll("/opt/bevyframe")
	if err != nil {
		return
	}
	err = os.MkdirAll("/opt/bevyframe", 0755)
	if err != nil {
		return
	}

	outputFile := "/opt/bevyframe/bevyframe.tar.gz"
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(f)

	_, err = f.Write(embeddedTarball)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	script := "#!/bin/sh\ncd /opt/bevyframe || exit\ntar -xf bevyframe.tar.gz"
	err = os.WriteFile("/opt/bevyframe/installer-script", []byte(script), 0755)
	if err != nil {
		fmt.Println("Error writing script:", err)
		return
	}
	cmd := exec.Command("/opt/bevyframe/installer-script")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running script:", err)
		return
	}
	err = os.Remove("/opt/bevyframe/installer-script")
	if err != nil {
		return
	}

	err = os.Mkdir("/opt/bevyframe/sockets/", 0777)
	if err != nil {
		return
	}
	err = os.Chmod("/opt/bevyframe/sockets/", 0777)
	if err != nil {
		return
	}

	fmt.Println("BevyFrame installed successfully.")
	fmt.Println("Add to path:")
	fmt.Println("\t/opt/bevyframe/bin")
}
