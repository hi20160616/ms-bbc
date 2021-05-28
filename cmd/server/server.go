package server

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hi20160616/ms-bbc/internal/job"
	"github.com/hi20160616/ms-bbc/internal/server"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := server.Stop(ctx); err != nil {
			panic(err)
		}
	}(ctx)

	g, ctx := errgroup.WithContext(ctx)
	// MS
	g.Go(func() error {
		return server.Start(ctx)
	})
	g.Go(func() error {
		<-ctx.Done() // wait for stop signal
		return server.Stop(ctx)
	})
	if err := g.Wait(); err != nil {
		log.Printf("service stopped.")
		return
	}

	// Job
	g.Go(func() error {
		return job.Crawl()
	})
	g.Go(func() error {
		<-ctx.Done() // wait for stop signal
		// TODO: return job.Stop()
		return nil
	})
	if err := g.Wait(); err != nil {
		log.Printf("job stopped.")
		return
	}

	// Elegent stop
	c := make(chan os.Signal, 1)
	sigs := []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
	signal.Notify(c, sigs...)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				return server.Stop(ctx)
			}
		}
	})
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("signal caught: %s ready to quit...", err)
		return
	}
}
