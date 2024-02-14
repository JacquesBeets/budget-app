package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02") // Change this to your desired format
}

func FormatPrice(p float64) string {
	// Convert the float to a string with two decimal places
	str := fmt.Sprintf("%.2f", p)

	// Check if the number is negative
	isNegative := str[0] == '-'
	if isNegative {
		// Remove the negative sign for now
		str = str[1:]
	}

	// Split the string into the integer and decimal parts
	parts := strings.Split(str, ".")

	// Convert the integer part to an int
	i, _ := strconv.Atoi(parts[0])

	// Format the integer part with commas
	integerPart := strconv.FormatInt(int64(i), 10)

	// Add commas as thousand separators
	for i := len(integerPart) - 3; i > 0; i -= 3 {
		integerPart = integerPart[:i] + "," + integerPart[i:]
	}

	// Add the negative sign back if necessary
	if isNegative {
		integerPart = "-" + integerPart
	}

	// Return the formatted price
	return integerPart + "." + parts[1]
}
