package db

import (
	"context"
	"fmt"
	"log"
	"s4s-backend/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var Redis *redis.Client
var RabbitConn *amqp.Connection

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.GetString("DB_HOST", "localhost"),
		config.GetString("DB_USER", "postgres"),
		config.GetString("DB_PASSWORD", "postgres"),
		config.GetString("DB_NAME", "s4s"),
		config.GetString("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error | logger.Info | logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	return db, nil
}

func ConnectRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.GetString("REDIS_ADDR", "localhost:6379"),
		Password: config.GetString("REDIS_PASSWORD", ""),
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	Redis = client
	return client, nil
}

func ConnectRabbitMQ() (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for i := 0; i < 20; i++ {
		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err == nil {
			log.Println("RabbitMQ connected successfully")
			return conn, nil
		}
		log.Printf("RabbitMQ not ready, retry %d/20... (%v)", i+1, err)
		time.Sleep(3 * time.Second)
	}
	return nil, fmt.Errorf("failed to connect to RabbitMQ after retries: %w", err)
}
