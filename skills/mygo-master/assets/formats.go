package main

import (
	"fmt"

	"github.com/egoist/mygo/i18n"
)

func grouped(v float64) string {
	return i18n.FormatNumber("en", v, i18n.NumberOptions{Integer: true})
}

func dollars(v float64) string {
	return "$" + grouped(v)
}

func cents(v float64) string {
	return "$" + i18n.FormatNumber("en", v, i18n.NumberOptions{MinFractionDigits: 2, MaxFractionDigits: 2})
}

func thousands(v float64) string {
	return fmt.Sprintf("$%.1fK", v/1000)
}

func rate(v float64) string {
	return fmt.Sprintf("%.2f%%", v)
}

func axisDollars(v float64) string {
	if v == 0 {
		return "$0"
	}
	return fmt.Sprintf("$%.0fK", v/1000)
}
