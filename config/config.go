package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type configuration struct {
	MS MicroService `json:"microservice"`
}

type MicroService struct {
	Title   string   `json:"title"`
	Domain  string   `json:"domain"`
	URL     []string `json:"url"`
	Addr    string   `json:"addr"`
	Timeout string   `json:"timeout"`
}

var Data = &configuration{}

func load() error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	root = "../" // for config test
	f, err := os.ReadFile(filepath.Join(root, "config/config.json"))
	if err != nil {
		return err
	}
	return json.Unmarshal(f, Data)
}

func init() {
	if err := load(); err != nil {
		log.Printf("config init error: %#v", err)
	}
}
