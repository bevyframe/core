package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Manifest struct {
	Context     string      `json:"@context"`
	App         App         `json:"app"`
	Publishing  Publishing  `json:"publishing"`
	Accounts    Accounts    `json:"accounts"`
	Development Environment `json:"development"`
	Production  Environment `json:"production"`
}

type App struct {
	Name                  string            `json:"name"`
	ShortName             string            `json:"short_name"`
	Orientation           string            `json:"orientation"`
	Version               string            `json:"version"`
	Package               string            `json:"package"`
	Style                 string            `json:"style"`
	Icon                  string            `json:"icon"`
	LoginView             string            `json:"loginview"`
	ShareView             string            `json:"shareview"`
	OfflineView           string            `json:"offlineview"`
	AcceptMedia           []string          `json:"accept_media"`
	AllowMultipleInstance bool              `json:"allow_multiple_instance"`
	Shortcuts             map[string]string `json:"shortcuts"`
	Cors                  bool              `json:"cors"`
	Routing               map[string]string `json:"routing"`
}

type Publishing struct {
	Description string   `json:"description"`
	Screenshots []string `json:"screenshots"`
}

type Accounts struct {
	DefaultNetwork string   `json:"default_network"`
	Permissions    []string `json:"permissions"`
}

type Environment struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func loadManifest() (*Manifest, error) {
	file, err := os.Open("manifest.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open manifest.json: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest.json: %w", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &manifest, nil
}
