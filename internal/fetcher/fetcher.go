package fetcher

type Fetcher struct {
	Links []string
}

func (f *Fetcher) Fetch() ([]*Article, error) {
	if err := f.fetchLinks(); err != nil {
		return nil, err
	}
	return nil, nil
}
