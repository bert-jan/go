package main

import "fmt"

// Define a custom type
type Distance float64

// Add method to custom type
func (d Distance) ToMiles() float64 {
	return float64(d) * 0.621371
}

func main() {
	var d Distance = 10.0 // Distance in kilometers
	fmt.Printf("%.2f km is %.2f miles\n", d, d.ToMiles())
}
