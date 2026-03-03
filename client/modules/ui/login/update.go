package login

import (
	"crypto/x509"
	"errors"
	"log"

	"github.com/Mattcazz/Chat-TUI/client/internal/commands"
	"github.com/Mattcazz/Chat-TUI/client/internal/logger"
	"github.com/Mattcazz/Chat-TUI/client/types"
	"golang.org/x/crypto/ssh"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) doLoginCmd() (tea.Model, tea.Cmd) {
	if m.signature == nil {
		logger.Log.Printf("[NORMAL] Requesting Challenge from server...")
		var err error
		if m.nonce == nil {
			logger.Log.Printf("[NORMAL] nonce is nil, requesting challenge...")
			m.nonce, err = m.client.RequestChallenge(m.pk)
			if err != nil {
				logger.Log.Printf("[NORMAL] Couldn't get challenge from server, user does not exist, opening register state")
				m.state = types.NeedsUsername
				m.username_input.Reset()

				return m, nil // TODO make this automatically go to next part
			}
		}
		logger.Log.Printf("[NORMAL] Got nonce: %s", m.nonce)

		logger.Log.Printf("[NORMAL] Creating signature...")
		m.signature, err = createSignature(string(m.nonce), m.sk, nil)
		if _, ok := err.(*ssh.PassphraseMissingError); ok {
			logger.Log.Printf("[NORMAL] Couldn't create signature, private SSH key is encrypted, asking for password")
			m.state = types.NeedsSSHPassword
			return m, nil // TODO make this automatically go to next part
		}
		logger.Log.Printf("[NORMAL] Signature created: %s", m.signature)
	}
	logger.Log.Printf("[NORMAL] Attempting to log in...")
	_, err := m.client.Login(m.pk, m.signature) // TODO get the login response with username
	if err != nil {
		log.Panic(err.Error())
	}
	logger.Log.Printf("[NORMAL] Login succeeded, returning username: %s", "I don't know how to get this yet :D")

	// TODO get username somehow
	return m, commands.NewLogInCmd("test_username")
}

func (m Model) Init() tea.Cmd {
	return commands.NewLoadSSHKeysCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case commands.LoadSSHKeysMsg:
		logger.Log.Printf("Loading ssh keys...")
		m.pk, m.sk = getSSHKeys() // TODO Let user pick at some point
		return m, commands.NewDoLoginCmd()
	case commands.DoLogInMsg:
		logger.Log.Printf("Attempting log in...")

		return m.doLoginCmd()
	}

	switch m.state {
	case types.Normal:
		return m, nil
	case types.NeedsUsername:
		// Show username input and intercept enter key
		m.username_input.Focus()

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				logger.Log.Printf("[USERNAME] Attempting to register with username: %s", m.username_input.Value())
				m.client.Register(m.pk, m.username_input.Value())

				logger.Log.Printf("[USERNAME] Successfully registered, returning to NORMAL mode...")
				m.signature = nil
				m.state = types.Normal

				return m, commands.NewDoLoginCmd()
			}
		}

		m.username_input, cmd = m.username_input.Update(msg)
		return m, cmd
	case types.NeedsSSHPassword:
		// Show password input and intercept enter key
		m.password_input.Focus()

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				// use password idfk man
				var err error
				if m.nonce == nil {
					logger.Log.Printf("[SSH] nonce is nil, requesting challenge...")
					m.nonce, err = m.client.RequestChallenge(m.pk)
					if err != nil {
						m.state = types.NeedsUsername // should never happen, since this should trip before ssh password
						m.username_input.Reset() // maybe panic? maybe fatal?
						// TODO return
					}
				}

				m.signature, err = createSignature(string(m.nonce), m.sk, []byte(m.password_input.Value()))
				if err != nil {
					if errors.Is(err, x509.IncorrectPasswordError) {
						// TODO fucking explode idfk, you shouldn't be allowed to not know your password idk
						// Probably a state for wrong password
						log.Fatal("imagine not knowing your password holy shit")
					} else {
						log.Panic(err.Error())
					}
				}

				m.state = types.Normal
				return m, commands.NewDoLoginCmd()
			}
		}

		m.password_input, cmd = m.password_input.Update(msg)
		return m, cmd
	}

	return m, cmd
}
