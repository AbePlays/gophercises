package handler

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/AbePlays/gophercises/03-choose-your-own-adventure/chapter"
	"github.com/AbePlays/gophercises/03-choose-your-own-adventure/constants"
)

type Handler struct {
	story chapter.Story
	tpl   *template.Template
}

func New(story chapter.Story) http.Handler {
	tpl := template.Must(template.New("").Parse(constants.AdventureTemplate))
	return &Handler{story: story, tpl: tpl}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	path := strings.TrimSpace(request.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = strings.TrimPrefix(path, "/")

	chapter, ok := h.story[path]
	if !ok {
		http.Error(writer, "Chapter Not Found", http.StatusNotFound)
		return
	}

	if err := h.tpl.Execute(writer, chapter); err != nil {
		log.Println(err)
		http.Error(writer, "Something went wrong", http.StatusInternalServerError)
	}
}
