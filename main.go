package main

import (
	"log"

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
}

// NewModel: initial model
func NewModel() Model {
	return Model{
		title: "hello world",
	}
}

// Init: event loop
func (m Model) Init() tea.Cmd {
	return nil
}

// Update: handle mssgs
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View: return a string based on the state of our model
func (m Model) View() string {
	return m.title
}
