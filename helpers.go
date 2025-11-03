package main

import (
	"strings"
	"time"
)

func countLines(text string) int {
	return len(strings.Split(strings.TrimSpace(text), "\n"))
}

func normalizeText(text string) string {
	// Remove leading whitespace from each line for comparison
	lines := strings.Split(text, "\n")
	var normalized []string
	for _, line := range lines {
		normalized = append(normalized, strings.TrimLeft(line, " \t"))
	}
	return strings.Join(normalized, "\n")
}

func calculateWPM(text string, duration time.Duration) float64 {
	if len(text) == 0 {
		return 0
	}
	words := float64(len(strings.Fields(text)))
	minutes := duration.Minutes()
	if minutes == 0 {
		return 0
	}
	return words / minutes
}

func calculateAccuracy(target, input string) float64 {
	if len(input) == 0 {
		return 0
	}

	correct := 0
	for i := 0; i < len(input) && i < len(target); i++ {
		if input[i] == target[i] {
			correct++
		}
	}

	return float64(correct) / float64(len(target)) * 100
}
