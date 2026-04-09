package types

import "time"

type Message struct {
	Author string
	Message string
	Timestamp time.Time
}
