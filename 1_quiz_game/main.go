package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Problem struct {
	Question string
	Answer string

}

const DEFAULT_PROBLEMS_FILE = "problems.csv"

func main() {
	filename := flag.String("f", DEFAULT_PROBLEMS_FILE, "Specifies the file to pull questions and answers from.")
	flag.Parse()

	p, err := getProblemsFrom(*filename)

	if err != nil {
		log.Fatal("Problem getting input from file")
	}

	num_correct := 0
	num_wrong := 0

	reader := bufio.NewReader(os.Stdin)
	for _, problem := range p {
		fmt.Println(problem.Question + "=")
		a, _ := reader.ReadString('\n')
		a = strings.TrimSpace(a)

		if a == problem.Answer {
			num_correct++
		} else {
			num_wrong++
		}
	}

	fmt.Printf("\n\nCorrect: %v", num_correct)
	fmt.Printf("\nIncorrect: %v", num_wrong)
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
