package cmd

import (
	"github.com/jim-ww/todo-go/task"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add task...",
	Short:   "Add new tasks",
	Example: "todo add 'buy milk' 'sign petition'\ntodo add task1 task2 task3",
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"a"},
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			task.Tasks = append(task.Tasks, task.NewTask(arg))
		}
		if listAfter, _ := cmd.Flags().GetBool("list-after"); listAfter {
			task.PrintTasks()
		}
		task.WriteTasksCSV(cfg.todosFilepath)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
