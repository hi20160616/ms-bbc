package fetcher

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/hi20160616/ms-bbc/config"
)

var dbfile = filepath.Join(config.Data.RootPath, config.Data.DBPath, "articles.json")

func storage(as []*Article) error {
	data, err := json.Marshal(as)
	if err != nil {
		return err
	}
	return os.WriteFile(dbfile, data, os.FileMode(os.O_CREATE|os.O_RDWR))
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
