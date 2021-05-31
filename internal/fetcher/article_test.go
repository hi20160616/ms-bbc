package fetcher

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

// pass test
func TestFetchArticle(t *testing.T) {
	tests := []struct {
		url string
		err error
	}{
		{"https://www.bbc.com/zhongwen/simp/world-57255390", ErrTimeOverDays},
		{"https://www.bbc.com/zhongwen/simp/uk-57264136", nil},
	}
	for _, tc := range tests {
		a := NewArticle()
		a, err := a.fetchArticle(tc.url)
		if err != nil {
			if !errors.Is(err, ErrTimeOverDays) {
				t.Error(err)
			} else {
				fmt.Println("ignore pass test: ", tc.url)
			}
		} else {
			fmt.Println("pass test: ", a.Title)
		}
	}
}
