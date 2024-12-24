package cmd

import (
	"github.com/jim-ww/todo-go/task"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:     "reset",
	Short:   "Deletes all tasks",
	Example: "todo reset",
	Aliases: []string{"rs"},
	Run: func(cmd *cobra.Command, args []string) {
		task.Tasks = []*task.Task{}
		writeChanges()
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
