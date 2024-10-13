package main

import (
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := NewModel()

	// NewProgram with initial model and program options
	p := tea.NewProgram(m)

	// Run
	_, err := p.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

// Model: app state
type Model struct {
	title string

	textinput textinput.Model
}

// NewModel: initial model
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter search term"
	ti.Focus()

	return Model{
		title:     "hello world",
		textinput: ti,
	}
}

// Init: event loop
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update: handle mssgs
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)

	return m, cmd
}

// View: return a string based on the state of our model
func (m Model) View() string {
	s := m.textinput.View()
	return s
}
