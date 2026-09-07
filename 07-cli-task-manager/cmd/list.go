package cmd

import (
	"fmt"

	"github.com/AbePlays/gophercises/07-cli-task-manager/db"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all of your tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(tasks) == 0 {
			fmt.Println("You have no tasks.")
			return
		}

		fmt.Println("You have the following tasks:")
		for index, task := range tasks {
			fmt.Printf("%d: %s\n", index+1, task.Value)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
