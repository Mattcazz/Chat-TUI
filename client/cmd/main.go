package main

import (
	main_model "github.com/Mattcazz/Chat-TUI/client/internal/app"
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
