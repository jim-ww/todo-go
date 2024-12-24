package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

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

func extractAndCheckArgID(arg string, lenOffcet int) int {
	id, err := strconv.Atoi(arg)
	if err != nil {
		log.Fatalf("invalid 'id' argument, must be number(int)")
	}
	id -= lenOffcet
	if id > len(task.Tasks) || id < 1 {
		fmt.Println("task id must be in range of list")
		os.Exit(1)
	}
	return id
}

func writeChanges() {
	task.WriteTasksCSV(cfg.todosFilepath)
	if cfg.listAfter {
		task.PrintTasks()
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cfg.todosFilepath = *rootCmd.PersistentFlags().StringP("path", "p", filepath.Join(os.Getenv("HOME"), ".local", "share", "todos.csv"), "path to todos file") // TODO make platform independent?
	cfg.listAfter = *rootCmd.PersistentFlags().BoolP("list-after", "l", true, "List tasks after adding/removing one")
	// TODO
	//cfg.listAfter, _ = cmd.Flags().GetBool("list-after")
	//cfg.todosFilepath, _ = cmd.Flags().GetString("path")
	task.Tasks = task.ReadTasksCSV(cfg.todosFilepath)
}
