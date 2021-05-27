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
	"golang.org/x/net/html"
)

type Article struct {
	Id            string
	Title         string
	Content       string
	WebsiteId     string
	WebsiteDomain string
	WebsiteTitle  string
	UpdateTime    time.Time
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

func New(rawurl string) (a *Article) {
	a.WebsiteDomain = config.Data.MS.Domain
	a.WebsiteTitle = config.Data.MS.Title
	a.WebsiteId = fmt.Sprintf("%x", md5.Sum([]byte(a.WebsiteDomain)))
	return
}

func (a *Article) Get(rawurl string) (*Article, error) {
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
	// TODO: optimized partten
	a.Id = fmt.Sprintf("%x", md5.Sum([]byte(rawurl)))
	a.Title, err = a.getTitle()
	if err != nil {
		return nil, err
	}
	a.UpdateTime, err = a.getUpdateTime()
	if err != nil {
		return nil, err
	}
	a.Content, err = a.getContent()
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Article) getTitle() (string, error) {
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

// TODO: parseWithZone: youtube_web: render.go
func (a *Article) getUpdateTime() (time.Time, error) {
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
		return time.Time{}, fmt.Errorf("bbc setData got nothing.")
	}
	t := cs[0]
	return time.Parse(time.RFC3339, t)
}

func (a *Article) getContent() (string, error) {
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
	h1 := a.UpdateTime.Format("# [02.01] [1504H] " + a.Title)
	u, err := url.QueryUnescape(a.u.String())
	if err != nil {
		u = a.u.String() + "\n\nunescape url error:\n" + err.Error()
	}
	body = h1 + "\n\n" + body + "\n\n原地址：" + u
	return body, nil
}
