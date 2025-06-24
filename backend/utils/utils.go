package utils

import (
	"math"
	"strconv"
	"strings"
)

// translating string number to int
func ToInt(val string) int {
	var myInt int = 0

	for i, char := range val {
		myInt += (int(char) - 48) * (int(math.Pow(10, float64(len(val)-i-1))))
	}
	return myInt
}

// translating string number to float
func ToFloat64(val string) float64 {
	var myFloat float64 = 0
	decimalPosition := len(val)

	for i, char := range val {
		if char == '.' {
			decimalPosition = i
			break
		}
	}
	for i, char := range val {
		if i < decimalPosition {
			myFloat += float64(int(char)-48) * (math.Pow(10, float64(decimalPosition-1-i)))
		} else if i > decimalPosition {
			myFloat += float64(int(char)-48) * 1 / (math.Pow(10, float64(i-decimalPosition)))
		}
	}
	return myFloat
}

func CalculateMonthlyInstallment(amount float64, duration int, TAEG float64) float64 {
	if amount == 0 || duration == 0 || TAEG == 0 {
		return 0
	}
	var monthlyInterestsVal float64 = amount * (TAEG / 12 * math.Pow(1+TAEG/12, float64(duration)*12)) / (math.Pow(1+TAEG/12, float64(duration)*12) - 1)

	return monthlyInterestsVal
}

func CalculateYearlyInstallment(amount float64, duration int, TAEG float64) float64 {

	return CalculateMonthlyInstallment(amount, duration, TAEG) * 12
}

func ToBool(val string) bool {
	if val[0] == 't' {
		return true
	}
	return false
}

func ToFloatArray(str string, defaultValues []float64) []float64 {
	if str == "" {
		return defaultValues
	}

	var result []float64
	parts := strings.Split(str, ",")

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if val, err := strconv.ParseFloat(trimmed, 64); err == nil {
			result = append(result, val)
		}
	}

	if len(result) == 0 {
		return defaultValues
	}

	return result
}

func SimulateSavingsCashVsMortgage(
	yearlyIncome []float64,
	installment float64,
	duration int,
	taeg float64,
	yearlyGrowthRate float64,
	mortgageAmount float64,
	housePrice float64,
) []float64 {
	savings := mortgageAmount
	savingsOverTime := make([]float64, 0, duration)
	houseSaving := housePrice - mortgageAmount // add the mortgage expenses based on the type of mortgage
	totalMortgageAmount := mortgageAmount

	// Convert percentages to decimals
	yearlyGrowthRate = yearlyGrowthRate / 100

	for year := 0; year < duration; year++ {
		var leftover float64

		leftover = yearlyIncome[year] - installment
		houseSaving += installment - totalMortgageAmount*taeg
		totalMortgageAmount -= installment - totalMortgageAmount*taeg

		if savings+leftover >= 0 {
			savings = (savings + leftover) * (1 + yearlyGrowthRate)
		} else {
			savings = savings + leftover
		}

		savingsOverTime = append(savingsOverTime, math.Round(savings+houseSaving))
	}

	return savingsOverTime
}
