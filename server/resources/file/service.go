package file

type Service struct {
	fileRepo FileRepository
}

func NewService(fr FileRepository) *Service {
	return &Service{
		fileRepo: fr,
	}
}
