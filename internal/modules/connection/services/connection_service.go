package services

import (
	"context"
	"s4s-backend/internal/modules/connection/models"
	connectionRepo "s4s-backend/internal/modules/connection/repository"
)

type ConnectionService interface {
	CreateConnection(ctx context.Context, conn *models.Connection) error
	GetConnection(ctx context.Context, id, userID string) (*models.Connection, error)
	ListConnections(ctx context.Context, userID string) ([]*models.Connection, error)
	UpdateConnection(ctx context.Context, conn *models.Connection) error
	DeleteConnection(ctx context.Context, id, userID string) error
	TestConnection(ctx context.Context, conn *models.Connection) (bool, error)
}

type connectionService struct {
	repo connectionRepo.ConnectionRepository
}

func NewConnectionService(repo connectionRepo.ConnectionRepository) ConnectionService {
	return &connectionService{repo: repo}
}

func (s *connectionService) CreateConnection(ctx context.Context, conn *models.Connection) error {
	return s.repo.Create(conn)
}

func (s *connectionService) GetConnection(ctx context.Context, id, userID string) (*models.Connection, error) {
	return s.repo.GetByID(id, userID)
}

func (s *connectionService) ListConnections(ctx context.Context, userID string) ([]*models.Connection, error) {
	return s.repo.FindByUserID(userID)
}

func (s *connectionService) UpdateConnection(ctx context.Context, conn *models.Connection) error {
	return s.repo.Update(conn)
}

func (s *connectionService) DeleteConnection(ctx context.Context, id, userID string) error {
	return s.repo.Delete(id, userID)
}

func (s *connectionService) TestConnection(ctx context.Context, conn *models.Connection) (bool, error) {
	// Implementation depends on the connection type
	// This is a basic implementation that always returns true
	return true, nil
}
