package job

import (
	"log"
	"time"

	"github.com/hi20160616/ms-bbc/config"
	"github.com/hi20160616/ms-bbc/internal/fetcher"
)

func Crawl() error {
	t, err := time.ParseDuration(config.Data.MS.Heartbeat)
	if err != nil {
		return err
	}
	for {
		select {
		case <-time.Tick(t):
			if err := fetcher.Fetch(); err != nil {
				log.Printf("%#v", err)
			}
		}
	}
}
