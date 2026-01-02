package main

import (
	"fmt"
	"math"
)

// thing := 10

func main() {
	investmentAmount := 1000
	expectedReturnRate := 5.5
	years := 10
	futureValue := float64(investmentAmount) * math.Pow(1+(expectedReturnRate/100), float64(years))
	fmt.Println("Future Value:", futureValue)
}
