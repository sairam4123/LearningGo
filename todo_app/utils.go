package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func promptInput(text string) (string, error) {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(text)

	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	text = strings.TrimSpace(text)
	return text, nil
}

func promptInteger(text string) (int, error) {
	res, err := promptInput(text)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func PrintTodo(t Todo, i int) {
	fmt.Printf("%d. %s - %s (%d)\n", i+1, t.Message, t.Status.String(), t.ID)
}
