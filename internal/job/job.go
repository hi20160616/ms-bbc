package job

import (
	"context"
	"log"
	"time"

	"github.com/hi20160616/ms-bbc/config"
	"github.com/hi20160616/ms-bbc/internal/fetcher"
	"github.com/pkg/errors"
)

var Done chan struct{} = make(chan struct{}, 1)

func Crawl(ctx context.Context) error {
	t, err := time.ParseDuration(config.Data.MS.Heartbeat)
	if err != nil {
		return err
	}
	for {
		select {
		case <-time.Tick(t):
			if err := fetcher.Fetch(); err != nil {
				if !errors.Is(err, fetcher.ErrTimeOverDays) {
					log.Printf("%#v", err)
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		case <-Done:
			ctx.Done()
		}
	}
}

// Stop is nil now
func Stop(ctx context.Context) error {
	log.Println("Job gracefully stopping.")
	Done <- struct{}{}
	// return error can define here, so it will display on frontend
	return nil
}
