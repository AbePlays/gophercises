package cmd

import (
	"fmt"
	"strconv"

	"github.com/AbePlays/gophercises/07-cli-task-manager/db"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Marks a task as done",
	Run: func(cmd *cobra.Command, args []string) {
		var taskKeys []string

		for _, arg := range args {
			if _, err := strconv.Atoi(arg); err != nil {
				continue
			}
			taskKeys = append(taskKeys, arg)
		}

		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, key := range taskKeys {
			id, _ := strconv.Atoi(key)
			if id <= 0 || id > len(tasks) {
				fmt.Printf("Invalid task ID: %d\n", id)
				continue
			}

			task := tasks[id-1]
			err := db.DeleteTaskByKey(task.Key)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Task %d marked as done\n", id)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}
