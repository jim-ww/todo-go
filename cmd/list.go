package cmd

import (
	"github.com/jim-ww/todo-go/task"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Lists all tasks",
	Aliases: []string{"l"},
	Example: "todo list",
	Run: func(cmd *cobra.Command, args []string) {
		task.PrintTasks()
	},
}

func init() {
	// TODO add flag to print time
	rootCmd.AddCommand(listCmd)
}
