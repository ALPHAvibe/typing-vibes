package main

import (
	"strings"
	"time"
)

func countLines(text string) int {
	return len(strings.Split(strings.TrimSpace(text), "\n"))
}

// isCharacterCorrect checks if the character about to be typed matches the expected character
func isCharacterCorrect(currentInput, target, charToAdd string) bool {
	// Get position accounting for skipped leading whitespace
	targetRunes := []rune(target)
	inputPos := 0
	targetPos := 0
	atLineStart := true

	for targetPos < len(targetRunes) && inputPos < len([]rune(currentInput)) {
		targetChar := targetRunes[targetPos]
		isLeadingWhitespace := atLineStart && (targetChar == ' ' || targetChar == '\t')

		if !isLeadingWhitespace {
			inputPos++
			atLineStart = false
		}

		if targetChar == '\n' {
			atLineStart = true
		}

		targetPos++
	}

	// Skip any remaining leading whitespace to find the next expected character
	for targetPos < len(targetRunes) {
		targetChar := targetRunes[targetPos]
		isLeadingWhitespace := atLineStart && (targetChar == ' ' || targetChar == '\t')

		if !isLeadingWhitespace {
			break
		}

		targetPos++
	}

	// Check if the character to add matches the expected character
	if targetPos < len(targetRunes) {
		expectedChar := string(targetRunes[targetPos])
		return charToAdd == expectedChar
	}

	return false
}

// getCurrentPosition returns the current position in the target text
func getCurrentPosition(currentInput, target string) int {
	targetRunes := []rune(target)
	inputPos := 0
	targetPos := 0
	atLineStart := true

	for targetPos < len(targetRunes) && inputPos < len([]rune(currentInput)) {
		targetChar := targetRunes[targetPos]
		isLeadingWhitespace := atLineStart && (targetChar == ' ' || targetChar == '\t')

		if !isLeadingWhitespace {
			inputPos++
			atLineStart = false
		}

		if targetChar == '\n' {
			atLineStart = true
		}

		targetPos++
	}

	// Skip any remaining leading whitespace
	for targetPos < len(targetRunes) {
		targetChar := targetRunes[targetPos]
		isLeadingWhitespace := atLineStart && (targetChar == ' ' || targetChar == '\t')

		if !isLeadingWhitespace {
			break
		}

		targetPos++
	}

	return targetPos
}

func disableLigatures(text string) string {
	// Insert zero-width space (U+200B) between characters that form ligatures
	// This prevents the terminal from rendering them as ligatures
	ligatures := []struct {
		pattern     string
		replacement string
	}{
		// Common Go operators and symbols
		{":=", ":\u200B="},
		{"!=", "!\u200B="},
		{"==", "=\u200B="},
		{"<=", "<\u200B="},
		{">=", ">\u200B="},
		{"->", "-\u200B>"},
		{"<-", "<\u200B-"},
		{"||", "|\u200B|"},
		{"&&", "&\u200B&"},
		{"++", "+\u200B+"},
		{"--", "-\u200B-"},
		{"::", ":\u200B:"},
		{"..", ".\u200B."},
		{"/*", "/\u200B*"},
		{"*/", "*\u200B/"},
		{"//", "/\u200B/"},
		{"=>", "=\u200B>"},
		{"===", "=\u200B=\u200B="},
	}

	result := text
	for _, lig := range ligatures {
		result = strings.ReplaceAll(result, lig.pattern, lig.replacement)
	}

	return result
}

func normalizeText(text string) string {
	// Remove zero-width spaces (used for ligature breaking)
	text = strings.ReplaceAll(text, "\u200B", "")

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
	// This function is no longer used for real-time accuracy
	// It's kept for compatibility but real accuracy uses tracked counters
	// See calculateAccuracyFromCounters instead
	return 100.0
}

// calculateAccuracyFromCounters uses the real-time tracked correct/incorrect counts
func calculateAccuracyFromCounters(correctChars, incorrectChars int) float64 {
	total := correctChars + incorrectChars
	if total == 0 {
		return 100.0
	}
	return float64(correctChars) / float64(total) * 100
}

