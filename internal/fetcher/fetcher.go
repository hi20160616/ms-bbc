package fetcher

import (
	"context"
	"errors"
	"log"
)

// Fetch fetch and storage all stuffs to `db/articles.json`
func Fetch() error {
	log.Println("Fetching ...")
	as, err := fetch(context.Background())
	if err != nil {
		return err
	}
	log.Println("Done")
	return storage(as)
}

// fetch fetch all articles by url set in config.json
func fetch(ctx context.Context) (as []*Article, err error) {
	links, err := fetchLinks()
	if err != nil {
		return
	}
	for _, link := range links {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			a := NewArticle()
			a, err = a.fetchArticle(link)
			if err != nil {
				if !errors.Is(err, ErrTimeOverDays) {
					log.Printf("fetch error: %v, link: %s", err, link)
				}
				err = nil
				continue
			}
			as = append(as, a)
		}
	}
	return
}
