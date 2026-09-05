package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/AbePlays/gophercises/04-html-link-parser/link"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	rawURL := flag.String("url", "https://wrongabhishek.com", "the url you want your sitemap built for")
	maxDepth := flag.Int("maxDepth", 0, "maximum number of links to follow from the root page (0 = unlimited)")
	flag.Parse()

	pages, err := crawl(*rawURL, *maxDepth)
	if err != nil {
		panic(err)
	}

	toXml := urlset{Xmlns: xmlns}
	for _, p := range pages {
		toXml.Urls = append(toXml.Urls, loc{Value: p})
	}

	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}

	fmt.Println()
}

func crawl(rootURL string, maxDepth int) ([]string, error) {
	rootBody, rootBase, err := fetchPage(rootURL)
	if err != nil {
		return nil, err
	}
	root, err := url.Parse(rootBase)
	if err != nil {
		return nil, err
	}

	type item struct {
		url   string
		depth int
	}

	visited := map[string]bool{rootBase: true}
	queue := []item{{rootBase, 0}}
	var pages []string

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if maxDepth > 0 && current.depth >= maxDepth {
			continue
		}

		pages = append(pages, current.url)

		var body []byte
		var base string
		if current.url == rootBase {
			body, base = rootBody, rootBase
		} else {
			body, base, err = fetchPage(current.url)
			if err != nil {
				continue
			}
		}

		links, err := link.Parse(string(body))
		if err != nil {
			continue
		}

		for _, l := range links {
			abs, ok := resolve(root, l, base)
			if !ok {
				continue
			}
			if visited[abs] {
				continue
			}
			visited[abs] = true
			queue = append(queue, item{abs, current.depth + 1})
		}
	}

	return pages, nil
}

func fetchPage(rawURL string) (body []byte, base string, err error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	return body, baseURL.String(), nil
}

func resolve(root *url.URL, l link.Link, base string) (string, bool) {
	href := strings.TrimSpace(l.Href)
	if href == "" {
		return "", false
	}

	baseURL, err := url.Parse(base)
	if err != nil {
		return "", false
	}

	parsed, err := url.Parse(href)
	if err != nil {
		return "", false
	}

	resolved := baseURL.ResolveReference(parsed)
	if resolved.Scheme != root.Scheme || resolved.Host != root.Host {
		return "", false
	}
	return resolved.String(), true
}
