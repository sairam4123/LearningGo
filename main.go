package main

import (
	"fmt"
	"math/rand/v2"
)

func printMinMaxBounds(text string, foundLowerBound int, foundUpperBound int, lowerBound int, upperBound int) {
	fmt.Print(text, " ")
	fmt.Print("Min Bound: ")
	if foundLowerBound != lowerBound {
		fmt.Print(foundLowerBound)
	} else {
		fmt.Print("NYF")
	}
	fmt.Print(" ")
	fmt.Print("Max Bound: ")
	if foundUpperBound != upperBound {
		fmt.Print(foundUpperBound)
	} else {
		fmt.Print("NYF")
	}
	fmt.Println()
}

func main() {
	fmt.Println("Guess the number")
	fmt.Println("Computer (Imperfect Search) vs Human")

	var specialNumber = rand.N(10000) + 1

	var guessCount = 0
	var lowerBound = 1
	var upperBound = 10000
	var computerGuessCount = 0
	var computerGuessLowerBound = -1
	var computerGuessUpperBound = 10001
	var computerGuessed = false
	for {
		printMinMaxBounds("Guess Might be located in bound", lowerBound, upperBound, -1, 10001)
		var number int
		fmt.Print("Enter your guess: ")
		var _, error = fmt.Scanln(&number)
		if error != nil {
			fmt.Println("Error getting input", error)
			return
		}
		guessCount++

		if !computerGuessed {
			var computerGuess int = (computerGuessLowerBound+computerGuessUpperBound)/2 + rand.N(10) - 1
			computerGuessCount++
			if computerGuess > specialNumber {
				fmt.Println("Computer has guessed:", computerGuess)
				fmt.Println("Computer Guess: Too high")
				computerGuessUpperBound = min(computerGuessUpperBound, computerGuess+1)
			}
			if computerGuess < specialNumber {
				fmt.Println("Computer has guessed:", computerGuess)
				fmt.Println("Computer Guess: Too low")
				computerGuessLowerBound = max(computerGuessLowerBound, computerGuess-1)
			}
			if computerGuess == specialNumber {
				fmt.Println("Computer has guessed the number in", computerGuessCount, "steps.")
				computerGuessed = true
			}
		}

		if number > specialNumber {
			fmt.Println("Your Number too high")
			upperBound = min(upperBound, number)
		}
		if number < specialNumber {
			fmt.Println("Your Number too low")
			lowerBound = max(lowerBound, number)
		}
		if number == specialNumber {
			fmt.Println("You guessed the number in", guessCount, "steps.")
			if !computerGuessed {
				fmt.Println("You were faster than computer! Super human! Current Computer Guess:", computerGuessCount)
			}
			if guessCount == computerGuessCount {
				fmt.Println("You guessed the number in the same number of steps as computer! Incredible!")
			}
			if guessCount > computerGuessCount {
				fmt.Println("Get better at guessing numbers, you guessed the number in", guessCount-computerGuessCount, "steps worse.")
			}
			return
		}
	}

}
