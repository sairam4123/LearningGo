package main

import (
	"encoding/json"
	"os"
)

func loadTodos() ([]Todo, error) {

	var todo []Todo

	data, err := os.ReadFile(TODO_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			return []Todo{}, nil
		}
	}

	err = json.Unmarshal(data, &todo)
	if err != nil {
		return []Todo{}, err
	}

	return todo, nil
}

func saveTodos(todo []Todo) error {
	data, err := json.Marshal(todo)
	if err != nil {
		return err
	}
	return os.WriteFile(TODO_FILE, data, 0644)
}
