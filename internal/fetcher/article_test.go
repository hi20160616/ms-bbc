package fetcher

import (
	"fmt"
	"testing"
)

// pass test
func TestFetchArticle(t *testing.T) {
	rawurl := "https://www.bbc.com/zhongwen/simp/uk-57264136"
	a := NewArticle()
	a, err := a.fetchArticle(rawurl)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a.Id)
	fmt.Println(a.Title)
	fmt.Println(a.Content)
}
