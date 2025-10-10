package regexview

import (
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	matchStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("205")).
		Bold(true)
)

type Model struct {
	expression *regexp.Regexp
	value      string
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	if m.expression == nil {
		return m.value
	}

	var b strings.Builder
	lastIndex := 0

	matches := m.expression.FindAllStringIndex(m.value, -1)
	for _, match := range matches {
		b.WriteString(m.value[lastIndex:match[0]])
		b.WriteString(matchStyle.Render(m.value[match[0]:match[1]]))
		lastIndex = match[1]
	}

	b.WriteString(m.value[lastIndex:])

	return b.String()
}

func (m *Model) SetExpression(expression *regexp.Regexp) {
	m.expression = expression
}

func (m *Model) SetExpressionString(expression string) error {
	expr, err := regexp.Compile(expression)

	m.SetExpression(expr)

	return err
}

func (m *Model) SetValue(value string) {
	m.value = value
}
