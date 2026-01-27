package msg

type Service struct {
	MsgStore MsgStore
}

func NewService(ms MsgStore) *Service {
	return &Service{
		MsgStore: ms,
	}
}
