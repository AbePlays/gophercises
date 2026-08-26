package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

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

func Evaluate(problems []problem.Problem, timeLimit int) int {
	correctCount := 0

	fmt.Print("Press enter to start the quiz... ")
	fmt.Scanln()

	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	ch := make(chan string)

	go func() {
		for index, problem := range problems {
			fmt.Printf("Problem %d: %s = \n", index+1, problem.Question)
			var userAnswer string
			fmt.Scanln(&userAnswer)
			ch <- userAnswer
		}
		close(ch)
	}()

	for index := range problems {
		select {
		case <-timer.C:
			return correctCount
		case a := <-ch:
			if strings.TrimSpace(a) == problems[index].Answer {
				correctCount++
			}
		}
	}

	return correctCount
}
