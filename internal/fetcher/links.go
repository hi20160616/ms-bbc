package fetcher

import (
	"net/url"

	"github.com/hi20160616/exhtml"
	"github.com/hi20160616/gears"
	"github.com/hi20160616/ms-bbc/config"
	"github.com/pkg/errors"
)

func getLinks(rawurl string) ([]string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if links, err := exhtml.ExtractLinks(u.String()); err != nil {
		return nil, errors.WithMessage(err, "cannot extract links from "+rawurl)
	} else {
		return gears.StrSliceDeDupl(links), nil
	}
}

func (f *Fetcher) fetchLinks() error {
	for _, rawurl := range config.Data.MS.URL {
		links, err := getLinks(rawurl)
		if err != nil {
			return err
		}
		f.Links = append(f.Links, links...)
	}
	return nil
}
