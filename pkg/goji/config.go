package goji

import (
	"github.com/shibukawa/configdir"
	"encoding/json"
)

type Config struct {
	Username string `json:"username"`
	Url      string `json:"url"`
}

var DefaultConfig = Config{
	Username: "",
	Url:      "",
}

func GetConfig() Config {
	var config Config
	configDirs := configdir.New("pd", "goji")
	folder := configDirs.QueryFolderContainsFile("goji.settings.json")
	if folder != nil {
		data, _ := folder.ReadFile("goji.settings.json")
		json.Unmarshal(data, &config)
	} else {
		config = DefaultConfig
	}
	return config
}

func SaveConfig() {
	configDirs := configdir.New("pd", "goji")

	var config Config
	data, _ := json.Marshal(&config)

	// Stores to user folder
	folders := configDirs.QueryFolders(configdir.Global)
	folders[0].WriteFile("goji.settings.json", data)
}
