package main

import (
	"fmt"
	"os"
	"time"
)

func convertToTimestamp(dateInput string) (int64, error) {
	dateFormat := "02-01-2006"
	date, err := time.Parse(dateFormat, dateInput)
	if err != nil {
		return 0, fmt.Errorf("Invalid date format. Please use the dd-mm-yyyy format.")
	}
	timestamp := date.Unix()
	return timestamp, nil
}

func main() {
	// Prompt the user to enter a date
	var dateInput string
	fmt.Print("Enter a date (dd-mm-yyyy): ")
	fmt.Scanln(&dateInput)

	// Convert the date to timestamp
	timestamp, err := convertToTimestamp(dateInput)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the RTC timestamp
	fmt.Printf("RTC Timestamp: %d\n", timestamp)
}

