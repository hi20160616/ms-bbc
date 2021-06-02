package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var RootPath = ""

type configuration struct {
	MS       MicroService `json:"microservice"`
	RootPath string
	DBPath   string `json:"dbpath"`
}

type MicroService struct {
	Title     string   `json:"title"`
	Domain    string   `json:"domain"`
	URL       []string `json:"url"`
	Addr      string   `json:"addr"`
	Timeout   string   `json:"timeout"`
	Heartbeat string   `json:"heartbeat"`
}

var Data = &configuration{}

func setRootPath() error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	RootPath = root
	return nil
}

func get() error {
	f, err := os.ReadFile(filepath.Join(RootPath, "config/config.json"))
	if err != nil {
		return err
	}
	return json.Unmarshal(f, Data)
}

func init() {
	if err := setRootPath(); err != nil {
		log.Printf("config init error: %v", err)
	}
	if err := get(); err != nil {
		log.Printf("config get() error: %v", err)
	}
}

// Reset is for test to reset RootPath and invoke get()
func Reset(pwd string) error {
	RootPath = pwd
	return get()
}
