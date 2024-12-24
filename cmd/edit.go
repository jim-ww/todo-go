package cmd

import (
	"github.com/jim-ww/todo-go/task"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:     "edit [id newtask]...",
	Short:   "Edit task assosiated with id/s",
	Example: "todo edit 4 'new name'\ntodo e 1 hello 2 world",
	Args:    cobra.MinimumNArgs(2),
	Aliases: []string{"e"},
	Run: func(cmd *cobra.Command, args []string) {
		for i := 0; i < len(args); i += 2 {
			task.Tasks[extractAndCheckArgID(args[i], 0)-1].Name = args[i+1]
		}
		writeChanges()
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
