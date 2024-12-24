package cmd

import (
	"os"
	"path/filepath"

	"github.com/jim-ww/todo-go/task"
	"github.com/spf13/cobra"
)

type config struct {
	todosFilepath string
	listAfter     bool
}

var cfg = new(config)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "Todo-go is a fast and simple tasks organizer written in Go",
	Run:   listCmd.Run,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cfg.todosFilepath = *rootCmd.Flags().String("path", filepath.Join(os.Getenv("HOME"), ".local", "share", "todos.csv"), "path to todos file") // TODO make platform independent?
	rootCmd.PersistentFlags().BoolP("list-after", "l", true, "List tasks after adding/removing one")
	task.Tasks = task.ReadTasksCSV(cfg.todosFilepath)
}
