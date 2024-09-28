package main

import "testing"

func TestNewTask(t *testing.T) {
	task := NewTask(0, "test")

	if task.ID != 0 {
		t.Errorf("task ID should be equal to 0, but got %d", task.ID)
	}

	if task.Name != "test" {
		t.Errorf("task name should be 'test', but got %s", task.Name)
	}

	if task.Completed != false {
		t.Errorf("task name should be 'false', but got %t", task.Completed)
	}

	// if task.Date != REGEX HERE {
	// 	t.Errorf("task name should be 'false', but got %t", task.Completed)
	// }
}
