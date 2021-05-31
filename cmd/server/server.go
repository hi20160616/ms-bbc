package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hi20160616/ms-bbc/config"
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
		log.Println("MS start at: ", config.Data.MS.Addr)
		return server.Start(ctx)
	})
	g.Go(func() error {
		<-ctx.Done() // wait for stop signal
		return server.Stop(ctx)
	})

	// Job
	g.Go(func() error {
		log.Println("Job start.")
		return job.Crawl(ctx)
	})
	g.Go(func() error {
		<-ctx.Done() // wait for stop signal
		return job.Stop(ctx)
	})

	// Elegant stop
	c := make(chan os.Signal, 1)
	sigs := []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
	signal.Notify(c, sigs...)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case sig := <-c:
				log.Printf("signal caught: %s ready to quit...", sig.String())
				if err := server.Stop(ctx); err != nil {
					return err
				}
				if err := job.Stop(ctx); err != nil {
					return err
				}
				os.Exit(0)
			}
		}
	})
	if err := g.Wait(); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("not canceled by context: %s", err)
		} else {
			log.Printf("%#v", err)
		}
	}
}
