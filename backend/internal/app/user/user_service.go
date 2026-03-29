package user

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/Kushian01100111/Tickermaster/internal/domain/user"
	"github.com/Kushian01100111/Tickermaster/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailInvalid     = errors.New("invalid email type")
	ErrPasswordRequired = errors.New("password is required")
	ErrRole             = errors.New("invalid role type")
	ErrAuthMethod       = errors.New("invalid authmethod")
	ErrAuthMethodLen    = errors.New("invalid amount of authmethods")
)

type UserParams struct {
	Email      string
	Role       string
	Password   string
	AuthMethod string
}

type UpdateUserParams struct {
	Role             string
	Password         string
	AuthMethods      []string
	FailedLoginCount int32
	LastFailedLogin  *time.Time

	BookedEvents []string
}

type UserService interface {
	GetAllUsers(ctx context.Context) ([]user.User, error)
	CreateUser(params UserParams, ctx context.Context) (*user.User, error)
	GetUser(idhex string, ctx context.Context) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)

	UpdateUser(idhex string, params UpdateUserParams, ctx context.Context) (*user.User, error)
	DeleteUser(idhex string, ctx context.Context) error

	FailedLogin(ctx context.Context, user *user.User) error
	ResetFailedLogin(ctx context.Context, user *user.User) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s userService) GetAllUsers(ctx context.Context) ([]user.User, error) {
	return s.userRepo.GetAllUser(ctx)
}

func (s userService) CreateUser(params UserParams, ctx context.Context) (*user.User, error) {
	if err := validateParam(params); err != nil {
		return nil, err
	}

	failedLoginDate := time.Now().Add(-1 * time.Hour)

	var User *user.User

	if params.AuthMethod == "password" {
		hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		temp := string(hash)
		User = &user.User{
			Email:           params.Email,
			Role:            params.Role,
			PasswordHash:    &temp,
			AuthMethods:     []string{params.AuthMethod},
			LastFailedLogin: &failedLoginDate,
		}
	} else if params.AuthMethod == "email_otp" {
		User = &user.User{
			Email:           params.Email,
			Role:            params.Role,
			AuthMethods:     []string{params.AuthMethod},
			LastFailedLogin: &failedLoginDate,
		}
	}

	id, err := s.userRepo.Create(User, ctx)
	if err != nil {
		return nil, err
	}

	User.ID = id
	return User, err
}

func (s userService) GetUser(idhex string, ctx context.Context) (*user.User, error) {
	id, err := bson.ObjectIDFromHex(idhex)
	if err != nil {
		return nil, err
	}

	return s.userRepo.GetByID(id, ctx)
}

func (s userService) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	mail := strings.TrimSpace(strings.ToLower(email))
	if err := validateEmail(mail); err != nil {
		return nil, err
	}
	return s.userRepo.GetByEmail(mail, ctx)
}

func (s userService) UpdateUser(idhex string, params UpdateUserParams, ctx context.Context) (*user.User, error) {
	if err := validateUpdateParam(params); err != nil {
		return nil, err
	}

	id, err := bson.ObjectIDFromHex(idhex)
	if err != nil {
		return nil, err
	}

	var User *user.User

	if passwordMethodOk(params.AuthMethods) {
		hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		temp := string(hash)

		User = &user.User{
			Role:             params.Role,
			PasswordHash:     &temp,
			AuthMethods:      params.AuthMethods,
			FailedLoginCount: params.FailedLoginCount,
			LastFailedLogin:  params.LastFailedLogin,
		}
	} else {
		User = &user.User{
			Role:             params.Role,
			AuthMethods:      params.AuthMethods,
			FailedLoginCount: params.FailedLoginCount,
			LastFailedLogin:  params.LastFailedLogin,
		}
	}

	User, err = bookedEventsIDhexToObjectId(User, params)
	if err != nil {
		return nil, err
	}

	User.ID = id
	if err := s.userRepo.UpdateUser(User, ctx); err != nil {
		return nil, err
	}

	return User, nil
}

func (s userService) DeleteUser(idhex string, ctx context.Context) error {
	id, err := bson.ObjectIDFromHex(idhex)
	if err != nil {
		return err
	}

	if err := s.userRepo.DeleteUser(id, ctx); err != nil {
		return err
	}

	return nil
}

func (s *userService) FailedLogin(ctx context.Context, user *user.User) error
func (s *userService) ResetFailedLogin(ctx context.Context, user *user.User) error

///

func validateParam(params UserParams) error {
	email := strings.TrimSpace(strings.ToLower(params.Email))
	if err := validateEmail(email); err != nil {
		return err
	}

	if err := validateRole(params.Role); err != nil {
		return err
	}

	if err := validateAuthMethod(params.AuthMethod); err != nil {
		return ErrAuthMethod
	}

	if params.AuthMethod == "password" {
		if strings.TrimSpace(params.Password) == "" {
			return ErrPasswordRequired
		}
	}

	return nil
}

func validateUpdateParam(params UpdateUserParams) error {
	if err := validateRole(params.Role); err != nil {
		return err
	}

	if err := validateAuthMethods(params.AuthMethods); err != nil {
		return err
	}

	if passwordMethodOk(params.AuthMethods) {
		if strings.TrimSpace(params.Password) == "" {
			return ErrPasswordRequired
		}
	}

	return nil
}

func bookedEventsIDhexToObjectId(User *user.User, params UpdateUserParams) (*user.User, error) {
	bookedEvents := make([]bson.ObjectID, len(params.BookedEvents))

	for i, event := range params.BookedEvents {
		eventId, err := bson.ObjectIDFromHex(event)
		if err != nil {
			return nil, err
		}
		bookedEvents[i] = eventId
	}

	User.BookedEvents = bookedEvents
	return User, nil
}

func validateEmail(email string) error {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return ErrEmailInvalid
	}

	if addr.Address != email {
		return ErrEmailInvalid
	}

	return nil
}

func passwordMethodOk(methods []string) bool {
	for _, val := range methods {
		if val == "password" {
			return true
		}
	}
	return false
}

type Role string

const (
	RoleCostumer Role = "costumer"
	RoleEditor   Role = "editor"
	RoleAdmin    Role = "admin"
)

func validateRole(role string) error {
	switch Role(role) {
	case RoleCostumer, RoleEditor, RoleAdmin:
		return nil
	default:
		return ErrRole
	}
}

type AuthMethod string

const (
	AuthPass  AuthMethod = "password"
	AuthEmail AuthMethod = "email_otp"
)

func validateAuthMethod(method string) error {
	switch AuthMethod(method) {
	case AuthPass, AuthEmail:
		return nil
	default:
		return ErrAuthMethod
	}
}

func validateAuthMethods(methods []string) error {
	if len(methods) > 2 {
		return ErrAuthMethodLen
	}

	for _, val := range methods {
		if err := validateAuthMethod(val); err != nil {
			return err
		}
	}

	return nil
}
