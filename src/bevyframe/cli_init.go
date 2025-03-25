package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func mainInit() {
	name, err := input("\nName: ")
	if err != nil {
		fmt.Println(err)
	}
	packageName, err := input("Package: ")
	if err != nil {
		fmt.Println(err)
	}
	description, err := input("Description: ")
	if err != nil {
		fmt.Println(err)
	}
	style, err := input("Style: ")
	if err != nil {
		fmt.Println(err)
	}
	defaultNetwork, err := input("Default Network: ")
	if err != nil {
		fmt.Println(err)
	}
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
		SDKs: map[string]string{},
	}
	manifestJson, err := json.Marshal(manifest)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("manifest.json", manifestJson, 0644)
	if err != nil {
		fmt.Println(err)
	}
	for _, dir := range []string{"functions", "assets", "pages", "src", "strings"} {
		err = os.Mkdir(dir, 0744)
		if err != nil {
			fmt.Println(err)
		}
	}
	secret := mainSecret()
	err = os.WriteFile(".secret", []byte(secret), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
