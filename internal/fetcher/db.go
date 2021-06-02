package fetcher

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/hi20160616/ms-bbc/config"
)

var dbfile = filepath.Join(config.Data.RootPath, config.Data.DBPath, "articles.json")

func storage(as []*Article) error {
	log.Println("Storage ...")
	data, err := json.Marshal(as)
	if err != nil {
		return err
	}
	log.Println("Done")
	return os.WriteFile(dbfile, data, 0755)
}

func load() (as []*Article, err error) {
	data, err := os.ReadFile(dbfile)
	if err != nil {
		return
	}
	if err = json.Unmarshal(data, &as); err != nil {
		return
	}
	return
}
