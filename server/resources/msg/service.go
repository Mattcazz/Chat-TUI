package msg

type Service struct {
	msgRepo MsgRepository
}

func NewService(msgRepo MsgRepository) *Service {
	return &Service{
		msgRepo: msgRepo,
	}
}
