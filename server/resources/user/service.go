package user

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Mattcazz/Chat-TUI/pkg"
	"github.com/Mattcazz/Chat-TUI/server/db"
	"github.com/Mattcazz/Chat-TUI/server/resources/middleware"
	"github.com/Mattcazz/Chat-TUI/server/utils"
)

type Service struct {
	userRepo      UserRepository
	contactRepo   ContactRepository
	challengeRepo ChallengeRepository
	tx            *db.TxManager
}

func NewService(userRepo UserRepository, contactRepo ContactRepository, challengeRepo ChallengeRepository, tx *db.TxManager) *Service {
	return &Service{
		userRepo:      userRepo,
		contactRepo:   contactRepo,
		challengeRepo: challengeRepo,
		tx:            tx,
	}
}

func (s *Service) CreateUser(ctx context.Context, publicKey, username string) (*pkg.UserResponse, error) {
	log.Printf("CreateUser: Checking if user with username %s already exists", username)
	user, err := s.userRepo.GetUserByPublicKey(ctx, publicKey)

	if err == nil && user != nil {
		log.Printf("CreateUser: User with username %s already exists", username)
		return nil, fmt.Errorf("User already exists")
	}

	log.Printf("CreateUser: Creating new user with username %s", username)
	u := &User{
		PublicKey: publicKey,
		Username:  username,
	}

	user, err = s.userRepo.CreateUser(ctx, u)
	if err != nil {
		log.Printf("CreateUser: Failed to create user with username %s: %v", username, err)
		return nil, err
	}

	log.Printf("CreateUser: Successfully created user with ID %d, username %s", user.ID, user.Username)
	return &pkg.UserResponse{
		Username: user.Username,
		ID:       user.ID,
	}, nil
}

func (s *Service) DeleteUser(ctx context.Context, userID int64) error {
	log.Printf("DeleteUser: Attempting to delete user with ID %d", userID)
	err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		log.Printf("DeleteUser: Failed to delete user with ID %d: %v", userID, err)
		return err
	}
	log.Printf("DeleteUser: Successfully deleted user with ID %d", userID)
	return nil
}

func (s *Service) GenerateChallenge(ctx context.Context, publicKey string) (string, error) {
	log.Printf("GenerateChallenge: Generating challenge for user with public key")
	user, err := s.userRepo.GetUserByPublicKey(ctx, publicKey)
	if err != nil {
		log.Println("GenerateChallenge: User does not exist for the provided public key")
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

	log.Printf("GenerateChallenge: Creating challenge for user ID %d with nonce expiring at %v", user.ID, expires_at)
	err = s.challengeRepo.CreateChallenge(ctx, challenge)
	if err != nil {
		log.Printf("GenerateChallenge: Failed to create challenge for user ID %d: %v", user.ID, err)
		return "", err
	}

	log.Printf("GenerateChallenge: Successfully generated challenge for user ID %d", user.ID)
	return nonce, nil
}

func (s *Service) VerifyAndLogin(ctx context.Context, publicKey, signature string) (string, error) {
	log.Printf("VerifyAndLogin: Verifying login for user")
	user, err := s.userRepo.GetUserByPublicKey(ctx, publicKey)
	if err != nil {
		log.Printf("VerifyAndLogin: User not found for login verification: %v", err)
		err = NewUserDoesNotExistError()
		return "", err
	}

	log.Printf("VerifyAndLogin: Retrieving challenge for user ID %d", user.ID)
	challenge, err := s.challengeRepo.GetChallenge(ctx, user.ID)

	defer s.challengeRepo.DeleteChallenge(ctx, user.ID, challenge.Nonce)

	if err != nil {
		log.Printf("VerifyAndLogin: Error retrieving challenge for user %d: %v", user.ID, err)
		return "", fmt.Errorf("Challenge not created")
	}

	if challenge.ExpiresAt.Before(time.Now()) {
		log.Printf("VerifyAndLogin: Challenge expired for user ID %d", user.ID)
		return "", fmt.Errorf("Challenge expired")
	}

	log.Printf("VerifyAndLogin: Validating SSH signature for user ID %d", user.ID)
	if err := middleware.IsValidSshSignature(publicKey, challenge.Nonce, signature); err != nil {
		log.Printf("VerifyAndLogin: Invalid SSH signature for user ID %d: %v", user.ID, err)
		return "", err
	}

	log.Printf("VerifyAndLogin: Creating JWT token for user ID %d", user.ID)
	token, err := middleware.CreateJWT(user.ID)
	if err != nil {
		log.Printf("VerifyAndLogin: Failed to create JWT for user ID %d: %v", user.ID, err)
		return "", err
	}

	log.Printf("VerifyAndLogin: Successfully created JWT token for user ID %d", user.ID)
	return token, nil
}

func (s *Service) GetInbox(ctx context.Context, userID int64) (*pkg.InboxResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("User with ID %d does not exist", userID)
	}

	conversations, err := s.userRepo.GetUserConversations(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := &pkg.InboxResponse{
		Conversations: conversations,
		User: &pkg.UserResponse{
			Username: user.Username,
			ID:       user.ID,
		},
	}

	return response, nil
}

func (s *Service) GetContacts(ctx context.Context, userID int64) ([]*pkg.ContactDetails, error) {
	log.Printf("GetContacts: Fetching contacts for user ID %d", userID)
	contacts, err := s.contactRepo.GetContactsByUserID(ctx, userID)
	if err != nil {
		log.Printf("GetContacts: Failed to retrieve contacts for user ID %d: %v", userID, err)
		return nil, err
	}

	// if the query returns nil, we want to show that there are no contacts from this user
	if contacts == nil {
		log.Printf("GetContacts: No contacts found for user ID %d", userID)
		return []*pkg.ContactDetails{}, nil
	}

	log.Printf("GetContacts: Successfully retrieved %d contacts for user ID %d", len(contacts), userID)
	return contacts, nil
}

func (s *Service) ContactRequest(ctx context.Context, fromUserID int64, toPk, nickname string) error {
	log.Printf("ContactRequest: Processing contact request from user ID %d", fromUserID)

	fromUser, err := s.userRepo.GetUserByID(ctx, fromUserID)
	if err != nil {
		log.Printf("ContactRequest: Failed to retrieve from user with ID %d: %v", fromUserID, err)
		return fmt.Errorf("From User with ID %d does not exist", fromUserID)
	}

	if fromUser.PublicKey == toPk {
		log.Printf("ContactRequest: User %d attempted to add themselves as a contact", fromUserID)
		return fmt.Errorf("Ya can't be friends with yourself, buddy")
	}

	toUser, err := s.userRepo.GetUserByPublicKey(ctx, toPk)
	if err != nil {
		log.Printf("ContactRequest: Failed to retrieve target user: %v", err)
		return fmt.Errorf("User with public key %s does not exist", toPk)
	}

	log.Printf("ContactRequest: Creating contact request from user %d to user %d with nickname %s", fromUserID, toUser.ID, nickname)
	status := StatusPending

	tx, err := s.tx.StartTx(ctx)
	defer tx.Rollback()

	if err != nil {
		log.Printf("ContactRequest: Failed to start transaction: %v", err)
		return err
	}

	contact, err := s.contactRepo.GetContactByPair(ctx, toUser.ID, fromUserID)

	if err == nil && contact != nil {
		// contact already exists, we accept the contact request (both ways)
		log.Printf("ContactRequest: Bidirectional contact found, accepting request from user %d", fromUserID)
		contact.Status = StatusAccept
		status = StatusAccept
		contact.UpdatedAt = time.Now()

		err = s.contactRepo.WithTx(tx).UpdateContact(ctx, contact)
		if err != nil {
			log.Printf("ContactRequest: Failed to update existing contact: %v", err)
			return err
		}
	} else {
		log.Printf("ContactRequest: No bidirectional contact found: %v", err)
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
		CreatedAt:  time.Now(),
	}

	err = s.contactRepo.WithTx(tx).CreateContact(ctx, c)
	if err == nil {
		log.Printf("ContactRequest: Successfully created contact request from user %d to user %d", fromUserID, toUser.ID)
		tx.Commit()
	} else {
		log.Printf("ContactRequest: Failed to create contact: %v", err)
	}

	return err
}

func (s *Service) GetContactRequests(ctx context.Context, userID int64) ([]*pkg.ContactDetails, error) {
	log.Printf("GetContactRequests: Fetching contact requests for user ID %d", userID)
	contacts, err := s.contactRepo.GetContactRequestsByUserID(ctx, userID)
	if err != nil {
		log.Printf("GetContactRequests: Failed to retrieve contact requests for user ID %d: %v", userID, err)
		return nil, err
	}

	// if the query returns nil, we want to show that there are no contact requests to this user
	if contacts == nil {
		log.Printf("GetContactRequests: No contact requests found for user ID %d", userID)
		return []*pkg.ContactDetails{}, nil
	}

	log.Printf("GetContactRequests: Successfully retrieved %d contact requests for user ID %d", len(contacts), userID)
	return contacts, nil
}

func (s *Service) BlockContact(ctx context.Context, userID, contactID int64) error {
	log.Printf("BlockContact: Blocking contact ID %d for user ID %d", contactID, userID)
	c := &Contact{
		ID:        contactID,
		Status:    StatusBlocked,
		UpdatedAt: time.Now(),
	}

	err := s.contactRepo.UpdateContact(ctx, c)
	if err != nil {
		log.Printf("BlockContact: Failed to block contact ID %d for user ID %d: %v", contactID, userID, err)
		return err
	}

	log.Printf("BlockContact: Successfully blocked contact ID %d for user ID %d", contactID, userID)
	return nil
}

func (s *Service) UnblockContact(ctx context.Context, userID, contactID int64) error {
	log.Printf("UnblockContact: Unblocking contact ID %d for user ID %d", contactID, userID)
	c := &Contact{
		ID:        contactID,
		Status:    StatusAccept,
		UpdatedAt: time.Now(),
	}

	err := s.contactRepo.UpdateContact(ctx, c)
	if err != nil {
		log.Printf("UnblockContact: Failed to unblock contact ID %d for user ID %d: %v", contactID, userID, err)
		return err
	}

	log.Printf("UnblockContact: Successfully unblocked contact ID %d for user ID %d", contactID, userID)
	return nil
}
