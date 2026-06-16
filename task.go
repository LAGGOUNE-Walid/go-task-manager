package main

import (
	"fmt"
	"time"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusDone    Status = "done"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskService struct {
	store Store
}

func NewTaskService(store Store) *TaskService {
	return &TaskService{
		store: store,
	}
}

func (t *Task) Done() bool {
	return t.Status == StatusDone
}

func (t *Task) Complete() {
	t.Status = StatusDone
}

func (t *TaskService) Add(title string) (Task, error) {
	var task Task
	tasks, err := t.store.GetAll()
	if err != nil {
		return Task{}, err
	}
	task.Title = title
	task.ID = nextID(tasks)
	task.Status = StatusPending
	task.CreatedAt = time.Now()

	tasks = append(tasks, task)
	err = t.store.Save(tasks)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func (t *TaskService) Completed(id int) error {
	tasks, err := t.store.GetAll()
	if err != nil {
		return err
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = StatusDone
			return t.store.Save(tasks)
		}
	}
	return fmt.Errorf("task with id %d not found", id)
}

func (t *TaskService) Delete(id int) error {
	tasks, err := t.store.GetAll()
	if err != nil {
		return err
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return t.store.Save(tasks)
		}
	}
	return fmt.Errorf("task with id %d not found", id)
}

func nextID(tasks []Task) int {
	max := 0
	for _, task := range tasks {
		if task.ID > max {
			max = task.ID
		}
	}
	return max + 1
}
