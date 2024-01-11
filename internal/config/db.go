package config

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDB(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Db.DSN))
	if err != nil {
		return nil, err
	}

	sqldb, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqldb.SetMaxOpenConns(cfg.Db.MaxOpenConn)
	sqldb.SetMaxIdleConns(cfg.Db.MaxIdleConn)

	idletime, err := time.ParseDuration(cfg.Db.MaxIdleTime)
	if err != nil {
		return nil, fmt.Errorf("failed parsing maxidletime: %v", err)
	}
	sqldb.SetConnMaxIdleTime(idletime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = sqldb.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
