package cmd

import (
	"slices"

	"github.com/jim-ww/todo-go/task"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove id...",
	Short:   "Removes task/s by id",
	Example: "todo remove 1 2 3\ntodo r 4",
	Aliases: []string{"r", "rm"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for i, arg := range args {
			id := extractAndCheckArgID(arg, i)
			task.Tasks = slices.Delete(task.Tasks, id-1, id)
		}
		writeChanges()
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
