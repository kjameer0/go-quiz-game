package main

import (
	"encoding/csv"
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

func handleTimeEnd(timer *time.Timer, wg *sync.WaitGroup) {
	defer wg.Done()
	<-timer.C
	fmt.Println("Time is up for the quiz")
	os.Exit(0)
}
func main() {
	fmt.Print("Press Enter to begin quiz")
	fmt.Scanf("n")
	var wg sync.WaitGroup
	wg.Add(1)
	timer := time.NewTimer(time.Duration(3) * time.Second)
	go handleTimeEnd(timer, &wg)
	filetext, err := getFileText("problems.csv")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer wg.Done()
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
			fmt.Scan(&userAnswer)
			if userAnswer == curAnswer {
				correctAnswers += 1
			}
		}
		fmt.Printf("You got %v out of %v\n", correctAnswers, len(records))
	}()
	wg.Wait()
}
