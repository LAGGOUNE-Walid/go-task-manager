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

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: task add <title>")
			os.Exit(1)
		}
		task, err := svc.Add(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("task added (id: %d)\n", task.ID)

	case "list":
		tasks, err := svc.store.GetAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if len(tasks) == 0 {
			fmt.Println("no tasks yet. add one with: task add <title>")
			return
		}
		fmt.Printf("%-4s  %-8s  %s\n", "ID", "STATUS", "TITLE")
		for _, task := range tasks {
			fmt.Printf("%-4d  %-8s  %s\n", task.ID, task.Status, task.Title)
		}

	case "done":
		id, err := parseID("done")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := svc.Completed(id); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("task %d marked as done\n", id)

	case "delete":
		id, err := parseID("delete")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := svc.Delete(id); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("task %d deleted\n", id)

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func parseID(cmd string) (int, error) {
	if len(os.Args) < 3 {
		return 0, fmt.Errorf("usage: task %s <id>", cmd)
	}
	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return 0, fmt.Errorf("invalid id %q: must be a number", os.Args[2])
	}
	return id, nil
}

func printUsage() {
	fmt.Print(`usage:
  task add <title>    add a new task
  task list           list all tasks
  task done <id>      mark a task as done
  task delete <id>    delete a task
`)
}
