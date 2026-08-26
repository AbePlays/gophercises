package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/AbePlays/gophercises/01-quiz-game/utils"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		utils.Exit(fmt.Sprintf("Failed to open the csv file: %s\n", *csvFilename))
	}

	csvReader := csv.NewReader(file)
	lines, err := csvReader.ReadAll()
	if err != nil {
		utils.Exit("Failed to parse the csv file")
	}

	problems := utils.ParseLines(lines)
	correctCount := utils.Evaluate(problems, *timeLimit)

	fmt.Printf("You got %d out of %d correct\n", correctCount, len(problems))
}
