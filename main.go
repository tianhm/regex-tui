package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/textarea"
	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/vitor-mariano/regex-tui/pkg/regexview"
)

type keyMap struct {
	Exit        key.Binding
	SwitchInput key.Binding
}

var keys = keyMap{
	Exit: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
	),
	SwitchInput: key.NewBinding(
		key.WithKeys("tab", "shift+tab"),
	),
}

type inputType int

const (
	inputTypeExpression inputType = iota
	inputTypeSubject
)

const (
	initialExpression = "[A-Z]\\w+"
	initialSubject    = "Hello World!"
)

var (
	primaryColor = lipgloss.Color("12")
	mutedColor   = lipgloss.Color("240")
	lightColor   = lipgloss.Color("15")
	errorColor   = lipgloss.Color("9")

	titleStyle = lipgloss.NewStyle().
			Background(primaryColor).
			Bold(true).
			Foreground(lightColor).
			Padding(0, 1).
			MarginLeft(1).
			MarginTop(1)

	inputContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Padding(0, 1)
	focusedInputContainerStyle = inputContainerStyle.
					BorderForeground(primaryColor)
	errorInputContainerStyle = inputContainerStyle.
					BorderForeground(errorColor)

	helpContainerStyle = lipgloss.NewStyle().
				MarginLeft(1)
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor)
	helpCommandStyle = helpStyle.
				Bold(true)
)

type model struct {
	expressionInput textinput.Model
	subjectInput    *textarea.Model
	subjectView     regexview.Model

	focusedInputType inputType
	expression       string
	subject          string
	width            int
	height           int
	initCmds         []tea.Cmd
}

func initialModel() model {
	m := model{
		expressionInput: textinput.New(),
		subjectInput:    textarea.New(),
		subjectView:     regexview.New(0, 0),
	}

	m.expressionInput.SetValue(initialExpression)
	m.expressionInput.Prompt = ""
	m.expressionInput.Placeholder = "Expression"
	m.expressionInput.SetVirtualCursor(true)
	m.expressionInput.SetStyles(textinput.Styles{
		Cursor: textinput.CursorStyle{
			Color: primaryColor,
			Blink: true,
		},
	})
	if cmd := m.expressionInput.Focus(); cmd != nil {
		m.initCmds = append(m.initCmds, cmd)
	}
	m.expressionInput.Validate = func(s string) error {
		_, err := regexp.Compile(s)
		return err
	}

	m.subjectInput.SetValue(initialSubject)
	m.subjectInput.Prompt = ""
	m.subjectInput.ShowLineNumbers = false
	m.subjectInput.SetVirtualCursor(true)
	m.subjectInput.SetStyles(textarea.Styles{
		Cursor: textarea.CursorStyle{
			Color: primaryColor,
			Blink: true,
		},
		Focused: textarea.StyleState{
			CursorLine: lipgloss.NewStyle().UnsetBackground(),
		},
	})

	m.subjectView.SetExpressionString(initialExpression)
	m.subjectView.SetValue(initialSubject)

	return m
}

func (m model) Init() tea.Cmd {
	cmds := append([]tea.Cmd{textinput.Blink, textarea.Blink}, m.initCmds...)
	return tea.Batch(cmds...)
}

func (m *model) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.focusedInputType == inputTypeSubject {
		sm, cmd := m.subjectInput.Update(msg)
		m.subjectInput = sm
		m.subjectView.SetValue(m.subjectInput.Value())

		return m, cmd
	}

	cmds := make([]tea.Cmd, 2)

	m.expressionInput, cmds[0] = m.expressionInput.Update(msg)

	err := m.subjectView.SetExpressionString(m.expressionInput.Value())
	if err == nil {
		m.subjectView, cmds[1] = m.subjectView.Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 2)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		const (
			expInputHSpacing    = 6
			subjInputHSpacing   = 5
			subjInputTopSpacing = 8
		)

		m.width = msg.Width
		m.height = msg.Height
		m.expressionInput.SetWidth(m.width - expInputHSpacing)
		m.subjectView.SetSize(m.width-subjInputHSpacing, m.height-subjInputTopSpacing)
		m.subjectInput.SetWidth(m.width - subjInputHSpacing)
		m.subjectInput.SetHeight(m.height - subjInputTopSpacing)

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.Exit):
			return m, tea.Quit

		case key.Matches(msg, keys.SwitchInput):
			var cmd tea.Cmd
			switch m.focusedInputType {
			case inputTypeExpression:
				m.focusedInputType = inputTypeSubject
				m.expressionInput.Blur()
				cmd = m.subjectInput.Focus()

			case inputTypeSubject:
				m.focusedInputType = inputTypeExpression
				m.subjectInput.Blur()
				cmd = m.expressionInput.Focus()
			}

			cmds = append(cmds, cmd)
		}
	}

	_, cmd := m.updateInputs(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Regex TUI"))
	b.WriteRune('\n')

	s := &inputContainerStyle
	if m.expressionInput.Err != nil {
		s = &errorInputContainerStyle
	} else if m.focusedInputType == inputTypeExpression {
		s = &focusedInputContainerStyle
	}
	b.WriteString(s.Render(m.expressionInput.View()))
	b.WriteRune('\n')

	if m.focusedInputType == inputTypeSubject {
		b.WriteString(focusedInputContainerStyle.Render(m.subjectInput.View()))
	} else {
		b.WriteString(inputContainerStyle.Render(m.subjectView.View()))
	}
	b.WriteRune('\n')

	h := helpCommandStyle.Render("esc/ctrl+c ") +
		helpStyle.Render("exit • ") +
		helpCommandStyle.Render("tab ") +
		helpStyle.Render("switch input")
	b.WriteString(helpContainerStyle.Render(h))

	return tea.NewView(b.String())
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
