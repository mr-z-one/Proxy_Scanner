package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
)

type user_config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DbName   string `json:"dbName"`
	Port     int    `json:"port"`
	Proxy    string `json:"proxy"`

	Scope    []string `json:"scope"`
	OutScope []string `json:"outScope"`
}

func ReadeConfig(name string) *user_config {

	r := regexp.MustCompile(`[^.]+\.json$`)

	if !r.MatchString(name) {
		fmt.Println("[-]", "this not a json file !!")
		return nil
	}
	config, err := os.Open(name)

	if err != nil {
		fmt.Println("[-]", err)
		return nil
	}

	data, err := io.ReadAll(config)
	if err != nil {
		fmt.Println("[-]", err)
		return nil
	}

	var userConfig user_config

	err = json.Unmarshal(data, &userConfig)

	if err != nil {
		fmt.Println("[-]", err)
		return nil
	}

	return &userConfig
}
