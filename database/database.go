package database

import (
	"context"
	"go.uber.org/zap"
	"webserver/database/repository"
)

type Database interface {
	Open(driverName string, dataSourceName string) error
	Close() error

	BeginTransaction(ctx context.Context) (Database, error)
	Commit() error
	Rollback() error

	UserRepository() repository.UserRepository
}

func NewDatabase(logger *zap.Logger) Database {
	return &sqlxDatabase{
		logger: logger,
	}
}
