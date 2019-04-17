package configdir_test

import (
	"encoding/json"
	"github.com/shibukawa/configdir"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

type Config struct {
	UserName string `json:"user-name"`
}

var DefaultConfig = Config{
	UserName: "baron", // Do you remember BeOS?
}

// Sample for reading configuration
func ExampleConfigDir() {

	var config Config

	configDirs := configdir.New("vendor-name", "application-name")
	// optional: local path has the highest priority
	configDirs.LocalPath, _ = filepath.Abs(".")
	folder := configDirs.QueryFolderContainsFile("setting.json")
	if folder != nil {
		data, _ := folder.ReadFile("setting.json")
		json.Unmarshal(data, &config)
	} else {
		config = DefaultConfig
	}
}

// Sample for reading configuration
func ExampleConfigDir_QueryFolders() {
	configDirs := configdir.New("vendor-name", "application-name")

	var config Config
	data, _ := json.Marshal(&config)

	// Stores to local folder
	folders := configDirs.QueryFolders(configdir.Local)
	folders[0].WriteFile("setting.json", data)

	// Stores to user folder
	folders = configDirs.QueryFolders(configdir.Global)
	folders[0].WriteFile("setting.json", data)

	// Stores to system folder
	folders = configDirs.QueryFolders(configdir.System)
	folders[0].WriteFile("setting.json", data)
}

// Sample for getting cache folder
func ExampleConfigDir_QueryCacheFolder() {
	configDirs := configdir.New("vendor-name", "application-name")
	cache := configDirs.QueryCacheFolder()

	resp, err := http.Get("http://examples.com/sdk.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	cache.WriteFile("sdk.zip", body)
}
