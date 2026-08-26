package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/AbePlays/gophercises/01-quiz-game/problem"
)

func Exit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func ParseLines(lines [][]string) []problem.Problem {
	res := make([]problem.Problem, len(lines))
	for index, line := range lines {
		res[index] = problem.Problem{
			Question: line[0],
			Answer:   strings.TrimSpace(line[1]),
		}
	}

	return res
}

func Evaluate(problems []problem.Problem) int {
	correctCount := 0

	for index, problem := range problems {
		fmt.Printf("Problem %d: %s = \n", index+1, problem.Question)
		var userAnswer string
		fmt.Scanln(&userAnswer)

		if userAnswer == problem.Answer {
			correctCount++
		}
	}

	return correctCount
}
