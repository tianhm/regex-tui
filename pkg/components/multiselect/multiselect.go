package multiselect

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/vitor-mariano/regex-tui/internal/styles"
	"github.com/vitor-mariano/regex-tui/pkg/utils"
)

type Model struct {
	cursor           string
	cursorStyle      lipgloss.Style
	checkIcon        string
	checkIconStyle   lipgloss.Style
	uncheckIcon      string
	uncheckIconStyle lipgloss.Style

	items    []string
	selected *utils.Set[string]
	current  int
	onToggle func(item string, selected bool)
}

func New(items []string) *Model {
	return &Model{
		cursor: ">",
		cursorStyle: lipgloss.NewStyle().
			Foreground(styles.PrimaryColor).
			MarginRight(1),
		checkIcon: "✓",
		checkIconStyle: lipgloss.NewStyle().
			Foreground(styles.PrimaryColor).
			Bold(true).
			MarginRight(1),
		uncheckIcon: "•",
		uncheckIconStyle: lipgloss.NewStyle().
			Foreground(styles.MutedColor).
			MarginRight(1),

		items:    items,
		selected: utils.NewSet[string](),
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Up):
			m.current--
			if m.current < 0 {
				m.current = len(m.items) - 1
			}
		case key.Matches(msg, Keys.Down):
			m.current++
			if m.current >= len(m.items) {
				m.current = 0
			}
		case key.Matches(msg, Keys.Toggle):
			item := m.items[m.current]
			if m.selected.Contains(item) {
				m.selected.Remove(item)
			} else {
				m.selected.Add(item)
			}

			if m.onToggle != nil {
				m.onToggle(item, m.selected.Contains(item))
			}
		}
	}

	return cmd
}

func (m *Model) formattedItems() []string {
	items := make([]string, len(m.items))

	for i, item := range m.items {
		cursor := m.cursorStyle.Render(" ")
		if i == m.current {
			cursor = m.cursorStyle.Render(m.cursor)
		}

		selected := m.uncheckIconStyle.Render(m.uncheckIcon)
		if m.selected.Contains(item) {
			selected = m.checkIconStyle.Render(m.checkIcon)
		}

		items[i] = cursor + selected + item
	}

	return items
}

func (m *Model) View() string {
	return lipgloss.NewStyle().
		Render(lipgloss.JoinVertical(lipgloss.Left, m.formattedItems()...))
}

func (m *Model) SetItems(items []string) {
	m.items = items
}

func (m *Model) SetSelected(items ...string) {
	m.selected.Add(items...)

	if m.onToggle != nil {
		for _, item := range items {
			m.onToggle(item, m.selected.Contains(item))
		}
	}
}

func (m *Model) OnToggle(onToggle func(item string, selected bool)) {
	m.onToggle = onToggle
}
