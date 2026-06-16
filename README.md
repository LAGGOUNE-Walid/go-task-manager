# go-task-manager

A command-line task manager written in Go. Tasks are persisted to a local JSON file.

## Usage

```bash
task add "<title>"    # add a new task
task list             # list all tasks
task done <id>        # mark a task as done
task delete <id>      # delete a task
```

## Build

```bash
go build -o task .
```

## Test

```bash
go test ./... -v
```
