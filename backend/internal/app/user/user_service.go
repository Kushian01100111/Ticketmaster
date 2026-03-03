package user

import (
	"context"

	"github.com/Kushian01100111/Tickermaster/internal/domain/user"
	"github.com/Kushian01100111/Tickermaster/internal/repository"
)

// Falta construir las estructuras -> 3/03
type UserParams struct{}

type UpdateParams struct{}

type LoginParams struct{}

type PasswordLessLogin struct{}

type PasswordLessSignUp struct{}

type Token struct{}

type TokenSignUp struct{}

type UserService interface {
	GetAllUser(ctx context.Context) ([]user.User, error)
	CreateUser(user UserParams, ctx context.Context) (*user.User, error) // ACA falta el JWT de la sesion creada con usuario, tambien estaria faltando en el resto de los metodos
	GetUser(idHex string, ctx context.Context) (*user.User, error)
	UpdateUser(idHex string, user UpdateParams, ctx context.Context) (*user.User, error)
	DeleteUser(idHex string, ctx context.Context) error
	Login(user LoginParams, ctx context.Context) (*user.User, error)               //-> JWT
	RequestToken(user PasswordLessLogin, ctx context.Context) error                // -> JWT
	LoginPasswordless(token Token, ctx context.Context) error                      // -> JWT
	SignUpRequestToken(user PasswordLessSignUp, ctx context.Context) error         // -> JWT
	SignUpPasswordless(token TokenSignUp, ctx context.Context) (*user.User, error) // -> JWT
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserRepository(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (a *userService) GetAllUser(ctx context.Context) ([]user.User, error)
func (a *userService) CreateUser(user UserParams, ctx context.Context) (*user.User, error) // ACA falta el JWT de la sesion creada con usuario, tambien estaria faltando en el resto de los metodos
func (a *userService) GetUser(idHex string, ctx context.Context) (*user.User, error)
func (a *userService) UpdateUser(idHex string, user UpdateParams, ctx context.Context) (*user.User, error)
func (a *userService) DeleteUser(idHex string, ctx context.Context) error
func (a *userService) Login(user LoginParams, ctx context.Context) (*user.User, error)               //-> JWT
func (a *userService) RequestToken(user PasswordLessLogin, ctx context.Context) error                // -> JWT
func (a *userService) LoginPasswordless(token Token, ctx context.Context) error                      // -> JWT
func (a *userService) SignUpRequestToken(user PasswordLessSignUp, ctx context.Context) error         // -> JWT
func (a *userService) SignUpPasswordless(token TokenSignUp, ctx context.Context) (*user.User, error) // -> JWT
