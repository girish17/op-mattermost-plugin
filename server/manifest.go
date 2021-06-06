// This file is automatically generated. Do not modify it manually.

package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

var manifest *model.Manifest

const manifestStr = `
{
  "id": "com.girishm.op-mattermost-plugin",
  "name": "op-mattermost-plugin",
  "description": "OpenProject plugin for Mattermost",
  "homepage_url": "https://github.com/girish17/op-mattermost-plugin",
  "support_url": "https://github.com/girish17/op-mattermost-plugin/issues",
  "release_notes_url": "https://github.com/girish17/op-mattermost-plugin/releases/tag/v0.1.0",
  "icon_path": "assets/op_logo.svg",
  "version": "0.1.0",
  "min_server_version": "5.12.0",
  "server": {
    "executables": {
      "linux-amd64": "server/dist/plugin-linux-amd64",
      "darwin-amd64": "server/dist/plugin-darwin-amd64",
      "windows-amd64": "server/dist/plugin-windows-amd64.exe"
    },
    "executable": ""
  },
  "webapp": {
    "bundle_path": "webapp/dist/main.js"
  },
  "settings_schema": {
    "header": "",
    "footer": "",
    "settings": []
  }
}
`

func init() {
	manifest = model.ManifestFromJson(strings.NewReader(manifestStr))
}
