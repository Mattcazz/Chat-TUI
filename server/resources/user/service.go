package user

import (
	"context"
	"errors"

	"github.com/Mattcazz/Chat-TUI/server/resources/middleware"
)

type Service struct {
	userRepo      UserRepository
	contactRepo   ContactRepository
	challengeRepo ChallengeRepository
}

func NewService(userRepo UserRepository, contactRepo ContactRepository, challengeRepo ChallengeRepository) *Service {
	return &Service{
		userRepo:      userRepo,
		contactRepo:   contactRepo,
		challengeRepo: challengeRepo,
	}
}

func (s *Service) CreateOrLoginUser(ctx context.Context, user *User) error {
	return s.userRepo.CreateOrLoginUser(ctx, user)
}

func (s *Service) GenerateChallenge(ctx context.Context, publicKey string) (string, error) {
	_, err := s.userRepo.GetUserByPublicKey(ctx, publicKey)

	if err != nil {
		return "", err
	}

	nonce := generateRandomString(32)

	s.challengeRepo.CreateChallenge(ctx, publicKey, nonce)

	return nonce, nil
}

func (s *Service) VerifyAndLogin(ctx context.Context, publicKey, signature string) (string, error) {

	user, err := s.userRepo.GetUserByPublicKey(ctx, publicKey)

	if err != nil {
		return "", errors.New("User does not exist")
	}

	nonce, err := s.challengeRepo.GetNonceByPublicKey(ctx, publicKey)

	if err != nil {
		return "", errors.New("Challenge not requested")
	}

	if !isValidSshSignature(nonce, publicKey, signature) {
		return "", errors.New("Invalid signature")
	}

	return middleware.CreateJWT(nil, user.ID)
}

func (s *Service) GetContacts(ctx context.Context, userId int64) ([]*Contact, error) {
	return s.contactRepo.GetContactsByUserID(ctx, userId)
}
