package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func input(text string) (string, error) {
	var inp string
	fmt.Print(text)
	reader := bufio.NewReader(os.Stdin)
	inp, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	inp = strings.TrimSuffix(inp, "\n")
	inp = strings.TrimSuffix(inp, "\r")
	return inp, nil
}

func mainRun(isDebug bool) {
	manifest, err := loadManifest()
	if err != nil {
		fmt.Println("Failed to load manifest:", err)
		return
	}
	frame := newServer(*manifest)
	// go contextManager.Run(manifest.App.Package)
	frame.runServer(isDebug)
}

func mainSecret() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(key)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bevyframe <command>")
		os.Exit(1)
	} else {
		switch os.Args[1] {
		case "run":
			mainRun(true)
			break
		case "serve":
			mainRun(false)
			break
		case "version":
			fmt.Printf("BevyFrame 0.6 ⍺ (%s)\n", FindInstallation())
			break
		case "help":
			fmt.Println("Usage: bevyframe <command>")
			break
		case "init":
			mainInit()
			break
		case "secret":
			fmt.Println(mainSecret())
			break
		case "simulate_request":
			SimulatedRequest()
			break
		default:
			fmt.Println("Unknown command:", os.Args[1])
			os.Exit(1)
		}
	}
}
