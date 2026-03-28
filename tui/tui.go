package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	stateInputName sessionState = iota
	stateSelectDB
	stateSelectPipenv
	stateSelectVenv
)

type Model struct {
	State       sessionState
	TextInput   textinput.Model
	ProjectName string
	Cursor      int
	Choices     []string
	Selected    string
	UsePipenv   bool
	SetupVenv   bool
	Quitting    bool
}

func InitialModel() Model {
	return InitialModelWithName("")
}

func InitialModelWithName(name string) Model {
	ti := textinput.New()
	ti.Placeholder = "my-fastapi-app"
	ti.CharLimit = 32
	ti.Width = 20

	state := stateInputName
	if name != "" {
		state = stateSelectDB
	} else {
		ti.Focus()
	}

	return Model{
		State:       state,
		TextInput:   ti,
		ProjectName: name,
		Choices:     []string{"PostgreSQL (SQLAlchemy)", "MongoDB (PyMongo)"},
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, tea.Quit

		case "enter":
			if m.State == stateInputName {
				m.ProjectName = m.TextInput.Value()
				if m.ProjectName == "" {
					m.ProjectName = "my-fastapi-app"
				}
				m.State = stateSelectDB
				m.Cursor = 0
				return m, nil
			} else if m.State == stateSelectDB {
				m.Selected = m.Choices[m.Cursor]
				m.State = stateSelectPipenv
				m.Cursor = 0
				return m, nil
			} else if m.State == stateSelectPipenv {
				m.UsePipenv = m.Cursor == 0
				m.State = stateSelectVenv
				m.Cursor = 0
				return m, nil
			} else if m.State == stateSelectVenv {
				m.SetupVenv = m.Cursor == 0
				return m, tea.Quit
			}

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.State == stateSelectDB && m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			} else if (m.State == stateSelectPipenv || m.State == stateSelectVenv) && m.Cursor < 1 {
				m.Cursor++
			}
		}
	}

	if m.State == stateInputName {
		m.TextInput, cmd = m.TextInput.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.Quitting {
		return "Exiting...\n"
	}

	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)

	if m.State == stateInputName {
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			headerStyle.Render("Step 1: Project Name"),
			"What is your project called?",
			m.TextInput.View(),
		) + "\n\n(press enter to continue)\n"
	}

	if m.State == stateSelectDB {
		s := headerStyle.Render("Step 2: Database Selection") + "\n\n"
		s += fmt.Sprintf("Project: %s\n\nChoose a DB:\n", m.ProjectName)
		for i, choice := range m.Choices {
			cursor := " "
			if m.Cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
		return s + "\n(j/k to move, enter to select)\n"
	}

	// Pipenv selection
	if m.State == stateSelectPipenv {
		s := headerStyle.Render("Step 3: Package Manager") + "\n\n"
		s += fmt.Sprintf("Project: %s | DB: %s\n\nUse pipenv?\n", m.ProjectName, m.Selected)
		pipenvChoices := []string{"Yes (pipenv)", "No (requirements.txt)"}
		for i, choice := range pipenvChoices {
			cursor := " "
			if m.Cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
		return s + "\n(j/k to move, enter to select)\n"
	}

	// Venv setup selection
	pkgManager := "requirements.txt"
	if m.UsePipenv {
		pkgManager = "pipenv"
	}
	s := headerStyle.Render("Step 4: Virtual Environment") + "\n\n"
	s += fmt.Sprintf("Project: %s | DB: %s | PM: %s\n\nInstall pip packages now?\n", m.ProjectName, m.Selected, pkgManager)
	venvChoices := []string{"Yes", "No"}
	for i, choice := range venvChoices {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	return s + "\n(j/k to move, enter to select)\n"
}