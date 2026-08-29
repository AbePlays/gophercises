package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/AbePlays/gophercises/03-choose-your-own-adventure/chapter"
	"github.com/AbePlays/gophercises/03-choose-your-own-adventure/handler"
)

func main() {
	fileName := flag.String("file", "gopher.json", "the JSON file with the story")
	port := flag.String("port", ":8080", "the port to listen on")

	flag.Parse()

	story, err := chapter.Load(*fileName)
	if err != nil {
		panic(err)
	}

	fmt.Println("listening on", *port)
	log.Fatal(http.ListenAndServe(*port, handler.New(story)))
}
