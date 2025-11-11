package expression

import (
	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/vitor-mariano/regex-tui/internal/styles"
	"github.com/vitor-mariano/regex-tui/pkg/regex/pcre"
	"github.com/vitor-mariano/regex-tui/pkg/regex/re2"
)

type Model struct {
	input textinput.Model
	width int
	pcre  bool
}

func newValidate(usePCRE bool) func(s string) error {
	if usePCRE {
		return func(s string) error {
			if s == "" {
				return nil
			}

			_, err := pcre.New(s)
			return err
		}
	}

	return func(s string) error {
		if s == "" {
			return nil
		}

		_, err := re2.New(s)
		return err
	}
}

func New(initialValue string) *Model {
	m := textinput.New()
	m.SetValue(initialValue)
	m.SetVirtualCursor(true)
	m.SetStyles(textinput.Styles{
		Cursor: textinput.CursorStyle{
			Color: styles.PrimaryColor,
			Blink: true,
		},
	})
	m.Prompt = ""
	m.Placeholder = "Expression"

	model := &Model{input: m}
	model.input.Validate = newValidate(model.pcre)

	return model
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	return cmd
}

func (m *Model) View() string {
	s := &styles.InputContainerStyle
	if m.input.Err != nil {
		s = &styles.ErrorInputContainerStyle
	} else if m.input.Focused() {
		s = &styles.FocusedInputContainerStyle
	}

	return s.Width(m.width).Render(m.input.View())
}

func (m *Model) SetWidth(width int) {
	m.width = width
	m.input.SetWidth(width)
}

func (m *Model) GetInput() *textinput.Model {
	return &m.input
}

func (m *Model) SetPCRE(enabled bool) {
	m.pcre = enabled
	m.input.Validate = newValidate(m.pcre)

	m.input.SetValue(m.input.Value())
}
