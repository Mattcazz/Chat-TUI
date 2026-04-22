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
	logger.Log.Printf("[NORMAL] Attempting log in")
	if m.signature == nil {
		logger.Log.Printf("[NORMAL] Requesting Challenge from server")
		var err error
		if m.nonce == nil {
			logger.Log.Printf("[NORMAL] nonce is nil, requesting challenge")
			m.nonce, err = m.client.RequestChallenge(m.pk)
			if err != nil {
				if errors.Is(err, errors.ErrUnsupported) {
					// 400 from server, would be nice to have a better error at some point
					logger.Log.Printf("[NORMAL] Couldn't get challenge from server, duplicate something or other, check server logs")
					m.errorMsg = "Something went wrong, please check the server logs"
					return m, nil
				}
				logger.Log.Printf("[NORMAL] Couldn't get challenge from server, user does not exist, opening register state")
				m.state = types.NeedsUsername
				m.usernameInput.Reset()

				return m, nil
			}
		}
		logger.Log.Printf("[NORMAL] Got nonce: %s", m.nonce)

		logger.Log.Printf("[NORMAL] Creating signature")
		m.signature, err = createSignature(string(m.nonce), m.sk, nil)
		if _, ok := err.(*ssh.PassphraseMissingError); ok {
			logger.Log.Printf("[NORMAL] Couldn't create signature, private SSH key is encrypted, asking for password")
			m.state = types.NeedsSSHPassword
			return m, nil
		}
		logger.Log.Printf("[NORMAL] Signature created: %s", m.signature)
	}
	logger.Log.Printf("[NORMAL] Attempting to log in")
	_, err := m.client.Login(m.pk, m.signature)
	if err != nil {
		logger.Log.Panicf("Error while trying to log in: %s", err.Error())
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
		logger.Log.Printf("Loading ssh keys")
		m.pk, m.sk = getSSHKeys() // TODO Let user pick at some point
		return m, commands.NewDoLoginCmd()
	case commands.DoLogInMsg:
		logger.Log.Printf("Attempting log in")

		return m.doLoginCmd()
	}

	switch m.state {
	case types.Normal:
		return m, nil
	case types.NeedsUsername:
		// Show username input and intercept enter key
		m.usernameInput.Focus()

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				logger.Log.Printf("[USERNAME] Attempting to register with username: %s", m.usernameInput.Value())
				m.client.Register(m.pk, m.usernameInput.Value())

				logger.Log.Printf("[USERNAME] Successfully registered, returning to NORMAL mode")
				m.signature = nil
				m.state = types.Normal

				return m, commands.NewDoLoginCmd()
			}
		}

		m.usernameInput, cmd = m.usernameInput.Update(msg)
		return m, cmd
	case types.NeedsSSHPassword:
		// Show password input and intercept enter key
		m.passwordInput.Focus()

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				// use password idfk man
				var err error
				if m.nonce == nil {
					logger.Log.Printf("[SSH] nonce is nil, requesting challenge")
					m.nonce, err = m.client.RequestChallenge(m.pk)
					if err != nil {
						m.state = types.NeedsUsername // should never happen, since this should trip before ssh password
						m.usernameInput.Reset() // maybe panic? maybe fatal?
						logger.Log.Panicf("Error while requesting challenge: %s", err.Error())
					}
				}

				m.signature, err = createSignature(string(m.nonce), m.sk, []byte(m.passwordInput.Value()))
				if err != nil {
					if errors.Is(err, x509.IncorrectPasswordError) {
						// TODO fucking explode idfk, you shouldn't be allowed to not know your password idk
						// Probably a state for wrong password
						m.errorMsg = "Incorrect password. Try again."
						// log.Fatal("imagine not knowing your password holy shit")
						m.passwordInput.Reset()
						m.passwordInput, cmd = m.passwordInput.Update(msg)
						return m, cmd
					} else {
						logger.Log.Panicf("Error while creating signature: %s", err.Error())
					}
				}

				m.state = types.Normal
				return m, commands.NewDoLoginCmd()
			}
		}

		m.passwordInput, cmd = m.passwordInput.Update(msg)
		return m, cmd
	}

	return m, cmd
}
