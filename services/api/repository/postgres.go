package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type PostresConfig struct {
	Host     string
	Port     string
	PG_User  string
	PG_Pwd   string
	PG_DB    string
	SSL_Mode string
}

type postgres struct {
	db      *pgxpool.Pool
	querier Querier
}

var (
	pgInstance *postgres
	pgOnce     sync.Once
	initErr    error
)

func NewPostgresDB(ctx context.Context, cfg PostresConfig) (Repository, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.PG_User,
		cfg.PG_Pwd,
		cfg.PG_DB,
		cfg.SSL_Mode,
	)

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

		queries := New(dbpool)

		pgInstance = &postgres{dbpool, queries}
		logrus.Info("Successfully initialized Postgres connection pool")

		migrateDSN := fmt.Sprintf("pgx5://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.PG_User, cfg.PG_Pwd, cfg.Host, cfg.Port, cfg.PG_DB, cfg.SSL_Mode)

		m, err := migrate.New("file://migrations", migrateDSN+"&x-migrations-table=order_migrations")
		if err != nil {
			initErr = err
			logrus.Errorf("Unable to initialize migrations: %v", err.Error())
			return
		}

		if err = m.Up(); err != nil && err != migrate.ErrNoChange {
			initErr = err
			logrus.Errorf("Unable to apply migrations")
			return
		}

		logrus.Info("migrations apllied successfully")
	})

	if initErr != nil {
		return nil, initErr
	}

	if pgInstance == nil {
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

func (p *postgres) CreateOrder(ctx context.Context, params CreateOrderParams) (Order, error) {
	return p.querier.CreateOrder(ctx, params)
}

func (p *postgres) GetOrderByID(ctx context.Context, id int32) (Order, error) {
	return p.querier.GetOrderByID(ctx, id)
}

func (p *postgres) ListOrders(ctx context.Context, params ListOrdersParams) ([]Order, error) {
	return p.querier.ListOrders(ctx, params)
}

func (p *postgres) UpdateOrderStatus(ctx context.Context, params UpdateOrderStatusParams) (Order, error) {
	return p.querier.UpdateOrderStatus(ctx, params)
}

func (p *postgres) DeleteOrderByID(ctx context.Context, id int32) (int64, error) {
	return p.querier.DeleteOrderByID(ctx, id)
}
