package link

import (
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func Parse(htmlStr string) ([]Link, error) {
	links := []Link{}
	reader := strings.NewReader(htmlStr)

	doc, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	links = dfs(doc)

	return links, nil
}

func dfs(node *html.Node) []Link {
	links := []Link{}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "a" {
			href := ""
			for _, attr := range child.Attr {
				if attr.Key == "href" {
					href = attr.Val
				}
			}
			text := ""
			if child.FirstChild != nil {
				text = child.FirstChild.Data
			}
			links = append(links, Link{Href: href, Text: text})
		}
		links = append(links, dfs(child)...)
	}

	return links
}
