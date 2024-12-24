package cmd

import (
	"sort"

	"github.com/jim-ww/todo-go/task"
	"github.com/spf13/cobra"
)

var sortCmd = &cobra.Command{
	Use:     "sort",
	Short:   "Sorts tasks by completion status",
	Example: "todo sort\ntodo s",
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		sort.Slice(task.Tasks, func(i, j int) bool {
			// sort by completion status: incomplete (false) first, then complete (true)
			return !task.Tasks[i].Completed && task.Tasks[j].Completed
		})
		writeChanges()
	},
}

func init() {
	// TODO add flags to reverse order, order by time
	rootCmd.AddCommand(sortCmd)
}
