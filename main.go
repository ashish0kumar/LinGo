package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func main() {
	p := tea.NewProgram()

	_, err := p.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

