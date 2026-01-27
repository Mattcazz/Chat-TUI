package file

type Service struct {
	FileStore FileStore
}

func NewService(fs FileStore) *Service {
	return &Service{
		FileStore: fs,
	}
}
