package main

import (
	"fmt"
)

func printMenu() {
	fmt.Println("Todo App - Go")
	fmt.Println("1. Add Todo")
	fmt.Println("2. Remove Todo")
	fmt.Println("3. Toggle Done")
	fmt.Println("4. Show Todos")
	fmt.Println("9. Exit")
	fmt.Print("Your choice?> ")
}

func printStatusPicker() {
	fmt.Println("1. Pending")
	fmt.Println("2. Done")
	fmt.Println("3. Archived")
}

func statusFromInteger(statusInt int) Status {
	switch statusInt {
	case 1:
		return Pending
	case 2:
		return Completed
	case 3:
		return Archived
	default:
		return Pending
	}
}

func SelectTodo(todos []Todo) (int, error) {
	for i, t := range todos {
		PrintTodo(t, i)
	}

	choice, err := promptInteger("Select a todo:>")
	if err != nil {
		return 0, err
	}

	if choice < 1 || choice > len(todos) {
		return 0, fmt.Errorf("Choice is not in range")
	}

	return choice - 1, nil
}

func CLIAddTodo() (Todo, error) {

	todos, err := loadTodos()
	if err != nil {
		return Todo{}, err
	}

	text, err := promptInput("Enter message or text:> ")
	if err != nil {
		return Todo{}, err
	}

	todos, todo := CreateTodo(todos, text)
	fmt.Println("Creating new todo")
	PrintTodo(todo, len(todos)-1)

	err = saveTodos(todos)
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func CLIDeleteTodo() error {
	todos, err := loadTodos()
	if err != nil {
		return err
	}
	choice, err := SelectTodo(todos)
	if err != nil {
		return err
	}

	todos, err = DeleteTodo(todos, choice)
	if err != nil {
		return err
	}

	fmt.Println("Deleted todo..")

	saveTodos(todos)
	return nil
}

func CLIUpdateStatus() error {
	todos, err := loadTodos()
	if err != nil {
		return err
	}
	choice, err := SelectTodo(todos)
	if err != nil {
		return err
	}

	statusChoice, err := promptInteger("Enter a status:> ")
	if err != nil {
		return err
	}

	status := statusFromInteger(statusChoice)
	todos, todo, err := UpdateStatus(todos, choice, status)
	if err != nil {
		return err
	}
	PrintTodo(todo, choice)

	saveTodos(todos)
	return nil
}

func CLIShowTodos() error {
	todos, err := loadTodos()
	if err != nil {
		return err
	}
	for i, t := range todos {
		PrintTodo(t, i)
	}
	return nil
}

func main() {
	fmt.Println("Welcome to the Todo App!")

	for {
		printMenu()
		choice, err := promptInteger("")
		if err != nil {
			fmt.Println("error", err)
		}

		switch choice {
		case 1:
			CLIAddTodo()
		case 2:
			CLIDeleteTodo()
		case 3:
			CLIUpdateStatus()
		case 4:
			CLIShowTodos()

		case 9:
			{
				fmt.Println("Quitting")
				return
			}
		default:
			fmt.Println("I don't understand...")
		}
	}
}
