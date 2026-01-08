package subject

import (
	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/vitor-mariano/regex-tui/internal/styles"
	"github.com/vitor-mariano/regex-tui/pkg/components/regexview"
)

type Model struct {
	input         textarea.Model
	view          *regexview.Model
	cursor        cursor.Model
	width, height int
}

func New(initialValue, initialExpression string) *Model {
	m := textarea.New()
	m.SetValue(initialValue)
	m.SetVirtualCursor(false)
	m.SetStyles(textarea.Styles{
		Cursor: textarea.CursorStyle{
			Color: styles.PrimaryColor,
			Blink: true,
		},
		Focused: textarea.StyleState{
			CursorLine: lipgloss.NewStyle().UnsetBackground(),
		},
	})
	m.Prompt = ""
	m.ShowLineNumbers = false
	m.Focus()

	sv := regexview.New(0, 0)
	sv.SetExpression(initialExpression)
	sv.SetValue(initialValue)

	c := cursor.New()
	c.Style = lipgloss.NewStyle().Foreground(styles.PrimaryColor)

	return &Model{input: m, view: sv, cursor: c}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Focus() tea.Cmd {
	return tea.Batch(m.input.Focus(), m.cursor.Focus())
}

func (m *Model) Blur() {
	m.input.Blur()
	m.cursor.Blur()
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0, 2)

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.view.SetValue(m.input.Value())
	cmds = append(cmds, cmd)

	m.cursor, cmd = m.cursor.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Model) View() string {
	s := &styles.InputContainerStyle
	if m.input.Err != nil {
		s = &styles.ErrorInputContainerStyle
	} else if m.input.Focused() {
		s = &styles.FocusedInputContainerStyle
	}

	v := m.view.View()

	m.cursor.SetChar(" ")
	c := m.cursor.View()
	ic := m.input.Cursor()

	layers := make([]*lipgloss.Layer, 0, 2)
	layers = append(layers, lipgloss.NewLayer(v))
	if ic != nil && !m.cursor.IsBlinked {
		layers = append(layers, lipgloss.NewLayer(c).X(ic.X).Y(ic.Y))
	}

	return s.Width(m.width).Render(lipgloss.NewCanvas(layers...).Render())
}

func (m *Model) SetSize(width, height int) {
	const subjectHSpacing = 4

	m.width = width
	m.height = height
	m.input.SetWidth(width - subjectHSpacing)
	m.input.SetHeight(height)

	m.view.SetSize(width-subjectHSpacing-1, height)
}

func (m *Model) GetInput() *textarea.Model {
	return &m.input
}

func (m *Model) GetView() *regexview.Model {
	return m.view
}

func (m *Model) SetExpression(expression string) error {
	return m.view.SetExpression(expression)
}

func (m *Model) SetValue(value string) {
	m.input.SetValue(value)
	m.view.SetValue(value)
}
