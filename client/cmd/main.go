package main

import (
	main_model "clit_client/internal/app"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(main_model.New(), tea.WithAltScreen())

	_, err := p.Run();
	if err != nil {
		log.Fatal(err)
	}
}
