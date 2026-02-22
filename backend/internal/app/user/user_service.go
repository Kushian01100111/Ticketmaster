package user

import "github.com/Kushian01100111/Tickermaster/internal/repository"

type UserService interface {
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserRepository(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}
