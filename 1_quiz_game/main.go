package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Problem struct {
	Question string
	Answer string

}

const DEFAULT_PROBLEMS_FILE = "problems.csv"
const DEFAULT_TIME_LIMIT = 10

func main() {
	filename := flag.String("f", DEFAULT_PROBLEMS_FILE, "Specifies the file to pull questions and answers from.")
	timeLimit := flag.Int("t", DEFAULT_TIME_LIMIT, "Specifies the time limit")
	flag.Parse()

	p, err := getProblemsFrom(*filename)

	if err != nil {
		log.Fatal("Problem getting input from file")
	}

	num_correct := 0

	reader := bufio.NewReader(os.Stdin)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	problemLoop:
	for i, problem := range p {

		answer := make(chan string)

		go func() {
			fmt.Printf("Problem #%d: %s = \n", i+1, problem.Question)

			a, _ := reader.ReadString('\n')
			a = strings.TrimSpace(a)
			answer <- a
		}()

		select {
			case <- timer.C:
				break problemLoop
			case userAnswer := <- answer: {
				if userAnswer == problem.Answer {
					num_correct++
				}
			}

		}
	}

	fmt.Printf("\nYou got %d out of %d correct!\n", num_correct, len(p))
}

func getProblemsFrom(filename string) (problems []Problem, err error){
	f, err := os.Open(filename)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	inputData, err := r.ReadAll()

	if err != nil {
		return nil, err
	}

	problems = make([]Problem, len(inputData))

	for index, problem := range inputData {
		problems[index] = Problem{Question: problem[0], Answer: problem[1]}
	}

	return problems, nil
}
