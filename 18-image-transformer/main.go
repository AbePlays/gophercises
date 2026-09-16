package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/AbePlays/gophercises/18-image-transformer/primitive"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `<html><body>
		<form action="/upload" method="POST" enctype="multipart/form-data">
			<input type="file" name="image" />
			<button type="submit">Upload</button>
		</form>
		</body></html>`

		fmt.Fprint(w, html)
	})

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		ext := filepath.Ext(header.Filename)[1:]
		onDisk, err := tempFile("", ext)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer onDisk.Close()
		_, err = io.Copy(onDisk, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/modify/"+filepath.Base(onDisk.Name()), http.StatusFound)

	})

	mux.HandleFunc("/modify/", func(w http.ResponseWriter, r *http.Request) {
		imgPath := r.URL.Path[len("/modify/"):]
		file, err := os.Open("./img/" + filepath.Base(imgPath))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		ext := filepath.Ext(file.Name())

		modeStr := r.FormValue("mode")
		if modeStr == "" {
			renderModeChoices(w, file, ext)
			return
		}

		mode, err := strconv.Atoi(modeStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		numShapes := r.FormValue("shapes")
		if numShapes == "" {
			renderNumShapeChoices(w, file, ext, primitive.Mode(mode))
			return
		}

		numShapesInt, err := strconv.Atoi(numShapes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_ = numShapesInt

		http.Redirect(w, r, "/img/"+filepath.Base(file.Name()), http.StatusFound)
	})

	mux.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("./img/"))))

	log.Fatal(http.ListenAndServe(":3000", mux))
}

func tempFile(prefix, ext string) (*os.File, error) {
	file, err := os.CreateTemp("./img/", prefix)
	if err != nil {
		return nil, errors.New("Failed to create a temporary file")
	}
	defer os.Remove(file.Name())

	return os.Create(fmt.Sprintf("%s.%s", file.Name(), ext))
}

func genImage(file io.Reader, ext string, numShapes int, mode primitive.Mode) (string, error) {
	out, err := primitive.Transform(file, ext, numShapes, primitive.WithMode(mode))
	if err != nil {
		return "", err
	}

	outFile, err := tempFile("out", ext)
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	io.Copy(outFile, out)

	return outFile.Name(), nil
}

func renderNumShapeChoices(w http.ResponseWriter, rs io.ReadSeeker, ext string, mode primitive.Mode) {
	opts := []genOpts{
		{100, mode},
		{150, mode},
		{200, mode},
		{250, mode},
	}

	imgs, err := genImages(rs, ext, opts...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := `<html><body>
		{{range .}}
			<a href="/modify/{{.Name}}?mode={{.Mode}}&shapes={{.Shapes}}">
				<img width="25%" src="/img/{{.Name}}">
			</a>
		{{end}}
		</body></html>`
	tpl := template.Must(template.New("").Parse(html))

	var data []struct {
		Name   string
		Mode   primitive.Mode
		Shapes int
	}
	for i, img := range imgs {
		data = append(data, struct {
			Name   string
			Mode   primitive.Mode
			Shapes int
		}{
			Name:   filepath.Base(img),
			Mode:   opts[i].M,
			Shapes: opts[i].N,
		})
	}

	tpl.Execute(w, data)
}

func renderModeChoices(w http.ResponseWriter, rs io.ReadSeeker, ext string) {
	opts := []genOpts{
		{100, primitive.ModeRotatedRect},
		{100, primitive.ModeBeziers},
		{100, primitive.ModePolygon},
		{100, primitive.ModeEllipse},
	}

	imgs, err := genImages(rs, ext, opts...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := `<html><body>
		{{range .}}
			<a href="/modify/{{.Name}}?mode={{.Mode}}">
				<img width="25%" src="/img/{{.Name}}">
			</a>
		{{end}}
		</body></html>`
	tpl := template.Must(template.New("").Parse(html))

	var data []struct {
		Name string
		Mode primitive.Mode
	}
	for i, img := range imgs {
		data = append(data, struct {
			Name string
			Mode primitive.Mode
		}{
			Name: filepath.Base(img),
			Mode: opts[i].M,
		})
	}

	tpl.Execute(w, data)
}

type genOpts struct {
	N int
	M primitive.Mode
}

func genImages(rs io.ReadSeeker, ext string, opts ...genOpts) ([]string, error) {
	var ret []string
	for _, opt := range opts {
		rs.Seek(0, 0)
		file, err := genImage(rs, ext, opt.N, opt.M)
		if err != nil {
			return nil, err
		}
		ret = append(ret, file)
	}

	return ret, nil
}
