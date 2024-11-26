package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"cyansnbrst.com/order-info/config"
	_ "github.com/lib/pq"
)

func OpenDB(cfg *config.Config) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.PostgreSQL.User,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.DBName,
		cfg.PostgreSQL.SSLMode,
	)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.PostgreSQL.MaxOpenConns)
	db.SetMaxIdleConns(cfg.PostgreSQL.MaxIdleConns)
	duration, err := time.ParseDuration(cfg.PostgreSQL.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}