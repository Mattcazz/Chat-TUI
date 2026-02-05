package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Mattcazz/Chat-TUI/server/resources/middleware"
	"github.com/Mattcazz/Chat-TUI/server/utils"
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

func (s *Service) CreateUser(ctx context.Context, publicKey, username string) (*User, error) {
	user, err := s.userRepo.GetUserByPublicKey(ctx, publicKey)

	if err == nil && user != nil {
		return nil, fmt.Errorf("User already exists")
	}

	u := &User{
		PublicKey: publicKey,
		Username:  username,
	}

	return s.userRepo.CreateUser(ctx, u)
}

func (s *Service) GenerateChallenge(ctx context.Context, publicKey string) (string, error) {
	user, err := s.userRepo.GetUserByPublicKey(ctx, publicKey)

	if err != nil {
		err = NewUserDoesNotExistError()
		return "", err
	}

	nonce := utils.RandomString(32)
	expires_at := time.Now().Add(5 * time.Minute)

	challenge := &Challenge{
		UserID:    user.ID,
		Nonce:     nonce,
		ExpiresAt: expires_at,
	}

	err = s.challengeRepo.CreateChallenge(ctx, challenge)

	if err != nil {
		return "", err
	}

	return nonce, nil
}

func (s *Service) VerifyAndLogin(ctx context.Context, publicKey, signature string) (string, error) {

	user, err := s.userRepo.GetUserByPublicKey(ctx, publicKey)

	if err != nil {
		err = NewUserDoesNotExistError()
		return "", err
	}

	challenge, err := s.challengeRepo.GetChallenge(ctx, user.ID)

	defer s.challengeRepo.DeleteChallenge(ctx, user.ID, challenge.Nonce)

	if err != nil {
		return "", fmt.Errorf("Challenge not created")
	}

	if challenge.ExpiresAt.Before(time.Now()) {
		return "", fmt.Errorf("Challenge expired")
	}

	if err := middleware.IsValidSshSignature(publicKey, challenge.Nonce, signature); err != nil {
		return "", err
	}

	return middleware.CreateJWT(user.ID)
}
