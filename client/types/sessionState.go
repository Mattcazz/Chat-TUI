package types

type SessionState int
const (
	LoginView SessionState = iota
	ChatView
)

