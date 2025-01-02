package main

import "fmt"

// Enumarated string constants
const (
	StatusPending  = "Pending"
	StatusApproved = "Approved"
	StatusRejected = "Rejected"
)

func printStatus(status string) {
	fmt.Println("Status:", status)
}

func main() {
	printStatus(StatusPending)
	printStatus(StatusApproved)
	printStatus(StatusRejected)
}
