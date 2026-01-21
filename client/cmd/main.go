package main

import (
	login_model "clit_client/modules/ui/login"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(login_model.New())

	_, err := p.Run();
	if err != nil {
		log.Fatal(err)
	}
}
