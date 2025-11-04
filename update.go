package main

import (
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tickMsg:
		if m.started && !m.finished {
			if m.config.MaxTimeLimit > 0 {
				elapsed := time.Since(m.startTime)
				maxDuration := time.Duration(m.config.MaxTimeLimit) * time.Second
				if elapsed >= maxDuration {
					m.finished = true
					m.endTime = m.startTime.Add(maxDuration)
					return m, nil
				}
			}
			// Continue ticking to update elapsed time
			return m, tickCmd()
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEsc:
			if m.showingConfig {
				// Cancel settings without saving
				m.showingConfig = false
				return m, nil
			}
			return m, tea.Quit

		case tea.KeyCtrlS:
			// Toggle settings
			if !m.showingConfig {
				m.showingConfig = true
				m.focusIndex = 0
				m.configInputs[0].Focus()
				for i := 1; i < len(m.configInputs); i++ {
					m.configInputs[i].Blur()
				}
			} else {
				// Cancel settings without saving
				m.showingConfig = false
			}
			return m, nil

		case tea.KeyEnter:
			if m.showingConfig {
				// Save config
				minLines, _ := strconv.Atoi(m.configInputs[1].Value())
				maxLines, _ := strconv.Atoi(m.configInputs[2].Value())
				maxTime, _ := strconv.Atoi(m.configInputs[3].Value())

				m.config = config{
					FolderPath:   m.configInputs[0].Value(),
					MinLines:     minLines,
					MaxLines:     maxLines,
					MaxTimeLimit: maxTime,
				}

				if err := saveConfig(m.config); err != nil {
					m.err = err
				}

				m.showingConfig = false

				// Load new function with new settings
				if m.targetText != "" {
					funcText, filePath, err := loadRandomFunction(m.config)
					if err != nil {
						m.err = err
					} else {
						m.targetText = funcText
						m.currentFile = filePath
						m.currentInput = ""
						m.started = false
						m.finished = false
						m.correctChars = 0
						m.incorrectChars = 0
						m.errorPositions = make(map[int]bool)
					}
				}
				return m, nil
			}

			if m.targetText == "" {
				// Initial load
				funcText, filePath, err := loadRandomFunction(m.config)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.targetText = funcText
				m.currentFile = filePath
				m.currentInput = ""
				m.correctChars = 0
				m.incorrectChars = 0
				m.errorPositions = make(map[int]bool)
				return m, nil
			}

			if m.finished {
				// Reset for another round
				funcText, filePath, err := loadRandomFunction(m.config)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.targetText = funcText
				m.currentFile = filePath
				m.currentInput = ""
				m.started = false
				m.finished = false
				m.correctChars = 0
				m.incorrectChars = 0
				m.errorPositions = make(map[int]bool)
				return m, nil
			}

			// Handle Enter during typing - add newline (rendering will skip whitespace)
			// Only allow if not finished
			if !m.finished && m.targetText != "" {
				// Find current position in target (accounting for skipped whitespace)
				targetPos := 0
				inputPos := 0
				atLineStart := true
				targetRunes := []rune(m.targetText)
				inputRunes := []rune(m.currentInput)

				for targetPos < len(targetRunes) && inputPos < len(inputRunes) {
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

				// Now targetPos is where we are in the target
				if targetPos < len(targetRunes) && targetRunes[targetPos] == '\n' {
					// Just add newline - rendering will skip the leading whitespace automatically
					m.currentInput += "\n"
					// Count the newline as correct
					m.correctChars++
					return m, nil
				}
			}

		case tea.KeyCtrlR:
			if !m.showingConfig && m.targetText != "" {
				// Reload with a new function
				funcText, filePath, err := loadRandomFunction(m.config)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.targetText = funcText
				m.currentFile = filePath
				m.currentInput = ""
				m.started = false
				m.finished = false
				m.correctChars = 0
				m.incorrectChars = 0
				m.errorPositions = make(map[int]bool)
				return m, nil
			}

		case tea.KeyTab, tea.KeyShiftTab:
			if m.showingConfig {
				// Navigate between config inputs
				if msg.Type == tea.KeyTab {
					m.focusIndex++
					if m.focusIndex >= len(m.configInputs) {
						m.focusIndex = 0
					}
				} else {
					m.focusIndex--
					if m.focusIndex < 0 {
						m.focusIndex = len(m.configInputs) - 1
					}
				}

				for i := range m.configInputs {
					if i == m.focusIndex {
						m.configInputs[i].Focus()
					} else {
						m.configInputs[i].Blur()
					}
				}
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Block all input when finished (except the special keys handled above)
	if m.finished {
		return m, nil
	}

	if m.showingConfig {
		m.configInputs[m.focusIndex], cmd = m.configInputs[m.focusIndex].Update(msg)
		return m, cmd
	}

	if !m.finished && m.targetText != "" {
		// Handle typing manually instead of using textInput
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyBackspace:
				if len(m.currentInput) > 0 {
					m.currentInput = m.currentInput[:len(m.currentInput)-1]
				}
			case tea.KeySpace:
				// Get current position in target
				pos := getCurrentPosition(m.currentInput, m.targetText)
				// Check if this character is correct before adding
				if isCharacterCorrect(m.currentInput, m.targetText, " ") {
					m.correctChars++
				} else {
					m.incorrectChars++
					m.errorPositions[pos] = true
				}
				m.currentInput += " "
			case tea.KeyRunes:
				// Get current position in target
				pos := getCurrentPosition(m.currentInput, m.targetText)
				// Check if this character is correct before adding
				char := string(msg.Runes)
				if isCharacterCorrect(m.currentInput, m.targetText, char) {
					m.correctChars++
				} else {
					m.incorrectChars++
					m.errorPositions[pos] = true
				}
				m.currentInput += char
			}
		}

		// Start timer on first keypress
		if !m.started && len(m.currentInput) > 0 {
			m.started = true
			m.startTime = time.Now()
			// Always start ticking to update elapsed time
			return m, tickCmd()
		}

		// Check if finished (compare without leading whitespace)
		if normalizeText(m.currentInput) == normalizeText(m.targetText) {
			m.finished = true
			m.endTime = time.Now()
		}

		return m, cmd
	}

	return m, nil
}

