package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Mattcazz/Chat-TUI/client/internal/commands"
	"github.com/Mattcazz/Chat-TUI/client/modules/ui/chat"
	"github.com/Mattcazz/Chat-TUI/client/modules/ui/login"
	"github.com/Mattcazz/Chat-TUI/client/styles"
	"github.com/Mattcazz/Chat-TUI/client/types"
	"github.com/Mattcazz/Chat-TUI/pkg"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/ssh"
)

type App struct {
	state types.SessionState
	login_model tea.Model
	chat_model tea.Model

	client *http.Client

	username string
	err error
	width int
	height int
}

func getSSHPrivateKey() []byte {
	keyPath := os.Getenv("HOME") + "/.ssh/id_ed25519"
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		panic("Could not read private key: " + err.Error())
	}

	return keyBytes
}

func createSignature(nonce string, sk []byte) string {
	// 2. Parse Private Key
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		panic(err)
	}

	// 3. Sign the Nonce
	sig, err := signer.Sign(rand.Reader, []byte(nonce))
	if err != nil {
		panic(err)
	}

	// 4. Encode to Base64 (This is what you paste into Curl)
	sigBytes := ssh.Marshal(sig) // Important: Marshal to wire format first!
	b64Sig := base64.StdEncoding.EncodeToString(sigBytes)

	// 5. Get Public Key String (for the request)
	pubKey := ssh.MarshalAuthorizedKey(signer.PublicKey())

	fmt.Println("\n--- COPY THESE FOR CURL ---")
	fmt.Printf("Public Key: %s", pubKey) // Contains newline
	fmt.Printf("Signature:  %s\n", b64Sig)
}

func New() App {
	app_client := &http.Client{
		Timeout: time.Second * 10,
	}

	return App{
		state: types.LoginView,
		login_model: login.New(),
		chat_model: chat.New(),
		client: app_client,
		err: nil,
	}
}

func (a App) Init() tea.Cmd {
	return nil
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case commands.ChangeStateMsg:
		m.state = msg.State
		return m, nil
	case commands.LogInMsg:
		m.username = msg.Username
		// Check that log in was successful eventually, but
		// TODO get pub key
		login_req := pkg.LoginRequest{
			PublicKey: "fuck",
		}
		body, err := json.Marshal(login_req)
		if err != nil {
			// TODO
		}
		req, err := http.NewRequest("POST", "localhost", bytes.NewBuffer(body))

		req.Header.Add("Authorization", "JWT") // TODO
		req.Header.Add("Content-Type", "application/json")

		resp, err := m.client.Do(req)
		if err != nil {
			// TODO
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			// Logged in successfully
		case http.StatusAccepted:
			// Nonce coming
			var challenge_resp pkg.ChallengeResponse
			json.NewDecoder(resp.Body).Decode(&challenge_resp)
			// TODO handle that shit
		case http.StatusTemporaryRedirect:
			// Oops, we need to register
			// TODO /register with username
		}


		var login_resp pkg.LoginResponse
		json.NewDecoder(resp.Body).Decode(&login_resp)

		m.chat_model, _ = m.chat_model.Update(msg)
		return m, commands.NewChangeStateCmd(types.ChatView)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		msg.Width -= lipgloss.NormalBorder().GetLeftSize()
		msg.Width -= lipgloss.NormalBorder().GetRightSize()
		msg.Height -= lipgloss.NormalBorder().GetTopSize()
		msg.Height -= lipgloss.NormalBorder().GetBottomSize()

		m.login_model, _ = m.login_model.Update(msg)
		m.chat_model, _ = m.chat_model.Update(msg)

		return m, nil
	}

	switch m.state {
	case types.LoginView:
		m.login_model, cmd = m.login_model.Update(msg)
		return m, cmd
	case types.ChatView:
		m.chat_model, cmd = m.chat_model.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m App) View() string {
	var model_content string

	switch m.state {
	case types.LoginView:
		model_content = m.login_model.View()
	case types.ChatView:
		model_content = m.chat_model.View()
	default:
		model_content = "Ya broke it"
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		styles.Default.Border.Render(model_content),
	)
}


