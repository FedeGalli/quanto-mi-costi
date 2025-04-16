package utils

import (
	"math"
)

// translating string number to int
func ToInt(val string) int {
	var myInt int = 0

	for i, char := range val {
		myInt += (int(char) - 48) * (int(math.Pow(10, float64(len(val)-i))))
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

func CalculateTotalInterests(amount float64, duration float64, TAEG float64) float64 {
	if amount == 0 {
		return 0
	}

	var interestsVal float64 = 0
	var yearlyCapitalAmount float64 = amount / duration
	for duration >= 0 {
		interestsVal += amount * TAEG
		amount -= yearlyCapitalAmount
		duration--
	}

	return interestsVal
}

func ToBool(val string) bool {
	if val[0] == 't' {
		return true
	}
	return false
}
