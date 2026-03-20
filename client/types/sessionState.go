package types

type SessionState int
const (
	LoginView SessionState = iota
	ChatView
)

func (s SessionState) String() string {
	switch s {
	case LoginView:
		return "LoginView"
	case ChatView:
		return "ChatView"
	default:
		return "UNIMPLEMENTED"
	}
}
