package user

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Mattcazz/Chat-TUI/pkg"
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

func (s *Service) DeleteUser(ctx context.Context, userID int64) error {
	return s.userRepo.DeleteUser(ctx, userID)
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
		log.Printf("Error retrieving challenge for user %d: %v", user.ID, err)
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

func (s *Service) GetContacts(ctx context.Context, userID int64) ([]*pkg.ContactDetails, error) {
	contacts, err := s.contactRepo.GetContactsByUserID(ctx, userID)

	// if the query returns nil, we want to show that there are no contacts from this user
	if err == nil && contacts == nil {
		return []*pkg.ContactDetails{}, nil
	}

	return contacts, err
}

func (s *Service) ContactRequest(ctx context.Context, fromUserID int64, toPk, nickname string) error {

	toUser, err := s.userRepo.GetUserByPublicKey(ctx, toPk)

	if err != nil {
		return fmt.Errorf("User with public key %s does not exist", toPk)
	}

	status := StatusPending

	contact, err := s.contactRepo.GetContactByPair(ctx, toUser.ID, fromUserID)

	if err == nil && contact != nil {
		// contact already exists, we accept the contact request (both ways)
		contact.Status = StatusAccept
		status = StatusAccept
		contact.UpdatedAt = time.Now()

		err = s.contactRepo.UpdateContact(ctx, contact)

		if err != nil {
			return err
		}
	}

	if nickname == "" {
		nickname = toUser.Username
	}

	c := &Contact{
		FromUserID: fromUserID,
		ToUserID:   toUser.ID,
		Nickname:   nickname,
		Status:     status,
		UpdatedAt:  time.Now(),
		Created_at: time.Now(),
	}

	return s.contactRepo.CreateContact(ctx, c)

}

func (s *Service) GetContactRequests(ctx context.Context, userID int64) ([]*pkg.ContactDetails, error) {
	contacts, err := s.contactRepo.GetContactRequestsByUserID(ctx, userID)

	// if the query returns nil, we want to show that there are no contact requests to this user
	if err == nil && contacts == nil {
		return []*pkg.ContactDetails{}, nil
	}

	return contacts, err
}
