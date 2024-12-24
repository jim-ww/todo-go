package cmd

import (
	"github.com/jim-ww/todo-go/task"
	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:     "done task...",
	Short:   "Marks task/s with selected id's as 'done'",
	Example: "todo done 1 2\ntodo d 5",
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"d"},
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			task.Tasks[extractAndCheckArgID(arg, 0)-1].Completed = true
		}
		writeChanges()
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
