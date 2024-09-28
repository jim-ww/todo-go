package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	options := Options{
		path:            "todos.json",
		listAfterChange: true,
		numbered_list:   true,
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

func main() {
	options := newOptions()
	tasks := readTasksFromDisk(options.path)

	if len(os.Args) < 2 {
		list(tasks)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "add", "a":
		checkIsEnoughArgs(3)
		tasks = add(tasks, os.Args[2:]...)
	case "edit", "e":
		checkIsEnoughArgs(3)
		taskID := getTaskIdFromArg(os.Args[2])
		tasks = edit(tasks, taskID, os.Args[3])
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
		taskId := getTaskIdFromArg(os.Args[2])
		Info(tasks, taskId)
	case "list", "l":
		list(tasks)
		os.Exit(0)
	case "reset", "rs":
		tasks = TODO{}
	case "remove", "rm":
		checkIsEnoughArgs(3)
		for _, arg := range os.Args[2:] {
			taskId := getTaskIdFromArg(arg)
			tasks = remove(tasks, taskId)
		}
	}
	if options.listAfterChange {
		list(tasks)
	}
	writeTasksToDisk(tasks, options.path)
	os.Exit(0)
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

func readTasksFromDisk(path string) TODO {
	file, err := os.Open(path)
	if err != nil {
		return TODO{}
	}
	decoder := json.NewDecoder(file)

	var tasks TODO
	err = decoder.Decode(&tasks)
	if err != nil {
		log.Fatal(err)
	}
	return tasks
}

func list(tasks TODO) {
	for _, task := range tasks {
		if task.Completed {
			fmt.Printf("%d \x1b[9m%s\x1b[0m\n", task.ID, task.Name)
		} else {
			fmt.Printf("%d %s \n", task.ID, task.Name)
		}
	}
}

func setCompleted(tasks TODO, completed bool, taskIDs ...int) TODO {
	for _, id := range taskIDs {
		task, _ := getTaskByID(tasks, id)
		if task != nil {
			task.Completed = completed
		}
	}
	return tasks
}

func add(tasks TODO, tasksToAdd ...string) TODO {
	for i, task := range tasksToAdd {
		tasks = append(tasks, NewTask(i+1+len(tasks), task))
	}
	return tasks
}

func edit(tasks TODO, id int, newName string) TODO {
	task, _ := getTaskByID(tasks, id)
	task.Name = newName
	return tasks
}

func remove(tasks TODO, idsToRemove ...int) TODO {
	for _, id := range idsToRemove {
		task, index := getTaskByID(tasks, id)
		if task != nil {
			tasks = slices.Delete(tasks, index, index+1)
		}
	}
	return tasks
}

func Info(tasks TODO, taskID int) {
	if t, _ := getTaskByID(tasks, taskID); t != nil {
		fmt.Printf("ID: %d\nname: %s\ncompleted: %t\ndate: %s\n", t.ID, t.Name, t.Completed, t.Date)
	}
}

func getTaskByID(tasks TODO, id int) (*Task, int) {
	index, exists := slices.BinarySearchFunc(tasks, id, func(t *Task, i int) int {
		return t.ID - i
	})
	if exists {
		return tasks[index], index
	} else {
		fmt.Println("no task found with id: ", id)
	}
	return nil, -1
}

func writeTasksToDisk(tasks TODO, path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(file).Encode(tasks)
	if err != nil {
		log.Fatal(err)
	}
}
