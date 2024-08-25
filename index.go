package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
	"bufio"
)

func getFileText(userPath string) (string, error) {
	// Open the file
	file, err := os.ReadFile(userPath)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

func handleTimeEnd(wg *sync.WaitGroup) {
	defer wg.Done()
	timer := time.NewTimer(time.Duration(3) * time.Second)
	<-timer.C
	fmt.Println("Time is up for the quiz")
	os.Exit(0)
}

func main() {
	fmt.Println("Press Enter to continue...")

	// Create a new reader from standard input
	reader := bufio.NewReader(os.Stdin)

	// Read a line of input until Enter is pressed
	// This method is waiting for a user to submit some text
	_, _ = reader.ReadString('\n')
	fmt.Print("Started\n")

	var wg sync.WaitGroup
	wg.Add(1)

	go handleTimeEnd(&wg)

	go func() {
		defer wg.Done()
		filetext, err := getFileText("problems.csv")
		if err != nil {
			log.Fatal(err)
		}

		r := csv.NewReader(strings.NewReader(filetext))
		records, err := r.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		var correctAnswers int = 0
		for i := 0; i < len(records); i += 1 {
			var curQuestion string = records[i][0]
			var curAnswer string = records[i][1]
			var userAnswer string
			fmt.Println(curQuestion)
			_, err := fmt.Scan(&userAnswer)
			if err != nil {
				log.Fatal("Something went wrong")
			}
			fmt.Println(userAnswer)
			if userAnswer == curAnswer {
				correctAnswers += 1
			}
		}
		fmt.Printf("You got %v out of %v\n", correctAnswers, len(records))
	}()

	wg.Wait()
}
