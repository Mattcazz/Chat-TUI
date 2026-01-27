package user

type Service struct {
	UserStore UserStore
}

func NewService(us UserStore) *Service {
	return &Service{
		UserStore: us,
	}
}
