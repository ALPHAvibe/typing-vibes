package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.showingConfig {
		return m.renderConfigView()
	}

	if m.err != nil {
		return fmt.Sprintf("\n%s\n\nError: %v\n\n%s\n",
			titleStyle.Render("⚡ Typing Vibes"),
			m.err,
			helpStyle.Render("Press Esc to quit • Ctrl+S for settings"))
	}

	if m.targetText == "" {
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s\n",
			titleStyle.Render("⚡ Typing Vibes"),
			"Press Enter to load a function from your configured folder.",
			helpStyle.Render("Ctrl+S for settings • Esc to quit"),
		)
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render("⚡ Typing Vibes"))
	b.WriteString("\n\n")

	// Calculate pane widths - make typing pane bigger
	leftWidth := 30
	rightWidth := m.width - leftWidth - 8

	// Left pane: Info panel
	leftContent := m.renderInfoPane()
	leftPane := infoPaneStyle.Width(leftWidth).Height(m.height - 10).Render(leftContent)

	// Right pane: Typing area
	rightContent := m.renderTypingPane(rightWidth)
	rightPane := typingPaneStyle.Width(rightWidth).Render(rightContent)

	// Join panes side by side
	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
	b.WriteString(panes)
	b.WriteString("\n\n")

	// Help text at bottom
	if m.finished {
		b.WriteString(helpStyle.Render("Enter for new test • Ctrl+R for new function • Ctrl+S for settings • Esc to quit"))
	} else {
		b.WriteString(helpStyle.Render("Ctrl+R for new function • Ctrl+S for settings • Esc to quit"))
	}

	return b.String()
}

func (m model) renderConfigView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("⚙️  Settings"))
	b.WriteString("\n\n")

	b.WriteString(formLabelStyle.Render("Folder Path:"))
	b.WriteString("\n")
	b.WriteString(m.configInputs[0].View())
	b.WriteString("\n\n")

	b.WriteString(formLabelStyle.Render("Minimum Lines:"))
	b.WriteString("\n")
	b.WriteString(m.configInputs[1].View())
	b.WriteString("\n\n")

	b.WriteString(formLabelStyle.Render("Maximum Lines:"))
	b.WriteString("\n")
	b.WriteString(m.configInputs[2].View())
	b.WriteString("\n\n")

	b.WriteString(formLabelStyle.Render("Max Time Limit (seconds, 0 = no limit):"))
	b.WriteString("\n")
	b.WriteString(m.configInputs[3].View())
	b.WriteString("\n\n")

	b.WriteString(helpStyle.Render("Tab/Shift+Tab to navigate • Enter to save • Ctrl+S or Esc to cancel"))

	return b.String()
}

func (m model) renderInfoPane() string {
	var b strings.Builder

	// Directory info
	b.WriteString(labelStyle.Render("📁 Directory:"))
	b.WriteString("\n")
	b.WriteString(valueStyle.Render(m.config.FolderPath))
	b.WriteString("\n\n")

	// File info
	b.WriteString(labelStyle.Render("📄 File:"))
	b.WriteString("\n")
	b.WriteString(valueStyle.Render(filepath.Base(m.currentFile)))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("📂 Path:"))
	b.WriteString("\n")
	// Truncate long paths
	displayPath := m.currentFile
	if len(displayPath) > 30 {
		displayPath = "..." + displayPath[len(displayPath)-27:]
	}
	b.WriteString(valueStyle.Render(displayPath))
	b.WriteString("\n\n")

	b.WriteString("─────────────────────────────────")
	b.WriteString("\n\n")

	// Live stats
	var elapsed time.Duration
	if m.started && !m.finished {
		elapsed = time.Since(m.startTime)
		if m.config.MaxTimeLimit > 0 {
			maxDuration := time.Duration(m.config.MaxTimeLimit) * time.Second
			if elapsed > maxDuration {
				elapsed = maxDuration
			}
		}
	} else if m.finished {
		elapsed = m.endTime.Sub(m.startTime)
	}

	// Timer (only if limit is set)
	if m.config.MaxTimeLimit > 0 {
		maxDuration := time.Duration(m.config.MaxTimeLimit) * time.Second
		remaining := maxDuration - elapsed
		if remaining < 0 {
			remaining = 0
		}
		b.WriteString(labelStyle.Render("⏱️  Time Remaining:"))
		b.WriteString("\n")
		timeColor := lipgloss.Color("86")
		if remaining < 10*time.Second {
			timeColor = lipgloss.Color("196")
		} else if remaining < 20*time.Second {
			timeColor = lipgloss.Color("214")
		}
		b.WriteString(lipgloss.NewStyle().Foreground(timeColor).Bold(true).Render(fmt.Sprintf("%.1fs", remaining.Seconds())))
		b.WriteString("\n\n")
	}

	// Duration
	b.WriteString(labelStyle.Render("⏰ Elapsed:"))
	b.WriteString("\n")
	b.WriteString(statsStyle.Render(fmt.Sprintf("%.2fs", elapsed.Seconds())))
	b.WriteString("\n\n")

	// WPM (live)
	if m.started {
		wpm := calculateWPM(m.currentInput, elapsed)
		b.WriteString(labelStyle.Render("⚡ WPM:"))
		b.WriteString("\n")
		b.WriteString(statsStyle.Render(fmt.Sprintf("%.1f", wpm)))
		b.WriteString("\n\n")
	}

	// Progress
	b.WriteString(labelStyle.Render("📊 Progress:"))
	b.WriteString("\n")
	normalizedTarget := normalizeText(m.targetText)
	normalizedInput := normalizeText(m.currentInput)
	progress := float64(len(normalizedInput)) / float64(len(normalizedTarget)) * 100
	b.WriteString(statsStyle.Render(fmt.Sprintf("%d/%d (%.1f%%)", len(normalizedInput), len(normalizedTarget), progress)))
	b.WriteString("\n\n")

	// Final accuracy if finished
	if m.finished {
		accuracy := calculateAccuracy(m.targetText, m.currentInput)
		b.WriteString(labelStyle.Render("✓ Accuracy:"))
		b.WriteString("\n")
		b.WriteString(statsStyle.Render(fmt.Sprintf("%.1f%%", accuracy)))
		b.WriteString("\n")
	}

	return b.String()
}

func (m model) renderTypingPane(width int) string {
	targetRunes := []rune(m.targetText)
	inputRunes := []rune(m.currentInput)

	// Split into lines
	targetLines := strings.Split(m.targetText, "\n")
	var result strings.Builder

	inputPos := 0
	targetPos := 0

	// Calculate cursor position in target
	cursorTargetPos := 0
	tempInputPos := 0
	tempTargetPos := 0
	tempAtLineStart := true

	for tempTargetPos < len(targetRunes) && tempInputPos < len(inputRunes) {
		targetChar := targetRunes[tempTargetPos]
		isLeadingWhitespace := tempAtLineStart && (targetChar == ' ' || targetChar == '\t')

		if !isLeadingWhitespace {
			tempInputPos++
			tempAtLineStart = false
		}

		if targetChar == '\n' {
			tempAtLineStart = true
		}

		tempTargetPos++
	}

	for tempTargetPos < len(targetRunes) {
		targetChar := targetRunes[tempTargetPos]
		isLeadingWhitespace := tempAtLineStart && (targetChar == ' ' || targetChar == '\t')

		if !isLeadingWhitespace {
			break
		}

		tempTargetPos++
	}

	cursorTargetPos = tempTargetPos

	// Process each line
	for lineIdx, targetLine := range targetLines {
		var topLine strings.Builder
		var bottomLine strings.Builder

		lineTargetRunes := []rune(targetLine)
		lineAtStart := true

		for charIdx := 0; charIdx < len(lineTargetRunes); charIdx++ {
			targetChar := lineTargetRunes[charIdx]
			isLeadingWhitespace := lineAtStart && (targetChar == ' ' || targetChar == '\t')
			isCursor := targetPos == cursorTargetPos && !m.finished

			if isLeadingWhitespace {
				// Leading whitespace
				topLine.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(string(targetChar)))
				bottomLine.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(string(targetChar)))
				targetPos++
			} else {
				lineAtStart = false

				if inputPos < len(inputRunes) {
					inputChar := inputRunes[inputPos]

					// Top line: ONLY show incorrect inputs
					if inputChar == targetChar {
						topLine.WriteString(" ") // Correct - show space
					} else {
						topLine.WriteString(incorrectStyle.Render(string(inputChar))) // Wrong - show in red
					}

					// Bottom line: target
					var targetStyle lipgloss.Style
					if inputChar == targetChar {
						targetStyle = correctStyle
					} else {
						targetStyle = incorrectStyle
					}

					if isCursor {
						targetStyle = targetStyle.Underline(true).UnderlineSpaces(true)
					}

					bottomLine.WriteString(targetStyle.Render(string(targetChar)))
					inputPos++
				} else {
					// Not typed yet
					topLine.WriteString(" ")

					style := lipgloss.NewStyle()
					if isCursor {
						style = style.Underline(true).UnderlineSpaces(true)
					}
					bottomLine.WriteString(style.Render(string(targetChar)))
				}
				targetPos++
			}
		}

		// Add the line pair to result
		result.WriteString(topLine.String())
		result.WriteString("\n")
		result.WriteString(bottomLine.String())

		// Add newline between lines (not after last line)
		if lineIdx < len(targetLines)-1 {
			result.WriteString("\n")
			targetPos++ // Account for the \n character in target

			// Skip the newline in input if we've typed it
			if inputPos < len(inputRunes) && inputRunes[inputPos] == '\n' {
				inputPos++
			}
		}
	}

	return result.String()
}
