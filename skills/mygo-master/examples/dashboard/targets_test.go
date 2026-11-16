package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

func smallTargets(tt *ui.Tester, least float32, labels ...string) []string {
	var small []string
	for _, label := range labels {
		r, ok := tt.Find(label)
		if !ok || r.W < least || r.H < least {
			small = append(small, fmt.Sprintf("%s %.1fx%.1f", label, r.W, r.H))
		}
	}
	return small
}
