package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var Tasks []*Task

type Task struct {
	Name      string
	Time      time.Time
	Completed bool
}

func NewTask(name string) *Task {
	return &Task{
		Name:      name,
		Time:      time.Now(),
		Completed: false,
	}
}

func (t *Task) String() string {
	if !t.Completed {
		return t.Name
	}
	strikeStart := "\x1b[9m"
	strikeEnd := "\x1b[0m"
	return fmt.Sprintf("%s %s %s\n", strikeStart, t.Name, strikeEnd)
}

func (t *Task) Info() {
	timeStr := t.Time.Format("02.01.2006") + "/" + t.Time.Format("03:04")
	fmt.Printf("\nname: %s\ncompleted: %t\ndate: %s\n", t.Name, t.Completed, timeStr)
}

var helpInfo = `Usage: todo [COMMAND] [ARGUMENTS]

Example: todo list
Available commands:
		- add [TASK/s]
        adds new task/s
        Examples:
				todo add \"buy carrots\" "buy milk"
				todo add carrots milk bread
    - - edit [INDEX] [EDITED TASK]
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

func ReadTasksFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("couldn't open path: %s: %w", path, err)
	}
	defer file.Close()
	if err = json.NewDecoder(file).Decode(&Tasks); err != nil {
		return fmt.Errorf("couldn't decode json tasks, %w", err)
	}
	return nil
}

func WriteTasksToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("couldn't create tasks file at %s: %w", path, err)
	}
	if err := json.NewEncoder(file).Encode(Tasks); err != nil {
		return fmt.Errorf("couldn't encode json tasks to %s: %w", path, err)
	}
	return nil
}
