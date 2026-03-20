package types

type LoginModelState int
const (
	Normal LoginModelState = iota
	NeedsUsername
	NeedsSSHPassword
)
