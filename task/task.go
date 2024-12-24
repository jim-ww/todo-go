package task

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
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
	return fmt.Sprintf("\033[9m%s\033[0m", t.Name) // strikethrough, if completed
}

func PrintTasks() {
	for i, task := range Tasks {
		fmt.Printf("\033[1m%d:\033[0m %s\n", i+1, task)
	}
}

func (t *Task) PrintInfo() {
	timeStr := t.Time.Format("02.01.2006") + "/" + t.Time.Format("03:04")
	fmt.Printf("\nname: %s\ncompleted: %t\ndate: %s\n", t.Name, t.Completed, timeStr)
}

func (t *Task) CSV() []string {
	return []string{t.Name, strconv.FormatBool(t.Completed), strconv.FormatInt(t.Time.Unix(), 10)}
}

func ReadTasksJSON(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("couldn't open path: %s: %w", filepath, err)
	}
	defer file.Close()
	if err = json.NewDecoder(file).Decode(&Tasks); err != nil {
		return fmt.Errorf("couldn't decode json tasks, %w", err)
	}
	return nil
}

func WriteTasksJSON(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("couldn't create tasks file at %s: %w", filepath, err)
	}
	defer file.Close()
	if err := json.NewEncoder(file).Encode(Tasks); err != nil {
		return fmt.Errorf("couldn't encode json tasks to %s: %w", filepath, err)
	}
	return nil
}

func ReadTasksCSV(filepath string) (tasks []*Task) {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("couldn't open path for writing: %s: %v", filepath, err)
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatalf("couldn't decode/read csv tasks: %v", err)
	}

	if len(records) > 1 {
		for _, record := range records[1:] {
			compl, _ := strconv.ParseBool(record[1])
			t, _ := strconv.Atoi(record[2])
			tasks = append(tasks, &Task{Name: record[0], Completed: compl, Time: time.Unix(int64(t), 0)})
		}
	}
	return
}

func WriteTasksCSV(filepath string) {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("couldn't open path for writing: %s: %v", filepath, err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	if err := w.Write([]string{"name", "completed", "date"}); err != nil {
		log.Fatalf("failed to write header: %v", err)
	}
	for _, task := range Tasks {
		if err = w.Write(task.CSV()); err != nil {
			log.Fatalf("failed to write task entry: %s: %v\n", task.Name, err)
		}
	}
}
