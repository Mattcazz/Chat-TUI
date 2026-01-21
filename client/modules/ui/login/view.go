package login

import (
	"fmt"
)

func (m Model) View() string {
	return fmt.Sprintf(
		"Login\n\n%s",
		m.text_input.View(),
	) + "\n"
}
