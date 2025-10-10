package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	expressionInput textinput.Model
	subjectInput    textarea.Model

	expression string
	subject    string
}

func initialModel() model {
	m := model{
		expressionInput: textinput.New(),
		subjectInput:    textarea.New(),
	}

	m.expressionInput.Focus()

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 2)
	m.expressionInput, cmds[0] = m.expressionInput.Update(msg)
	m.subjectInput, cmds[1] = m.subjectInput.Update(msg)

	return m, tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	return m.updateInputs(msg)
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(m.expressionInput.View())
	b.WriteRune('\n')
	b.WriteString(m.subjectInput.View())

	return b.String()
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
