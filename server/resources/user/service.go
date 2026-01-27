package user

type Service struct {
	userRepo    UserRepository
	contactRepo ContactRepository
}

func NewService(userRepo UserRepository, contactRepo ContactRepository) *Service {
	return &Service{
		userRepo:    userRepo,
		contactRepo: contactRepo,
	}
}
