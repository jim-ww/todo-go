package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"time"
)

type Options struct {
	path            string
	backup          bool
	backup_path     string
	listAfterChange bool
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

func NewTask(id int, name string) *Task {
	t := time.Now()
	return &Task{
		ID:        id,
		Name:      name,
		Date:      t.Format("02.01.2006") + "/" + t.Format("03:04"),
		Completed: false,
	}
}

func (t *Task) Len() int {
	return len(t.Name)
}

var todosFilename = "todos.json"

func main() {

	tasks := ReadTasksFromDisk()

	if len(os.Args) < 2 {
		List(tasks)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "add", "a":
		checkIsEnoughArgs(3)
		tasks = Add(tasks, os.Args[2:]...)
	case "edit", "e":
		checkIsEnoughArgs(3)
		taskID := GetTaskIdFromArg(os.Args[2])
		tasks = Edit(tasks, taskID, os.Args[3])
	case "done", "d":
		checkIsEnoughArgs(3)
		for _, arg := range os.Args[2:] {
			taskId := GetTaskIdFromArg(arg)
			tasks = SetCompleted(tasks, true, taskId)
		}
	case "undone", "u":
	case "info", "i":
		checkIsEnoughArgs(3)
		taskId := GetTaskIdFromArg(os.Args[2])
		Info(tasks, taskId)
	case "list", "l":
		List(tasks)
	case "reset", "rs":
		AskForConfirmation()
		tasks = TODO{}
	case "remove", "rm":
		checkIsEnoughArgs(3)
		for _, arg := range os.Args[2:] {
			taskId := GetTaskIdFromArg(arg)
			tasks = Remove(tasks, taskId)
		}
	}
	WriteTasksToDisk(tasks)
	os.Exit(0)
}

func checkIsEnoughArgs(need int) {
	if len(os.Args) < need {
		log.Fatal("Not enough arguments")
	}
}

func AskForConfirmation() {
	os.Exit(0)
	// TODO
}

func GetTaskIdFromArg(arg string) int {
	taskID, err := strconv.Atoi(arg)
	if err != nil {
		log.Fatal("Couldn't read taskID.ID must be number")
	}
	return taskID
}

func ReadTasksFromDisk() TODO {
	file, err := os.Open(todosFilename)
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

func List(tasks TODO) {
	for _, task := range tasks {
		if task.Completed {
			fmt.Printf("%d \x1b[9m%s\x1b[0m\n", task.ID, task.Name)
		} else {
			fmt.Printf("%d %s \n", task.ID, task.Name)
		}
	}
}

func SetCompleted(tasks TODO, completed bool, taskIDs ...int) TODO {
	for _, id := range taskIDs {
		task, _ := GetTaskByID(tasks, id)
		if task != nil {
			task.Completed = completed
		}
	}
	return tasks
}

func Add(tasks TODO, tasksToAdd ...string) TODO {
	for i, task := range tasksToAdd {
		tasks = append(tasks, NewTask(i+len(tasks)+1, task))
	}
	return tasks
}

func Edit(tasks TODO, taskID int, newName string) TODO {
	for _, task := range tasks {
		if task.ID == taskID {
			task.Name = newName
		}
	}
	return tasks
}

func Remove(tasks TODO, idsToRemove ...int) TODO {
	for _, id := range idsToRemove {
		task, index := GetTaskByID(tasks, id)
		if task.ID == id {
			fmt.Println("removing " + task.Name)
			tasks = slices.Delete(tasks, index, index)
		}
	}
	return tasks
}

func Info(tasks TODO, taskID int) {
	if t, _ := GetTaskByID(tasks, taskID); t != nil {
		fmt.Printf("ID: %d\nname: %s\ncompleted: %t\ndate: %s\n", t.ID, t.Name, t.Completed, t.Date)
	}
}

func GetTaskByID(tasks TODO, id int) (*Task, int) {
	index, exists := slices.BinarySearchFunc(tasks, id, func(t *Task, i int) int {
		return t.ID - i
	})
	if exists {
		return tasks[index], index
	}
	return nil, -1
}

func WriteTasksToDisk(tasks TODO) {
	file, err := os.Create(todosFilename)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(file).Encode(tasks)
	if err != nil {
		log.Fatal(err)
	}
}

func Sort(tasks TODO) {
	sort.Sort(tasks)
}
