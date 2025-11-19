package services

import (
	"errors"
	"s4s-backend/internal/modules/auth/dto"
	"s4s-backend/internal/modules/auth/models"
	authRepository "s4s-backend/internal/modules/auth/repository"
	notificationRepository "s4s-backend/internal/modules/notification/repository"
	subscriptionRepository "s4s-backend/internal/modules/subscription/repository"
	"s4s-backend/internal/pkg/utils"
)

type AuthService struct {
	userRepo         *authRepository.UserRepository
	subscriptionRepo *subscriptionRepository.SubscriptionRepository
	notificationRepo *notificationRepository.NotificationRepository
	jwtSecret        string
}

func NewAuthService(
	userRepo *authRepository.UserRepository,
	subscriptionRepo *subscriptionRepository.SubscriptionRepository,
	notificationRepo *notificationRepository.NotificationRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		subscriptionRepo: subscriptionRepo,
		notificationRepo: notificationRepo,
		jwtSecret:        jwtSecret,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Phone:        req.Phone,
		City:         req.City,
		Role:         req.Role,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Create default subscription (freemium)
	//subscription := &subscriptionModels.Subscription{
	//	UserID:          user.ID,
	//	Plan:            "freemium",
	//	Status:          "active",
	//	WorkflowsLimit:  5,
	//	ExecutionsLimit: 100,
	//	StartedAt:       time.Now(),
	//}

	//if err := s.subscriptionRepo.Create(subscription); err != nil {
	//	return nil, err
	//}

	// Create default notification settings
	//notifSettings := &notificationModels.NotificationSettings{
	//	UserID:   user.ID,
	//	Email:    true,
	//	Slack:    false,
	//	Channels: []string{"errors"},
	//}

	//if err := s.notificationRepo.Create(notifSettings); err != nil {
	//	return nil, err
	//}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Check if active
	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) ForgotPassword(req *dto.ForgotPasswordRequest) error {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// TODO: Generate reset token and send email
	_ = user

	return nil
}
