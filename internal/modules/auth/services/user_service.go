package service

import (
	"errors"
	"your-project/internal/dto"
	"your-project/internal/models"
	"your-project/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUser(userID string) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *UserService) UpdateUser(userID string, req *dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.City != "" {
		user.City = req.City
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(userID string) error {
	return s.userRepo.Delete(userID)
}

func (s *UserService) InviteUser(req *dto.InviteUserRequest) error {
	// TODO: Send invitation email
	// For now, just validate
	if req.Email == "" {
		return errors.New("email is required")
	}

	return nil
}
