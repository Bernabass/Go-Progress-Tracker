package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetInput(prompt string, r *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	input, err := r.ReadString('\n')
	return strings.TrimSpace(input), err
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	StudentName, _ := GetInput("Enter your name: ", reader)
	input, _ := GetInput("Enter the number of subjects you have taken: ", reader)
	N, err := strconv.Atoi(input)
	Subjects := map[string]int{}
	total_grade := 0

	for N < 0 || err != nil {
		fmt.Println("Error: Please enter a valid number of subjects.")
		input, _ := GetInput("Enter the number of subjects you have taken: ", reader)
		N, err = strconv.Atoi(input)
	}

	for i := 0; i < N; i++ {
		command_1 := fmt.Sprintf("Enter subject - %v: ", i+1)
		subject, _ := GetInput(command_1, reader)
		command_2 := fmt.Sprintf("Enter your grade for %v: ", subject)
		input, _ := GetInput(command_2, reader)
		grade, err := strconv.Atoi(input)

		for err != nil {
			fmt.Println("Error: Please enter an integer.")
			command_2 := fmt.Sprintf("Enter your grade for %v: ", subject)
			input, _ = GetInput(command_2, reader)
			grade, err = strconv.Atoi(input)
		}

		for grade < 0 || grade > 100 {
			fmt.Println("Error: Please enter a grade within the valid range (0-100).")
			command_2 := fmt.Sprintf("Enter your grade for %v: ", subject)
			input, _ = GetInput(command_2, reader)
			grade, err = strconv.Atoi(input)
		}

		Subjects[subject] = grade
		total_grade += grade
	}

	fmt.Println("\n----------------------")
	fmt.Printf("Student's Name: %s\n", StudentName)
	fmt.Println("----------------------")
	fmt.Printf("%-15s %s\n", "SUBJECT", "GRADE")
	fmt.Println("----------------------")

	for subject, grade := range Subjects {
		fmt.Printf("%-15s %d\n", subject, grade)
	}

	average := float64(total_grade) / float64(N)
	fmt.Println("----------------------")
	fmt.Printf("AVERAGE GRADE: %.2f\n", average)
	fmt.Println("----------------------")
	fmt.Println("SHALAMGANDO !!!")
}