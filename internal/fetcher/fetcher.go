package fetcher

import (
	"context"
)

// Fetch fetch and storage all stuffs to `db/articles.json`
func Fetch() error {
	as, err := fetch(context.Background())
	if err != nil {
		return err
	}
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
				return nil, err
			}
			as = append(as, a)

		}
	}
	return
}
