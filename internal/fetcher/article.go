package fetcher

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hi20160616/exhtml"
	"github.com/hi20160616/gears"
	"github.com/hi20160616/ms-bbc/config"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Article struct {
	Id            string
	Title         string
	Content       string
	WebsiteId     string
	WebsiteDomain string
	WebsiteTitle  string
	UpdateTime    *timestamppb.Timestamp
	u             *url.URL
	raw           []byte
	doc           *html.Node
}

var timeout = func() time.Duration {
	t, err := time.ParseDuration(config.Data.MS.Timeout)
	if err != nil {
		log.Printf("timeout init error: %v", err)
		return time.Duration(1 * time.Minute)
	}
	return t
}()

func NewArticle() *Article {
	return &Article{
		WebsiteDomain: config.Data.MS.Domain,
		WebsiteTitle:  config.Data.MS.Title,
		WebsiteId:     fmt.Sprintf("%x", md5.Sum([]byte(config.Data.MS.Domain))),
	}
}

// List get all articles from database
func (a *Article) List() ([]*Article, error) {
	return load()
}

// Get read database and return the data by rawurl.
func (a *Article) Get(id string) (*Article, error) {
	as, err := load()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		if a.Id == id {
			return a, nil
		}
	}
	return nil, fmt.Errorf("no article with id: %s", id)
}

// fetchArticle fetch article by rawurl
func (a *Article) fetchArticle(rawurl string) (*Article, error) {
	var err error
	a.u, err = url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	// Dail
	a.raw, a.doc, err = exhtml.GetRawAndDoc(a.u, timeout)
	if err != nil {
		return nil, err
	}
	a.Id = fmt.Sprintf("%x", md5.Sum([]byte(rawurl)))
	a.Title, err = a.fetchTitle()
	if err != nil {
		return nil, err
	}
	a.UpdateTime, err = a.fetchUpdateTime()
	if err != nil {
		return nil, err
	}
	// filter work
	if a, err = a.filter(3); errors.Is(err, ErrTimeOverDays) {
		return nil, err
	}
	// content should be the last step to fetch
	a.Content, err = a.fetchContent()
	if err != nil {
		return nil, err
	}
	return a, nil

}

func (a *Article) fetchTitle() (string, error) {
	n := exhtml.ElementsByTag(a.doc, "title")
	if n == nil {
		return "", fmt.Errorf("getTitle error, there is no element <title>")
	}
	title := n[0].FirstChild.Data
	title = strings.ReplaceAll(title, " - BBC News 中文", "")
	title = strings.TrimSpace(title)
	gears.ReplaceIllegalChar(&title)
	return title, nil
}

func (a *Article) fetchUpdateTime() (*timestamppb.Timestamp, error) {
	metas := exhtml.MetasByName(a.doc, "article:modified_time")
	cs := []string{}
	for _, meta := range metas {
		for _, a := range meta.Attr {
			if a.Key == "content" {
				cs = append(cs, a.Val)
			}
		}
	}
	if len(cs) <= 0 {
		return nil, fmt.Errorf("bbc setData got nothing.")
	}
	t, err := time.Parse(time.RFC3339, cs[0])
	if err != nil {
		return nil, err
	}
	return timestamppb.New(t), nil
}

func shanghai(t time.Time) time.Time {
	loc := time.FixedZone("UTC", 8*60*60)
	return t.In(loc)
}

var ErrTimeOverDays error = errors.New("article update time out of range")

// filter work for ignore articles by conditions
func (a *Article) filter(days int) (*Article, error) {
	// if article time out of days, return nil and `ErrTimeOverDays`
	// param days means fetch news during days from befor now.
	during := func(days int, ts *timestamppb.Timestamp) bool {
		t := shanghai(ts.AsTime())
		if time.Now().Day()-t.Day() <= days {
			return true
		}
		return false
	}
	// if during return false rt nil, and error as ErrTimeOverDays
	if !during(days, a.UpdateTime) {
		return nil, ErrTimeOverDays
	}
	return a, nil
}

func (a *Article) fetchContent() (string, error) {
	body := ""
	// Fetch content nodes
	nodes := exhtml.ElementsByTag(a.doc, "main")
	if len(nodes) == 0 {
		return "", fmt.Errorf("err at 111L, ElementsByTag match nothing from: %s", a.u.String())
	}
	articleDoc := nodes[0]
	plist := exhtml.ElementsByTag(articleDoc, "h2", "p")
	for _, v := range plist {
		if v.FirstChild != nil {
			if v.Parent.FirstChild.Data == "h2" {
				body += fmt.Sprintf("\n** %s **  \n", v.FirstChild.Data)
			} else if v.FirstChild.Data == "b" {
				body += fmt.Sprintf("\n** %s **  \n", v.FirstChild.FirstChild.Data)
			} else {
				body += v.FirstChild.Data + "  \n"
			}
		}
	}

	// Format content
	body = strings.ReplaceAll(body, "span  \n", "")
	h1 := shanghai(a.UpdateTime.AsTime()).Format("# [02.01] [1504H] " + a.Title)
	u, err := url.QueryUnescape(a.u.String())
	if err != nil {
		u = a.u.String() + "\n\nunescape url error:\n" + err.Error()
	}
	body = h1 + "\n\n" + body + "\n\n原地址：" + u
	return body, nil
}
