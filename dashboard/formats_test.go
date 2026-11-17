package main

import "testing"

func TestFormatsGroupRoundAndKeepTheSign(t *testing.T) {
	for _, c := range []struct{ got, want string }{
		{grouped(0), "0"},
		{grouped(999), "999"},
		{grouped(1000), "1,000"},
		{grouped(128430), "128,430"},
		{grouped(1234567), "1,234,567"},
		{grouped(-24819), "-24,819"},
		{grouped(1234.56), "1,235"},
		{dollars(61200), "$61,200"},
		{cents(0), "$0.00"},
		{cents(74.5), "$74.50"},
		{cents(1234.567), "$1,234.57"},
		{cents(-1234.5), "$-1,234.50"},
		{thousands(128430), "$128.4K"},
		{thousands(53900), "$53.9K"},
		{rate(3.86), "3.86%"},
		{rate(3), "3.00%"},
		{axisDollars(0), "$0"},
		{axisDollars(25000), "$25K"},
		{axisDollars(150000), "$150K"},
	} {
		if c.got != c.want {
			t.Errorf("got %q, want %q", c.got, c.want)
		}
	}
}
