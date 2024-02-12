package utils

import (
	"fmt"
	"time"
)

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02") // Change this to your desired format
}

func FormatPrice(p float64) string {
	return fmt.Sprintf("%.2f", p)
}
