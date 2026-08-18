package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anotherhadi/ilovetui/bubbles"
)

type model struct {
	input   textinput.Model
	spinner spinner.Model
}

func newModel() model {
	ti := bubbles.NewTextInput()
	ti.Placeholder = "type something"
	ti.Focus()
	return model{input: ti, spinner: bubbles.NewSpinner()}
}

func (m model) Init() tea.Cmd { return m.spinner.Tick }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok && k.String() == "ctrl+c" {
		return m, tea.Quit
	}
	var inputCmd, spinnerCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	m.spinner, spinnerCmd = m.spinner.Update(msg)
	return m, tea.Batch(inputCmd, spinnerCmd)
}

func (m model) View() tea.View {
	return tea.NewView(lipgloss.JoinHorizontal(lipgloss.Center, m.spinner.View(), " ", m.input.View()))
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
