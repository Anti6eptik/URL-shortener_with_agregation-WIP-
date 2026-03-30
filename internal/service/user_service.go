package service

import "goshop/internal/repository"

type UserService struct {
	userRepo *repository.UserRepository
	passSrv  *PasswordService
}

func NewUserService(
	userRepo *repository.UserRepository,
	passSrv *PasswordService,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		passSrv:  passSrv,
	}
}

func (s *UserService) Login(email, password string) error {

}
