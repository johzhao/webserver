package database

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"webserver/config"
	"webserver/database/repository"
)

type sqlxDatabase struct {
	logger  *zap.Logger
	db      *sqlx.DB
	tx      *sqlx.Tx
	queryer sqlx.QueryerContext
	execer  sqlx.ExecerContext
}

func (s *sqlxDatabase) Open(config config.DB) error {
	db, err := sqlx.Open(config.DriverName, config.DataSourceName)
	if err != nil {
		return err
	}

	s.db = db
	s.queryer = db
	s.execer = db

	return nil
}

func (s *sqlxDatabase) Close() {
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			s.logger.Error("close database failed",
				zap.Error(err))
		}
		s.db = nil
	}
}

func (s *sqlxDatabase) BeginTransaction(ctx context.Context) (Database, error) {
	if s.db == nil {
		return s, nil
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &sqlxDatabase{
		db:      nil,
		tx:      tx,
		queryer: tx,
		execer:  tx,
	}, nil
}

func (s *sqlxDatabase) Commit() error {
	if s.tx == nil {
		return fmt.Errorf("not in transaction")
	}
	return s.tx.Commit()
}

func (s *sqlxDatabase) Rollback() error {
	if s.tx == nil {
		return fmt.Errorf("not in transaction")
	}
	return s.tx.Rollback()
}

func (s *sqlxDatabase) UserRepository() repository.UserRepository {
	return repository.NewUserRepository(s.logger, s.queryer, s.execer)
}
