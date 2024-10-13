package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	url2 "net/url"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/muesli/reflow/wordwrap"
)

func main() {
	m := NewModel()

	// NewProgram with initial model and program options
	p := tea.NewProgram(m, tea.WithAltScreen())

	// Run
	_, err := p.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

// Model: app state
type Model struct {
	title     string
	textinput textinput.Model
	terms     Terms
	err       error
	width     int
	height    int
	ready     bool
}

// NewModel: initial model
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter search term"
	ti.Focus()

	return Model{
		title:     "Urban Dictionary CLI",
		textinput: ti,
	}
}

// Init: event loop
func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)
}

// Update: handle messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			v := m.textinput.Value()
			return m, handleQuerySearch(v)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case TermsResponseMsg:
		if msg.Err != nil {
			m.err = msg.Err
		}

		m.terms = msg.Terms
		return m, nil
	}

	m.textinput, cmd = m.textinput.Update(msg)

	return m, cmd
}

// View: return a string based on the state of our model
func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	s := m.textinput.View() + "\n\n"

	if len(m.terms.List) > 0 {
		definition := wordwrap.String(m.terms.List[0].Definition, int(m.width))
		definitionStyled := color.New(color.FgBlue, color.Bold).Sprint(definition)

		example := wordwrap.String(m.terms.List[0].Example, int(m.width))
		exampleStyled := color.New(color.Italic).Sprint(example)

		thumbsUpStyled := color.New(color.FgGreen).Sprintf("👍 %d", m.terms.List[0].ThumbsUp)
		thumbsDownStyled := color.New(color.FgRed).Sprintf("👎 %d", m.terms.List[0].ThumbsDown)

		wordStyled := color.New(color.FgYellow, color.Bold, color.Underline).Sprintf(m.terms.List[0].Word)

		s += fmt.Sprintf("%s\n\n", wordStyled)
		s += fmt.Sprintf("%s\n\n", definitionStyled)
		s += fmt.Sprintf("%s\n\n", exampleStyled)
		s += fmt.Sprintf("%s\t%s\n\n", thumbsUpStyled, thumbsDownStyled)
	}

	return s
}

// Terms struct for parsing API response
type Terms struct {
	List []struct {
		Definition  string    `json:"definition"`
		Permalink   string    `json:"permalink"`
		ThumbsUp    int       `json:"thumbs_up"`
		Author      string    `json:"author"`
		Word        string    `json:"word"`
		Defid       int       `json:"defid"`
		CurrentVote string    `json:"current_vote"`
		WrittenOn   time.Time `json:"written_on"`
		Example     string    `json:"example"`
		ThumbsDown  int       `json:"thumbs_down"`
	} `json:"list"`
}

// Cmd for querying Urban Dictionary API
func handleQuerySearch(q string) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("https://api.urbandictionary.com/v0/define?term=%s", url2.QueryEscape(q))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return TermsResponseMsg{
				Err: err,
			}
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return TermsResponseMsg{
				Err: err,
			}
		}

		defer res.Body.Close()

		var terms Terms
		err = json.NewDecoder(res.Body).Decode(&terms)
		if err != nil {
			return TermsResponseMsg{
				Err: err,
			}
		}

		return TermsResponseMsg{
			Terms: terms,
		}
	}
}

// Msg to handle the API response
type TermsResponseMsg struct {
	Terms Terms
	Err   error
}
