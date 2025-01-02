package helloworld

import (
	"errors"
	"fmt"
	"net/mail"
	"sync"

	"github.com/google/uuid"
)

const (
	maxUsernameLength = 50
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateEmail    = errors.New("user with email already exists")
	ErrDuplicateUsername = errors.New("user with username already exists")
	ErrEmptyUsername     = errors.New("username cannot be empty")
	ErrEmptyEmail        = errors.New("email cannot be empty")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrUsernameTooLong   = errors.New("username exceeds maximum length")
)

type User struct {
	Username string
	Email    string
	UUID     uuid.UUID
}

func (u *User) Validate() error {
	if u.Username == "" {
		return ErrEmptyUsername
	}
	if len(u.Username) > maxUsernameLength {
		return ErrUsernameTooLong
	}
	if u.Email == "" {
		return ErrEmptyEmail
	}
	if !isValidEmail(u.Email) {
		return ErrInvalidEmail
	}
	return nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

type UserRepository interface {
	AddUser(user User) (User, error)
	GetUserByEmail(email string) (User, error)
	GetUserByUsername(username string) (User, error)
	ListUsers() []*User
}

type inMemoryUserRepository struct {
	users []*User
	mu    sync.RWMutex
}

func NewInMemoryUserRepository() UserRepository {
	return &inMemoryUserRepository{
		users: make([]*User, 0),
	}
}

func (r *inMemoryUserRepository) AddUser(user User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existingUser := range r.users {
		if existingUser.Email == user.Email {
			return user, fmt.Errorf("%w: %s", ErrDuplicateEmail, user.Email)
		}
		if existingUser.Username == user.Username {
			return user, fmt.Errorf("%w: %s", ErrDuplicateUsername, user.Username)
		}
	}

	user.UUID = uuid.New()

	r.users = append(r.users, &user)
	return user, nil
}

func (r *inMemoryUserRepository) GetUserByEmail(email string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return *user, nil
		}
	}
	return User{}, fmt.Errorf("%w: %s", ErrUserNotFound, email)
}

func (r *inMemoryUserRepository) GetUserByUsername(username string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return *user, nil
		}
	}
	return User{}, fmt.Errorf("%w: %s", ErrUserNotFound, username)
}

func (r *inMemoryUserRepository) ListUsers() []*User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*User, len(r.users))
	copy(users, r.users)
	return users
}
