package pkg

type MsgType string

const (
	MsgTypeText MsgType = "text"
	MsgTypeFile MsgType = "file"

	// Here we can add more message types in the future, such as images, videos, etc.
)

const MsgPerPage int64 = 20
