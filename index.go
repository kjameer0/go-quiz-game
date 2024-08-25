package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func getFileText(userPath string) (string, error) {
	// Open the file
	file, err := os.ReadFile(userPath)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

// answerChannel chan int, recordLengthChannel chan int
func handleTimeEnd(wg *sync.WaitGroup, answerChannel chan int, questionCount int, limit int) {
	defer wg.Done()
	timer := time.NewTimer(time.Duration(limit) * time.Second)
	<-timer.C
	close(answerChannel)
	correctAnswers := <-answerChannel
	fmt.Printf("\nYou got %v out of %v\n", correctAnswers, questionCount)
	os.Exit(0)
}

func main() {
	fileFlagPtr := flag.String("f", "problems.csv", "Path to a valid csv file")
	limitFlagPtr := flag.Int("limit", 30, "Seconds you have to complete the quiz")
	flag.Parse()
	fmt.Println("Press Enter to continue...")

	// Create a new reader from standard input
	reader := bufio.NewReader(os.Stdin)

	// Read a line of input until Enter is pressed
	// This method is waiting for a user to submit some text
	_, _ = reader.ReadString('\n')
	fmt.Print("Started\n")

	var wg sync.WaitGroup
	wg.Add(1)

	filetext, err := getFileText(*fileFlagPtr)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(strings.NewReader(filetext))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	answerChannel := make(chan int, 1)
	go handleTimeEnd(&wg, answerChannel, len(records), *limitFlagPtr)

	go func() {
		defer wg.Done()

		// recordLengthChannel <- len(records)
		var correctAnswers int = 0
		for i := 0; i < len(records); i += 1 {
			var curQuestion string = records[i][0]
			var curAnswer string = records[i][1]
			var userAnswer string
			fmt.Printf("Problem %v: %v = ", i+1, curQuestion)
			_, err := fmt.Scan(&userAnswer)
			if err != nil {
				log.Fatal("Something went wrong")
			}
			if userAnswer == curAnswer {
				correctAnswers += 1
				select {
				case <-answerChannel: // If there's a message, receive it to remove it.
				default: // If the channel is empty, do nothing.
				}
				answerChannel <- correctAnswers
			}
		}
		fmt.Printf("\nYou got %v out of %v\n", <-answerChannel, len(records))
	}()

	wg.Wait()
}
