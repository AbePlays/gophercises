package primitive

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Mode int

const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedRect
	ModeBeziers
	ModeRotatedEllipse
	ModePolygon
)

func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", mode)}
	}
}

func Transform(image io.Reader, ext string, numShapes int, opts ...func() []string) (io.Reader, error) {
	var args []string
	for _, opt := range opts {
		args = append(args, opt()...)
	}

	in, err := tempFile("in", ext)
	if err != nil {
		return nil, err
	}
	defer os.Remove(in.Name())

	out, err := tempFile("out", ext)
	if err != nil {
		return nil, err
	}
	defer os.Remove(out.Name())

	_, err = io.Copy(in, image)
	if err != nil {
		return nil, errors.New("Failed to copy image to temp file")
	}

	stdCombo, err := primitive(in.Name(), out.Name(), numShapes, args...)
	if err != nil {
		return nil, errors.New("Failed to run primitive")
	}
	fmt.Println(stdCombo)

	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, out)
	if err != nil {
		return nil, errors.New("Failed to copy output to buffer")
	}
	return b, nil
}

func primitive(inputFile, outputFile string, numShapes int, args ...string) (string, error) {
	argStr := fmt.Sprintf("-i %s -o %s -n %d", inputFile, outputFile, numShapes)
	args = append(strings.Fields(argStr), args...)
	cmd := exec.Command("primitive", args...)
	b, err := cmd.CombinedOutput()
	return string(b), err
}

func tempFile(prefix, ext string) (*os.File, error) {
	file, err := os.CreateTemp("", prefix+"_")
	if err != nil {
		return nil, errors.New("Failed to create a temporary file")
	}
	defer os.Remove(file.Name())

	return os.Create(fmt.Sprintf("%s.%s", file.Name(), ext))
}
