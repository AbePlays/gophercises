package cmd

import (
	"fmt"
	"strconv"

	"github.com/AbePlays/gophercises/07-cli-task-manager/db"
	"github.com/AbePlays/gophercises/07-cli-task-manager/utils"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Marks a task as done",
	Run: func(cmd *cobra.Command, args []string) {
		var taskIds []int

		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				continue
			}

			taskIds = append(taskIds, id)
		}

		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, id := range taskIds {
			if id <= 0 || id > len(tasks) {
				fmt.Printf("Invalid task ID: %d\n", id)
				continue
			}

			task := tasks[id-1]
			err := db.DeleteTask(utils.Atoi(task.Key))
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
