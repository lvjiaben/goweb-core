package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type PostgresConfig struct {
	DSN             string        `yaml:"dsn"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
	LogLevel        string        `yaml:"log_level"`
}

func OpenPostgres(cfg PostgresConfig, logger *slog.Logger) (*gorm.DB, error) {
	if strings.TrimSpace(cfg.DSN) == "" {
		return nil, fmt.Errorf("postgres dsn is empty")
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: newGormLogger(logger, cfg.LogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}
	configurePool(sqlDB, cfg)
	return db, nil
}

func configurePool(sqlDB *sql.DB, cfg PostgresConfig) {
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}
}

func newGormLogger(logger *slog.Logger, level string) gormlogger.Interface {
	logLevel := gormlogger.Warn
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "silent":
		logLevel = gormlogger.Silent
	case "info":
		logLevel = gormlogger.Info
	case "error":
		logLevel = gormlogger.Error
	}

	return gormlogger.New(
		&slogWriter{logger: logger},
		gormlogger.Config{
			SlowThreshold:             500 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}

type slogWriter struct {
	logger *slog.Logger
}

func (w *slogWriter) Printf(format string, args ...any) {
	if w.logger == nil {
		return
	}
	w.logger.Info(fmt.Sprintf(format, args...))
}
