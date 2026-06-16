package main

import "testing"

type mockStore struct {
	tasks []Task
}

func (s *mockStore) GetAll() ([]Task, error) {
	return s.tasks, nil
}
func (s *mockStore) Save(tasks []Task) error {
	s.tasks = tasks
	return nil
}
func TestAdd(t *testing.T) {
	svc := NewTaskService(&mockStore{})
	task, err := svc.Add("a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.Title != "a" {
		t.Errorf("expected title 'a', got %q", task.Title)
	}
	if task.Status != StatusPending {
		t.Errorf("expected StatusPending, got %s", task.Status)
	}
	if task.ID != 1 {
		t.Errorf("expected id 1, got %d", task.ID)
	}
}

func TestCompleted(t *testing.T) {
	store := &mockStore{}
	svc := NewTaskService(store)
	task, err := svc.Add("a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = svc.Completed(task.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tasks, _ := svc.store.GetAll()
	if tasks[0].Status != StatusDone {
		t.Errorf("expected status to be StatusDone")
	}
}

func TestCompletedNotFound(t *testing.T) {
	svc := NewTaskService(&mockStore{})
	err := svc.Completed(99)
	if err == nil {
		t.Errorf("expected error but got null err")
	}
}

func TestNextID(t *testing.T) {
	tests := []struct {
		name  string
		tasks []Task
		want  int
	}{
		{"empty", []Task{}, 1},
		{"one task", []Task{{ID: 1}}, 2},
		{"two task", []Task{{ID: 1}, {ID: 3}}, 4},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := nextID(tc.tasks)
			if got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	store := &mockStore{tasks: []Task{
		{ID: 1, Title: "Keep me"},
		{ID: 2, Title: "Delete me"},
	}}
	svc := NewTaskService(store)

	err := svc.Delete(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.tasks) != 1 {
		t.Fatalf("expected 1 task remaining, got %d", len(store.tasks))
	}
	if store.tasks[0].ID != 1 {
		t.Errorf("wrong task was deleted")
	}
}

func TestDeleteNotFound(t *testing.T) {
	svc := NewTaskService(&mockStore{})

	err := svc.Delete(99)
	if err == nil {
		t.Error("expected an error for non-existent task")
	}
}
