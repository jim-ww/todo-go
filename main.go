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

var tasks []*Task

var todosFilename = "todos.json"

func main() {

	tasks = ReadTasksFromDisk()

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
			tasks = Done(tasks, taskId)
		}
	case "info", "i":
		checkIsEnoughArgs(3)
		taskId := GetTaskIdFromArg(os.Args[2])
		Info(tasks, taskId)
	case "list", "l":
		List(tasks)
	case "reset", "rs":
		WriteTasksToDisk([]*Task{})
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

func GetTaskIdFromArg(arg string) int {
	taskID, err := strconv.Atoi(arg)
	if err != nil {
		log.Fatal("Couldn't read taskID.ID must be number")
	}
	return taskID
}

func ReadTasksFromDisk() []*Task {
	file, err := os.Open(todosFilename)
	if err != nil {
		return []*Task{}
	}
	decoder := json.NewDecoder(file)

	var tasks []*Task

	err = decoder.Decode(&tasks)
	if err != nil {
		log.Fatal(err)
	}
	return tasks
}

func List(tasks []*Task) {
	for _, task := range tasks {
		if task.Completed {
			fmt.Printf("%d \x1b[9m%s\x1b[0m\n", task.ID, task.Name)
		} else {
			fmt.Printf("%d %s \n", task.ID, task.Name)
		}
	}
}

func Done(tasks []*Task, taskIDs ...int) []*Task {
	for _, task := range tasks {
		for _, id := range taskIDs {
			if task.ID == id {
				task.Completed = true
			}
		}
	}
	return tasks
}

func Add(tasks []*Task, tasksToAdd ...string) []*Task {
	for i, task := range tasksToAdd {
		tasks = append(tasks, NewTask(i+len(tasks)+1, task))
	}
	return tasks
}

func Edit(tasks []*Task, taskID int, newName string) []*Task {
	for _, task := range tasks {
		if task.ID == taskID {
			task.Name = newName
		}
	}
	return tasks
}

func Remove(tasks []*Task, taskIds ...int) []*Task {
	for i, task := range tasks {
		for _, idToRemove := range taskIds {
			if task.ID == idToRemove {
				fmt.Println("removing " + task.Name)
				tasks = slices.Delete(tasks, i, i)
			}
		}
	}
	return tasks
}

func Info(tasks []*Task, taskID int) {
	for _, task := range tasks {
		if task.ID == taskID {
			// TODO
			// result := slices.BinarySearch(tasks, func(task *Task) bool {
			// return task.ID == taskID
			// })
			fmt.Printf("ID: %d\nname: %s\ncompleted: %t\ndate: %s\n", task.ID, task.Name, task.Completed, task.Date)
			return
		}
	}
}

func WriteTasksToDisk(tasks []*Task) {
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
