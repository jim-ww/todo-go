package cmd

import (
	"github.com/jim-ww/todo-go/task"
	"github.com/spf13/cobra"
)

var undoneCmd = &cobra.Command{
	Use:     "undone task...",
	Short:   "Marks task/s with selected id's as 'undone'",
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"u"},
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			task.Tasks[extractAndCheckArgID(arg, 0)-1].Completed = false
		}
		writeChanges()
	},
}

func init() {
	rootCmd.AddCommand(undoneCmd)
}
