package api

import (
	"context"
	"time"
)

// Service описывает бизнес-логику API модуля.
type Service interface {
	Health(ctx context.Context) (map[string]any, error)
	Info(ctx context.Context) (map[string]any, error)
}

type service struct {
	// TODO: добавьте зависимости при необходимости (например, логгер, конфиг, db-клиент и т.п.)
	// db *ent.Client
}

func NewService(
// db *ent.Client,
) Service {
	return &service{
		// db: db,
	}
}

func (s *service) Health(ctx context.Context) (map[string]any, error) {
	// TODO: добавить реальные проверки (доступность БД, внешних сервисов и пр.)
	return map[string]any{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (s *service) Info(ctx context.Context) (map[string]any, error) {
	// TODO: подставьте реальные данные о приложении/окружении
	return map[string]any{
		"name":    "app",
		"version": "0.1.0",
		"env":     "dev",
	}, nil
}
