package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	var keys struct {
		Key    string
		Secret string
	}

	file, err := os.Open(".env.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decorder := json.NewDecoder(file)
	decorder.Decode(&keys)
	fmt.Println(keys)
}
