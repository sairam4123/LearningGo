package main

import (
	"fmt"
	"slices"
	"time"
)

func CreateTodo(todos []Todo, text string) ([]Todo, Todo) {
	todo := Todo{
		Message: text,
		Status:  Pending,
		ID:      int(time.Now().Unix()),
	}
	todos = append(todos, todo)
	return todos, todo

}

func DeleteTodo(todos []Todo, index int) ([]Todo, error) {

	if index < 0 || index >= len(todos) {
		return todos, fmt.Errorf("Choice is out of bounds")
	} // wont trigger, yet a better option

	todos = slices.Delete(todos, index, index+1)
	return todos, nil
}

func UpdateStatus(todos []Todo, index int, status Status) ([]Todo, Todo, error) {
	if index < 0 || index >= len(todos) {
		return todos, Todo{}, fmt.Errorf("Choice is out of bounds")
	}

	todo := todos[index]
	todo.Status = status
	todos[index] = todo

	return todos, todo, nil
}

