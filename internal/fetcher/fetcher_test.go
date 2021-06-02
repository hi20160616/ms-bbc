package fetcher

import (
	"testing"

	"github.com/hi20160616/ms-bbc/config"
)

func TestFetch(t *testing.T) {
	if err := config.Reset("../../"); err != nil {
		t.Error(err)
	}

	if err := Fetch(); err != nil {
		t.Error(err)
	}
}
