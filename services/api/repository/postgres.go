package repository

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *postgres
	pgOnce     sync.Once
	initErr    error
)

func NewPostgresDB(ctx context.Context, dsn string) (*postgres, error) {
	pgOnce.Do(func() {
		dbpool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			logrus.Errorf("Unable to create connection pool: %v\n", err.Error())
			return
		}

		if err := dbpool.Ping(ctx); err != nil {
			initErr = err
			logrus.Errorf("Unable to ping connection pool: %v", err.Error())
			dbpool.Close()
			return
		}

		pgInstance = &postgres{dbpool}
		logrus.Info("Successfully initialized Postgres connection pool")
	})

	if initErr != nil {
		return nil, initErr
	}

	return pgInstance, nil
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}
