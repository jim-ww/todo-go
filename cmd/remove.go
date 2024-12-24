package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jim-ww/todo-go/task"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove id...",
	Short:   "Removes task/s by id",
	Aliases: []string{"rm"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for i, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				log.Fatalf("invalid 'id' argument, must be number(int)")
			}
			id -= i
			if id > len(task.Tasks) || id < 1 {
				fmt.Println("task id must be in range of list")
				os.Exit(1)
			}
			task.Tasks = append(task.Tasks[:id-1], task.Tasks[id:]...)
		}
		writeChanges()
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
