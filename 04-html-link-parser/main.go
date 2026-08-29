package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AbePlays/gophercises/04-html-link-parser/link"
)

func main() {
	filePath := flag.String("file", "constants/ex1.html", "path to the HTML file")
	flag.Parse()

	file, err := os.ReadFile(*filePath)
	if err != nil {
		panic(err)
	}

	html := string(file)
	fmt.Println(html)

	links, err := link.Parse(html)
	if err != nil {
		panic(err)
	}

	fmt.Println(links)
}
