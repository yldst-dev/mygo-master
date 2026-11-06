package main

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

const (
	statusCompleted  = "Completed"
	statusProcessing = "Processing"
	statusCancelled  = "Cancelled"
)

var (
	months      = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	trendMonths = append([]string{"Dec"}, months...)
)

type metric struct {
	label        string
	icon         *ui.SVG
	format       func(float64) string
	lastDecember float64
	monthly      [6]float64
}

func (m metric) at(month int) (value, previous float64) {
	if month == 0 {
		return m.monthly[0], m.lastDecember
	}
	return m.monthly[month], m.monthly[month-1]
}

func (m metric) trend() []float64 {
	return append([]float64{m.lastDecember}, m.monthly[:]...)
}

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

var metrics = []metric{
	{"Total Revenue", databaseIcon, dollars, 61200, [6]float64{67700, 80600, 76300, 97700, 114260, 128430}},
	{"Active Users", usersIcon, grouped, 17800, [6]float64{18650, 20100, 19700, 21300, 22832, 24819}},
	{"Orders", cartIcon, grouped, 2840, [6]float64{3010, 3350, 3180, 3520, 3761, 3972}},
	{"Conversion Rate", percentIcon, rate, 3.42, [6]float64{3.51, 3.68, 3.59, 3.74, 3.895, 3.86}},
}

var categories = []string{"Electronics", "Apparel", "Home & Living", "Others"}

var categoryShares = [6][]float64{
	{38, 30, 19, 13},
	{39, 30, 18, 13},
	{40, 29, 18, 13},
	{41, 29, 18, 12},
	{41, 28, 19, 12},
	{42, 28, 18, 12},
}

type order struct {
	id       string
	customer string
	status   string
	amount   float64
	date     string
}

var orders = []order{
	{"#ORD-10051", "Ethan Miller", statusCompleted, 249, "Jun 23, 2025"},
	{"#ORD-10050", "Olivia Carter", statusProcessing, 129, "Jun 23, 2025"},
	{"#ORD-10049", "Noah Davis", statusCompleted, 399, "Jun 22, 2025"},
	{"#ORD-10048", "Emma Wilson", statusCancelled, 89, "Jun 22, 2025"},
	{"#ORD-10047", "Liam Johnson", statusProcessing, 199, "Jun 21, 2025"},
	{"#ORD-10046", "Sophia Brown", statusCompleted, 319, "Jun 20, 2025"},
	{"#ORD-10045", "Mason Lee", statusCompleted, 74.5, "Jun 20, 2025"},
	{"#ORD-10044", "Ava Martinez", statusProcessing, 158, "Jun 19, 2025"},
	{"#ORD-10043", "Ethan Miller", statusCompleted, 212, "Jun 18, 2025"},
	{"#ORD-10042", "Olivia Carter", statusCancelled, 96, "Jun 18, 2025"},
	{"#ORD-10041", "Noah Davis", statusCompleted, 449, "Jun 17, 2025"},
	{"#ORD-10040", "Isabella Clark", statusCompleted, 129, "Jun 16, 2025"},
}

func filterOrders(query string) []order {
	q := strings.ToLower(strings.TrimSpace(query))
	var matched []order
	for _, o := range orders {
		if strings.Contains(strings.ToLower(o.id+" "+o.customer+" "+o.status+" "+o.date), q) {
			matched = append(matched, o)
		}
	}
	return matched
}

type customer struct {
	name   string
	orders int
	spent  float64
	last   string
}

func summarizeCustomers(list []order) []customer {
	index := map[string]int{}
	var summary []customer
	for _, o := range list {
		i, seen := index[o.customer]
		if !seen {
			i = len(summary)
			index[o.customer] = i
			summary = append(summary, customer{name: o.customer, last: o.date})
		}
		summary[i].orders++
		if o.status != statusCancelled {
			summary[i].spent += o.amount
		}
	}
	return summary
}
