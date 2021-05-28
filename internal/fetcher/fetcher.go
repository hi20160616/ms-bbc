package fetcher

// Fetch fetch and storage all stuffs to `db/articles.json`
func Fetch() error {
	as, err := fetch()
	if err != nil {
		return err
	}
	return storage(as)
}

// fetch fetch all articles by url set in config.json
func fetch() (as []*Article, err error) {
	links, err := fetchLinks()
	if err != nil {
		return
	}
	for _, link := range links {
		a := NewArticle()
		a, err = a.fetchArticle(link)
		if err != nil {
			return nil, err
		}
		as = append(as, a)
	}
	return
}
