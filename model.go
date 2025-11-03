package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

type model struct {
	targetText    string
	currentInput  string
	textInput     textinput.Model
	startTime     time.Time
	endTime       time.Time
	started       bool
	finished      bool
	currentFile   string
	width         int
	height        int
	err           error
	config        config
	showingConfig bool
	configInputs  []textinput.Model
	focusIndex    int
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Start typing..."
	ti.Focus()
	ti.CharLimit = 10000
	ti.Width = 80

	cfg := loadConfig()

	// Create config form inputs
	inputs := make([]textinput.Model, 4)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Folder path"
	inputs[0].SetValue(cfg.FolderPath)
	inputs[0].Width = 50

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Minimum lines"
	inputs[1].SetValue(fmt.Sprintf("%d", cfg.MinLines))
	inputs[1].Width = 20

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Maximum lines"
	inputs[2].SetValue(fmt.Sprintf("%d", cfg.MaxLines))
	inputs[2].Width = 20

	inputs[3] = textinput.New()
	inputs[3].Placeholder = "Max time (0 = no limit)"
	inputs[3].SetValue(fmt.Sprintf("%d", cfg.MaxTimeLimit))
	inputs[3].Width = 20

	return model{
		textInput:    ti,
		config:       cfg,
		width:        120,
		height:       24,
		configInputs: inputs,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

