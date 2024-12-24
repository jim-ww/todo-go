package cmd

import (
	"github.com/jim-ww/todo-go/task"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove id...",
	Short:   "Removes task/s by id",
	Aliases: []string{"r", "rm"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for i, arg := range args {
			id := extractAndCheckArgID(arg, i)
			task.Tasks = append(task.Tasks[:id-1], task.Tasks[id:]...)
		}
		writeChanges()
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
