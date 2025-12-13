package screen

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/vitor-mariano/regex-tui/internal/components/expression"
	"github.com/vitor-mariano/regex-tui/internal/components/options"
	"github.com/vitor-mariano/regex-tui/internal/components/subject"
	"github.com/vitor-mariano/regex-tui/pkg/components/multiselect"
)

type inputType int

const (
	inputTypeExpression inputType = iota
	inputTypeSubject
)

type Config struct {
	InitialExpression string
	InitialSubject    string
	Global            bool
	Insensitive       bool
	Regexp2           bool
}

type model struct {
	expressionInput *expression.Model
	subjectInput    *subject.Model
	options         *options.Model
	help            help.Model

	focusedInputType inputType
	width, height    int
}

func New(config Config) model {
	si := subject.New(config.InitialSubject, config.InitialExpression)

	ei := expression.New(config.InitialExpression, si.GetView())
	ei.GetInput().Focus()

	d := options.New()
	d.OnToggle(func(item string, selected bool) {
		switch item {
		case options.GlobalOption:
			si.GetView().SetGlobal(selected)
		case options.InsensitiveOption:
			si.GetView().SetInsensitive(selected)
		case options.Regexp2Option:
			si.GetView().SetRegexp2(selected)
			ei.GetInput().Err = si.GetView().SetRegexp2(selected)
			// Force re-evaluation with the new engine.
			si.SetExpression(ei.GetInput().Value())
		}
	})

	var selectedOptions []string
	if config.Global {
		selectedOptions = append(selectedOptions, options.GlobalOption)
	}
	if config.Insensitive {
		selectedOptions = append(selectedOptions, options.InsensitiveOption)
	}
	if config.Regexp2 {
		selectedOptions = append(selectedOptions, options.Regexp2Option)
	}

	if len(selectedOptions) > 0 {
		d.SetSelected(selectedOptions...)
	}

	return model{
		expressionInput: ei,
		subjectInput:    si,
		options:         d,
		help:            help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return m.expressionInput.Init()
}

func (m *model) setSize(width, height int) {
	const subjectVSpacing = 8

	m.width = width
	m.height = height
	m.expressionInput.SetWidth(width)
	m.subjectInput.SetSize(width, height-subjectVSpacing)
}

func (m *model) updateScreen(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0, 2)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.SwitchInput):
			var cmd tea.Cmd
			switch m.focusedInputType {
			case inputTypeExpression:
				m.focusedInputType = inputTypeSubject
				m.expressionInput.GetInput().Blur()
				cmd = m.subjectInput.GetInput().Focus()

			case inputTypeSubject:
				m.focusedInputType = inputTypeExpression
				m.subjectInput.GetInput().Blur()
				cmd = m.expressionInput.GetInput().Focus()
			}

			cmds = append(cmds, cmd)

		case key.Matches(msg, keys.ToggleOptions):
			if !m.options.IsOpen() {
				m.options.Open()
			}
		}
	}

	if m.focusedInputType == inputTypeSubject {
		cmds = append(cmds, m.subjectInput.Update(msg))
	} else {
		cmds = append(cmds, m.expressionInput.Update(msg))
		m.subjectInput.SetExpression(m.expressionInput.GetInput().Value())
	}

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0, 2)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)

	case tea.KeyPressMsg:
		if key.Matches(msg, keys.Exit) {
			if m.options.IsOpen() {
				break
			}

			return m, tea.Quit
		}
	}

	if m.options.IsOpen() {
		cmds = append(cmds, m.options.Update(msg))
	} else {
		cmds = append(cmds, m.updateScreen(msg))
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	var helpKeyMap help.KeyMap = keys
	if m.options.IsOpen() {
		helpKeyMap = multiselect.Keys
	}

	baseLayer := lipgloss.NewLayer(lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		m.expressionInput.View(),
		m.subjectInput.View(),
		m.help.View(helpKeyMap),
	))

	layers := []*lipgloss.Layer{baseLayer}
	if m.options.IsOpen() {
		optionsLayer := lipgloss.NewLayer(m.options.View())
		optionsLayer.X((m.width - optionsLayer.GetWidth()) / 2)
		optionsLayer.Y((m.height - optionsLayer.GetHeight()) / 2)

		layers = append(layers, optionsLayer)
	}

	return tea.NewView(lipgloss.NewCanvas(layers...).Render())
}
