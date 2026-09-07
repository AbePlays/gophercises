package main

import (
	"path/filepath"

	"github.com/AbePlays/gophercises/07-cli-task-manager/cmd"
	"github.com/AbePlays/gophercises/07-cli-task-manager/db"
	"github.com/mitchellh/go-homedir"
)

func main() {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	dbPath := filepath.Join(homeDir, "tasks.db")
	err = db.Init(dbPath)
	if err != nil {
		panic(err)
	}

	cmd.RootCmd.Execute()
}
