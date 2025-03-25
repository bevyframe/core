package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func mainInit() {
	name, _ := input("\nName: ")
	packageName, _ := input("Package: ")
	description, _ := input("Description: ")
	style, _ := input("Style: ")
	defaultNetwork, _ := input("Default Network: ")
	fmt.Println()
	manifest := Manifest{
		Context: "https://bevyframe.islekcaganmert.me/ns/manifest",
		App: App{
			Name:                  name,
			ShortName:             name,
			Orientation:           "any",
			Version:               "1.0.0",
			Package:               packageName,
			Style:                 style,
			Icon:                  "/assets/favicon.png",
			LoginView:             "/login",
			ShareView:             "/share",
			OfflineView:           "/offline",
			AcceptMedia:           []string{},
			AllowMultipleInstance: false,
			Shortcuts:             map[string]string{},
			Cors:                  false,
			Routing:               map[string]string{},
		},
		Publishing: Publishing{
			Description: description,
			Screenshots: []string{},
		},
		Accounts: Accounts{
			DefaultNetwork: defaultNetwork,
			Permissions:    []string{},
		},
		Development: Environment{
			Host: "0.0.0.0",
			Port: 3000,
		},
		Production: Environment{
			Host: "0.0.0.0",
			Port: 80,
		},
	}
	manifestJson, _ := json.Marshal(manifest)
	_ = os.WriteFile("manifest.json", manifestJson, 0644)
	for _, dir := range []string{"functions", "assets", "pages", "src", "strings"} {
		_ = os.Mkdir(dir, 0744)
	}
	secret := mainSecret()
	_ = os.WriteFile(".secret", []byte(secret), 0644)
}
