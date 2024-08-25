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
	//return contents of file
	return string(file), nil
}

//goroutine for starting a timer and ending the program on timer end

func handleTimeEnd(wg *sync.WaitGroup, answerChannel chan int, questionCount int, limit int) {
	defer wg.Done()
	timer := time.NewTimer(time.Duration(limit) * time.Second)
	<-timer.C
	//do the following after timer ends
	//close answer channel(no more data to add)
	close(answerChannel)
	//read last count of correct answers
	correctAnswers := <-answerChannel
	//print formatted string containing final message to user
	fmt.Printf("\nYou got %v out of %v\n", correctAnswers, questionCount)
	//exit program
	os.Exit(0)
}

func main() {
	//recognize limit and file flags
	fileFlagPtr := flag.String("f", "problems.csv", "Path to a valid csv file")
	limitFlagPtr := flag.Int("limit", 30, "Seconds you have to complete the quiz")
	//grab actual user flags from command line
	flag.Parse()
	fmt.Println("Press Enter to continue...")

	// Create a new reader from standard input
	reader := bufio.NewReader(os.Stdin)

	// Read a line of input until Enter is pressed
	// This method is waiting for a user to submit some text
	_, _ = reader.ReadString('\n')
	//create a wait group so one of the 2 goroutines can complete
	//and not exit the program prematurely
	var wg sync.WaitGroup
	//only allow waiting for 1 routine to complete
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
	//create a channel to count correct answers from user
	answerChannel := make(chan int, 1)
	//start goroutine that will run timer
	go handleTimeEnd(&wg, answerChannel, len(records), *limitFlagPtr)
	//create and run goroutine to ask questions and receive answers
	go func() {
		//tell waitGroup this routine is done
		defer wg.Done()

		var correctAnswers int = 0
		for i := 0; i < len(records); i += 1 {
			var curQuestion string = records[i][0]
			var curAnswer string = strings.ToLower(strings.TrimSpace(records[i][1]))
			var userAnswer string
			fmt.Printf("Problem %v: %v = ", i+1, curQuestion)
			//get user input and put it at the address for userAnswer
			_, err := fmt.Scan(&userAnswer)
			if err != nil {
				log.Fatal("Something went wrong")
			}
			if userAnswer == curAnswer {
				correctAnswers += 1
				//select will respond to different events
				select {
				case <-answerChannel: // If there's a message in channel, receive it to remove it.
				default: // If the channel is empty, do nothing.
				}
				//add latest value for correctAnswers to channel
				answerChannel <- correctAnswers
			}
		}
		fmt.Printf("\nYou got %v out of %v\n", <-answerChannel, len(records))
	}()
	//Wait for one goroutine to finish
	wg.Wait()
}
