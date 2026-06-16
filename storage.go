package main

import (
	"encoding/json"
	"os"
)

type Store interface {
	GetAll() ([]Task, error)
	Save(tasks []Task) error
}

type JSONStore struct {
	FilePath string
}

func (s *JSONStore) GetAll() ([]Task, error) {
	tbytes, err := os.ReadFile(s.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return []Task{}, err
	}
	var tasks []Task
	err = json.Unmarshal(tbytes, &tasks)
	if err != nil {
		return []Task{}, err
	}
	return tasks, nil
}

func (s *JSONStore) Save(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.FilePath, data, 0664)
}
