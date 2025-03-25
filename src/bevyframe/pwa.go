package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
return {
        "name": manifest['app']['name'],
        "short_name": manifest['app']['short_name'],
        "description": manifest['publishing']['description'],
        "start_url": "/",
        "scope": "/",
        "id": manifest['app']['package'],
        "display": "standalone",
        "orientation": manifest['app']['orientation'],
        "icons": [icon_manifest(manifest['app']['icon'])],
        "shortcuts": [{"name": manifest['app']['shortcuts'][shortcut], "url": shortcut} for shortcut in manifest['app']['shortcuts']],
        "screenshots": [icon_manifest(icon) for icon in manifest['publishing']['screenshots']],
        "display_override": ["window-controls-overlay", "minimal-ui"],
        "launch_handler": {"client_mode": "navigate-new" if manifest['app']['allow_multiple_instance'] else "navigate-existing"},
        "share_target": {
            "action": manifest['app']['shareview'],
            "method": "POST",
            "enctype": "multipart/form-data",
            "params": {
                "title": "title",
                "text": "text",
                "url": "link",
                "files": [{"name": "media", "accept": manifest['app']['accept_media']}]
            }
        }
    }
*/

type Icon struct {
}

type Shortcut struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LaunchHandler struct {
	ClientMode string `json:"client_mode"`
}

type ShareTarget struct {
	Action  string            `json:"action"`
	Method  string            `json:"method"`
	Enctype string            `json:"enctype"`
	Params  map[string]string `json:"params"`
}

type PWA struct {
	Name            string        `json:"name"`
	ShortName       string        `json:"short_name"`
	Description     string        `json:"description"`
	StartUrl        string        `json:"start_url"`
	Scope           string        `json:"scope"`
	Id              string        `json:"id"`
	Display         string        `json:"display"`
	Orientation     string        `json:"orientation"`
	Icons           []Icon        `json:"icons"`
	Shortcuts       []Shortcut    `json:"shortcuts"`
	Screenshots     []Icon        `json:"screenshots"`
	DisplayOverride []string      `json:"display_override"`
	LaunchHandler   LaunchHandler `json:"launch_handler"`
	ShareTarget     ShareTarget   `json:"share_target"`
}

func (self Frame) processPWA() ([]byte, error) {
	manifest := PWA{
		Name:        self.manifest.App.Name,
		ShortName:   self.manifest.App.ShortName,
		Description: self.manifest.Publishing.Description,
		StartUrl:    "/",
		Scope:       "/",
		Id:          self.manifest.App.Package,
		Display:     "standalone",
		Orientation: self.manifest.App.Orientation,
		Icons:       []Icon{},
		Shortcuts:   []Shortcut{},
		Screenshots: []Icon{},
		DisplayOverride: []string{
			"window-controls-overlay",
			"minimal-ui",
		},
		LaunchHandler: LaunchHandler{
			ClientMode: "navigate-new",
		},
		ShareTarget: ShareTarget{
			Action:  self.manifest.App.ShareView,
			Method:  "POST",
			Enctype: "multipart/form-data",
			Params: map[string]string{
				"title": "title",
				"text":  "text",
				"url":   "link",
				"files": `[{"name": "media", "accept": self.manifest.App.AcceptMedia}]`,
			},
		},
	}
	return json.Marshal(manifest)
}

func (frame Frame) getServiceWorker() ([]byte, error) {
	file, err := os.Open("/opt/bevyframe/scripts/sw.js")
	if err != nil {
		return []byte{}, fmt.Errorf("failed to open manifest.json: %w", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println("failed to close manifest.json:", err)
		}
	}(file)
	data, err := io.ReadAll(file)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read manifest.json: %w", err)
	}
	sw := string(data)
	sw = strings.ReplaceAll(sw, "---offlineview---", frame.manifest.App.OfflineView)
	return []byte(sw), nil
}
