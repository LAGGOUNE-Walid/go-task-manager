# Project 1: CLI Task Manager
**Duration:** Days 1–2  
**Difficulty:** Beginner → Intermediate

---

## What You're Building

A command-line task manager that lets you add, list, complete, and delete tasks. Tasks are saved to a JSON file so they persist between runs.

```
$ task add "Buy groceries"
✓ Task added (id: 1)

$ task add "Write unit tests"
✓ Task added (id: 2)

$ task list
ID  STATUS  TITLE
1   [ ]     Buy groceries
2   [ ]     Write unit tests

$ task done 1
✓ Task 1 marked as done

$ task list
ID  STATUS  TITLE
1   [x]     Buy groceries
2   [ ]     Write unit tests

$ task delete 2
✓ Task 2 deleted
```

---

## Why This Project

This is the standard "hello world" of real Go code. It forces you to use:
- Structs and methods (the Go alternative to OOP classes)
- Interfaces (the most important concept in Go)
- Error handling (no exceptions in Go — learn this early)
- File I/O with JSON
- Writing and running tests

Every Go interview will probe these fundamentals.

---

## Folder Structure

```
01-cli-task-manager/
├── main.go
├── task.go          # Task struct and business logic
├── storage.go       # Load/save from JSON file
├── task_test.go     # Unit tests
└── go.mod
```

---

## Step-by-Step Guide

### Step 1 — Initialize the module

```bash
cd 01-cli-task-manager
go mod init github.com/yourname/task-manager
```

`go.mod` is Go's equivalent of `package.json`. It tracks your module name and dependencies.

---

### Step 2 — Define the Task struct (`task.go`)

```go
package main

import "time"

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

// Done returns true if the task is completed.
func (t *Task) Done() bool {
    return t.Status == StatusDone
}

// Complete marks the task as done.
func (t *Task) Complete() {
    t.Status = StatusDone
}
```

**Key concept:** `json:"id"` are struct tags. They tell the JSON encoder/decoder what field names to use. You'll see these everywhere in Go.

---

### Step 3 — Define a Storage interface (`storage.go`)

```go
package main

import (
    "encoding/json"
    "os"
)

// Store is the interface any storage backend must satisfy.
// This is how Go does abstraction — define behavior, not structure.
type Store interface {
    GetAll() ([]Task, error)
    Save(tasks []Task) error
}

// JSONStore saves tasks to a local JSON file.
type JSONStore struct {
    FilePath string
}

func (s *JSONStore) GetAll() ([]Task, error) {
    data, err := os.ReadFile(s.FilePath)
    if err != nil {
        if os.IsNotExist(err) {
            return []Task{}, nil // first run: no file yet
        }
        return nil, err
    }

    var tasks []Task
    if err := json.Unmarshal(data, &tasks); err != nil {
        return nil, err
    }
    return tasks, nil
}

func (s *JSONStore) Save(tasks []Task) error {
    data, err := json.MarshalIndent(tasks, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(s.FilePath, data, 0644)
}
```

**Key concept:** Interfaces in Go are implicit. `JSONStore` satisfies `Store` automatically because it has the right methods — no `implements` keyword needed. This is a core Go design philosophy.

---

### Step 4 — Write the service layer (`task.go` continued)

Add a `TaskService` that uses the `Store` interface:

```go
type TaskService struct {
    store Store
}

func NewTaskService(store Store) *TaskService {
    return &TaskService{store: store}
}

func (s *TaskService) Add(title string) (Task, error) {
    tasks, err := s.store.GetAll()
    if err != nil {
        return Task{}, err
    }

    newTask := Task{
        ID:        nextID(tasks),
        Title:     title,
        Status:    StatusPending,
        CreatedAt: time.Now(),
    }

    tasks = append(tasks, newTask)
    return newTask, s.store.Save(tasks)
}

func (s *TaskService) Complete(id int) error {
    tasks, err := s.store.GetAll()
    if err != nil {
        return err
    }

    for i, t := range tasks {
        if t.ID == id {
            tasks[i].Complete()
            return s.store.Save(tasks)
        }
    }
    return fmt.Errorf("task with id %d not found", id)
}

func (s *TaskService) Delete(id int) error {
    tasks, err := s.store.GetAll()
    if err != nil {
        return err
    }

    filtered := tasks[:0] // reuse the underlying array — a common Go idiom
    found := false
    for _, t := range tasks {
        if t.ID == id {
            found = true
            continue
        }
        filtered = append(filtered, t)
    }

    if !found {
        return fmt.Errorf("task with id %d not found", id)
    }
    return s.store.Save(filtered)
}

func nextID(tasks []Task) int {
    max := 0
    for _, t := range tasks {
        if t.ID > max {
            max = t.ID
        }
    }
    return max + 1
}
```

**Key concept:** `fmt.Errorf("...", id)` is how you create errors in Go. There are no exceptions — errors are just values returned from functions. Always check them.

---

### Step 5 — Wire it up in `main.go`

```go
package main

import (
    "fmt"
    "os"
    "strconv"
)

func main() {
    store := &JSONStore{FilePath: "tasks.json"}
    svc := NewTaskService(store)

    if len(os.Args) < 2 {
        printUsage()
        os.Exit(1)
    }

    command := os.Args[1]

    switch command {
    case "add":
        if len(os.Args) < 3 {
            fmt.Println("usage: task add <title>")
            os.Exit(1)
        }
        task, err := svc.Add(os.Args[2])
        if err != nil {
            fmt.Fprintf(os.Stderr, "error: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("✓ Task added (id: %d)\n", task.ID)

    case "list":
        tasks, err := svc.store.GetAll()
        if err != nil {
            fmt.Fprintf(os.Stderr, "error: %v\n", err)
            os.Exit(1)
        }
        printTasks(tasks)

    case "done":
        id, err := parseID(os.Args)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
        if err := svc.Complete(id); err != nil {
            fmt.Fprintf(os.Stderr, "error: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("✓ Task %d marked as done\n", id)

    case "delete":
        id, err := parseID(os.Args)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
        if err := svc.Delete(id); err != nil {
            fmt.Fprintf(os.Stderr, "error: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("✓ Task %d deleted\n", id)

    default:
        printUsage()
        os.Exit(1)
    }
}

func parseID(args []string) (int, error) {
    if len(args) < 3 {
        return 0, fmt.Errorf("usage: task %s <id>", args[1])
    }
    id, err := strconv.Atoi(args[2])
    if err != nil {
        return 0, fmt.Errorf("invalid id: %s", args[2])
    }
    return id, nil
}

func printTasks(tasks []Task) {
    if len(tasks) == 0 {
        fmt.Println("No tasks yet. Add one with: task add <title>")
        return
    }
    fmt.Printf("%-4s %-8s %s\n", "ID", "STATUS", "TITLE")
    for _, t := range tasks {
        status := "[ ]"
        if t.Done() {
            status = "[x]"
        }
        fmt.Printf("%-4d %-8s %s\n", t.ID, status, t.Title)
    }
}

func printUsage() {
    fmt.Println(`Usage:
  task add <title>     Add a new task
  task list            List all tasks
  task done <id>       Mark a task as done
  task delete <id>     Delete a task`)
}
```

---

### Step 6 — Write tests (`task_test.go`)

```go
package main

import (
    "testing"
)

// mockStore is an in-memory store for testing — no file I/O needed.
type mockStore struct {
    tasks []Task
}

func (m *mockStore) GetAll() ([]Task, error) {
    return m.tasks, nil
}

func (m *mockStore) Save(tasks []Task) error {
    m.tasks = tasks
    return nil
}

func TestAdd(t *testing.T) {
    svc := NewTaskService(&mockStore{})

    task, err := svc.Add("Write tests")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if task.Title != "Write tests" {
        t.Errorf("expected title 'Write tests', got '%s'", task.Title)
    }
    if task.Status != StatusPending {
        t.Errorf("expected status pending, got %s", task.Status)
    }
    if task.ID != 1 {
        t.Errorf("expected id 1, got %d", task.ID)
    }
}

func TestComplete(t *testing.T) {
    store := &mockStore{tasks: []Task{{ID: 1, Title: "Test", Status: StatusPending}}}
    svc := NewTaskService(store)

    if err := svc.Complete(1); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    tasks, _ := store.GetAll()
    if tasks[0].Status != StatusDone {
        t.Errorf("expected task to be done")
    }
}

func TestCompleteNotFound(t *testing.T) {
    svc := NewTaskService(&mockStore{})

    err := svc.Complete(99)
    if err == nil {
        t.Error("expected an error for missing task")
    }
}

func TestDelete(t *testing.T) {
    store := &mockStore{tasks: []Task{
        {ID: 1, Title: "Keep"},
        {ID: 2, Title: "Delete me"},
    }}
    svc := NewTaskService(store)

    if err := svc.Delete(2); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    tasks, _ := store.GetAll()
    if len(tasks) != 1 {
        t.Errorf("expected 1 task, got %d", len(tasks))
    }
    if tasks[0].ID != 1 {
        t.Errorf("wrong task was deleted")
    }
}
```

Run with:
```bash
go test ./... -v
```

**Key concept:** The `mockStore` satisfies the `Store` interface without you doing anything special. This is why interfaces matter — your tests never touch the filesystem.

---

## How to Run

```bash
cd 01-cli-task-manager
go run . add "Learn Go interfaces"
go run . add "Write unit tests"
go run . list
go run . done 1
go run . list
go run . delete 2

# Run tests
go test ./... -v

# Build a binary
go build -o task .
./task list
```

---

## Concepts You Will Have Practiced

| Concept | Where |
|---|---|
| Structs and methods | `Task`, `TaskService` |
| Interfaces (implicit) | `Store` interface + `mockStore` |
| Error handling pattern | Every function returns `error` |
| Pointer receivers `*T` | `Complete()`, `Done()` |
| Slices and range loops | `nextID`, `Delete` |
| JSON encoding/decoding | `storage.go` |
| `os.Args` CLI parsing | `main.go` |
| Table-driven tests | `task_test.go` |

---

## Bonus Challenges (do these if you finish early)

1. Add a `task filter --status=done` command that shows only completed tasks
2. Add a `task edit <id> <new-title>` command
3. Add a `--file` flag to specify a custom storage path using the `flag` package
4. Replace `os.Args` parsing with the `cobra` library (`go get github.com/spf13/cobra`) — this is what real CLI tools use

---

## Interview Topics This Covers

- "What is an interface in Go?" → You implemented one
- "How does Go handle errors?" → No exceptions, return values
- "What are struct tags?" → You used `json:"..."` 
- "How do you write testable Go code?" → Dependency injection via interfaces
- "What is a pointer receiver?" → `func (t *Task) Complete()`
