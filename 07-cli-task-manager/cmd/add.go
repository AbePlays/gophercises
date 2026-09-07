package cmd

import (
	"fmt"
	"strings"

	"github.com/AbePlays/gophercises/07-cli-task-manager/db"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to your tasks list",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		_, err := db.CreateTask(task)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Added task: %s\n", task)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
