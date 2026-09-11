package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

const DIR = "sample"

var re = regexp.MustCompile("^(.+?) ([0-9]{4}) [(]([0-9]+) of ([0-9]+)[)][.](.+?)$")
var replaceString = "$2 - $1 - $3 of $4.$5"

func main() {
	var toRename []string
	filepath.Walk(DIR, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if _, err := match(info.Name()); err == nil {
			toRename = append(toRename, path)
		}

		return nil
	})

	for _, oldPath := range toRename {
		dir := filepath.Dir(oldPath)
		filename := filepath.Base(oldPath)
		newFilename, _ := match(filename)
		newPath := filepath.Join(dir, newFilename)

		fmt.Printf("mv %s => %s\n", oldPath, newPath)

		err := os.Rename(oldPath, newPath)
		if err != nil {
			fmt.Println("Error Renaming:", oldPath, newPath, err.Error())
		}
	}
}

func match(filename string) (string, error) {
	if !re.MatchString(filename) {
		return "", fmt.Errorf("%s did not match our pattern", filename)
	}

	return re.ReplaceAllString(filename, replaceString), nil
}
