package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"encoding/json"

	"github.com/goxlarge/gsoup/doc"
	"golang.org/x/net/html"
)

type blog struct {
	Link    string
	Title   string
	Date    string
	Author  string
	Summary string
}

func parseByToken(r io.Reader) (data []string) {
	tkn := html.NewTokenizer(r)
	var vals []string
	var isA bool
	for {
		tt := tkn.Next()
		switch {
		case tt == html.ErrorToken:
			return vals
		case tt == html.StartTagToken:
			t := tkn.Token()
			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" && strings.HasPrefix(a.Val, "/blog/") {
						isA = true
						fmt.Printf("LINK: %s", a.Val)
					}
				}
			}
		case tt == html.TextToken:
			t := tkn.Token()
			if isA {
				vals = append(vals, t.Data)
				fmt.Printf("\ttitle: %s\n", strings.TrimSpace(t.Data))
			}
			isA = false
		}
	}
}
func parse_blogs(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" && strings.HasPrefix(a.Val, "/blog/") {
				fmt.Printf("LINK: %s", a.Val)
				//<a href="/blog/survey2016-results">Go 2016 Survey Results</a>
				// n.FirstChild.Data == Go 2016 Survey Results
				fmt.Printf("\ttitle: %s\n", strings.TrimSpace(n.FirstChild.Data))
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parse_blogs(c)
	}
}

func traverseBlog(n *html.Node) *blog {
	b := blog{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "a" {
			for _, a := range c.Attr {
				if a.Key == "href" && strings.HasPrefix(a.Val, "/blog/") {
					b.Link = a.Val
					b.Title = c.FirstChild.Data
				}
			}
		}
		if c.Data == "span" && doc.CheckClass(c, "date") {
			b.Date = c.FirstChild.Data
		}
		if c.Data == "span" && doc.CheckClass(c, "author") {
			if c.FirstChild != nil {
				b.Author = c.FirstChild.Data
			}
		}
	}
	return &b
}

func main() {
	response, err := http.Get("https://go.dev/blog/all")
	if err != nil {
		panic("get error")
	}
	defer response.Body.Close()
	page, err2 := html.Parse(response.Body)
	if err2 != nil {
		panic("something wrong!")
	}
	// parse_blogs(doc)

	//parseByToken(response.Body)

	targets := doc.GetElementsByClass(page, "blogtitle")

	var results []*blog
	for _, t := range *targets {
		results = append(results, traverseBlog(t))
	}

	b, err3 := json.Marshal(results)
	if err3 != nil {
		panic("something wrong!")
	}
	fmt.Println(string(b))
}
