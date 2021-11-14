package database

import (
	"context"
	"go.uber.org/zap"
	"webserver/config"
	"webserver/database/repository"
)

type Database interface {
	Open(config config.DB) error
	Close()

	BeginTransaction(ctx context.Context) (Database, error)
	Commit() error
	Rollback() error

	UserRepository() repository.UserRepository
}

func NewDatabase(logger *zap.Logger, config config.DB) (Database, error) {
	db := &sqlxDatabase{
		logger: logger,
	}

	if err := db.Open(config); err != nil {
		return nil, err
	}

	return db, nil
}
