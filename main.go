package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"
)

type Options struct {
	path            string
	listAfterChange bool
	numbered_list   bool
	// storingFormat   ENUM // json, markdown
	// backup          bool
	// backup_path     string
}

type Task struct {
	ID        int
	Name      string
	Date      string
	Completed bool
}

type TODO []*Task

func (t TODO) Len() int {
	return len(t)
}

func (t TODO) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TODO) Less(i, j int) bool {
	return t[i].ID < t[j].ID
}

func newOptions(opts ...func(o *Options)) *Options {

	todoData := filepath.Join(os.Getenv("HOME"), ".local", "share", "todos.json")

	options := Options{
		path:            todoData,
		listAfterChange: true,
	}
	for _, fn := range opts {
		fn(&options)
	}
	return &options
}

func NewTask(id int, name string) *Task {
	t := time.Now()
	return &Task{
		ID:        id,
		Name:      name,
		Date:      t.Format("02.01.2006") + "/" + t.Format("03:04"),
		Completed: false,
	}
}

var helpInfo = `Usage: todo [COMMAND] [ARGUMENTS]
Todo-go is a fast and simple tasks organizer written in go
Example: todo list
Available commands:
    - add [TASK/s]
        adds new task/s
        Examples:
				todo add \"buy carrots\" "buy milk"
				todo add carrots milk bread
    - edit [INDEX] [EDITED TASK]
        edits an existing task/s
        Example: todo edit 1 banana
    - list
        lists all tasks
        Example: todo list
    - done [INDEX]
        marks task as done
        Example: todo done 2 3 (marks second and third tasks as completed)
    - rm [INDEX]
        removes a task
        Example: todo rm 4
    - reset
        deletes all tasks
    - restore
        restore recent backup after reset
    - sort
        sorts completed and uncompleted tasks
        Example: todo sort
    - raw [todo/done]
        prints nothing but done/incompleted tasks in plain text, useful for scripting
        Example: todo raw done
`

func main() {
	options := newOptions()
	tasks, _ := readTasksFromFile(options.path)

	if len(os.Args) < 2 {
		printTasks(tasks)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "add", "a":
		checkIsEnoughArgs(3)
		tasks = addNewTasks(tasks, os.Args[2:]...)
	case "edit", "e":
		checkIsEnoughArgs(3)
		taskID := getTaskIdFromArg(os.Args[2])
		tasks = updateTaskByID(tasks, taskID, func(t *Task) {
			t.Name = os.Args[3]
		})
	case "done", "d":
		checkIsEnoughArgs(3)
		for _, arg := range os.Args[2:] {
			taskId := getTaskIdFromArg(arg)
			tasks = setCompleted(tasks, true, taskId)
		}
	case "undone", "u":
		checkIsEnoughArgs(3)
		for _, arg := range os.Args[2:] {
			taskId := getTaskIdFromArg(arg)
			tasks = setCompleted(tasks, false, taskId)
		}
	case "info", "i":
		checkIsEnoughArgs(3)
		id := getTaskIdFromArg(os.Args[2])
		printTaskInfoByID(tasks, id)
	case "list", "l":
		printTasks(tasks)
		os.Exit(0)
	case "reset", "rs":
		tasks = TODO{}
	case "remove", "rm":
		checkIsEnoughArgs(3)
		for _, arg := range os.Args[2:] {
			taskId := getTaskIdFromArg(arg)
			tasks = deleteTaskByIDs(tasks, taskId)
		}
	default:
		printHelp()
		os.Exit(0)
	}

	if options.listAfterChange {
		printTasks(tasks)
	}
	writeTasksToFile(tasks, options.path)
	os.Exit(0)
}

func printHelp() {
	fmt.Println(helpInfo)
}

func checkIsEnoughArgs(need int) {
	if len(os.Args) < need {
		log.Fatal("Not enough arguments")
	}
}

func getTaskIdFromArg(arg string) int {
	taskID, err := strconv.Atoi(arg)
	if err != nil {
		log.Fatal("Couldn't read task ID from arg. ID must be number")
	}
	return taskID
}

func printTasks(tasks TODO) {
	for _, task := range tasks {
		printTask(task)
	}
}

func printTask(task *Task) {
	if !task.Completed {
		fmt.Printf("%d %s\n", task.ID, task.Name)
		return
	}

	strikeStart := "\x1b[9m"
	strikeEnd := "\x1b[0m"
	fmt.Printf("%d %s %s %s\n", task.ID, strikeStart, task.Name, strikeEnd)
}

func setCompleted(tasks TODO, completed bool, taskIDs ...int) TODO {
	for _, id := range taskIDs {
		_, task := getTaskByID(tasks, id)
		if task != nil {
			task.Completed = completed
		}
	}
	return tasks
}

func addNewTasks(tasks TODO, tasksToAdd ...string) TODO {
	startID := len(tasks) + 1
	for i, task := range tasksToAdd {
		tasks = append(tasks, NewTask(startID+i, task))
	}
	return tasks
}

func updateTaskByID(tasks TODO, id int, update func(t *Task)) TODO {
	_, task := getTaskByID(tasks, id)
	update(task)
	return tasks
}

func deleteTaskByIDs(tasks TODO, IDs ...int) TODO {
	for _, id := range IDs {
		i, task := getTaskByID(tasks, id)
		if task != nil {
			tasks = slices.Delete(tasks, i, i+1)
		}
	}
	return tasks
}

func printTaskInfoByID(tasks TODO, id int) {
	if _, t := getTaskByID(tasks, id); t != nil {
		fmt.Printf("ID: %d\nname: %s\ncompleted: %t\ndate: %s\n", t.ID, t.Name, t.Completed, t.Date)
	}
}

func getTaskByID(tasks TODO, id int) (int, *Task) {
	index, exists := slices.BinarySearchFunc(tasks, id, func(t *Task, i int) int {
		return t.ID - i
	})
	if exists {
		return index, tasks[index]
	}
	return -1, nil
}

func readTasksFromFile(path string) (TODO, error) {
	file, err := os.Open(path)
	if err != nil {
		return TODO{}, errors.Join(errors.New("couldn't open provided path: "+path), err)
	}

	decoder := json.NewDecoder(file)
	var tasks TODO
	err = decoder.Decode(&tasks)
	if err != nil {
		return TODO{}, errors.New("couldn't decode json tasks from " + path)
	}
	return tasks, nil
}

func writeTasksToFile(tasks TODO, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Join(errors.New("couldn't create tasks file at "+path), err)
	}
	if err := json.NewEncoder(file).Encode(tasks); err != nil {
		return errors.Join(errors.New("couldn't encode json tasks to "+path), err)
	}
	return nil
}
