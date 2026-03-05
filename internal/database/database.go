package database

import (
	"context"
	"fmt"

	"sds-admin/internal/config"
	"sds-admin/internal/logger"
	"sds-admin/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var db *gorm.DB

// Init initializes the database connection with configuration.
// Init 使用配置初始化数据库连接
func Init(cfg *config.DatabaseConfig) error {
	var err error

	gormLogger := gormlogger.Default
	if logger.GetLogger().GetLevel() == logrus.DebugLevel {
		gormLogger = gormlogger.Default.LogMode(gormlogger.Info)
	}

	db, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	logger.Info("Database connection established successfully")
	return nil
}

// autoMigrate automatically migrates the database schema.
// autoMigrate 自动迁移数据库 schema
func autoMigrate() error {
	logger.Info("Running database migrations...")
	if err := db.AutoMigrate(
		&models.Domain{},
		&models.RecordType{},
		&models.Record{},
		&models.RecordValue{},
	); err != nil {
		return err
	}
	logger.Info("Database migrations completed successfully")
	return nil
}

// GetDB returns the database instance.
// GetDB 返回数据库实例
func GetDB() *gorm.DB {
	return db
}

// Close closes the database connection.
// Close 关闭数据库连接
func Close() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// Ping checks if the database connection is alive.
// Ping 检查数据库连接是否存活
func Ping() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// Transaction executes a transaction with the given function.
// Transaction 使用给定函数执行事务
func Transaction(fn func(tx *gorm.DB) error) error {
	return db.Transaction(fn)
}

// WithContext creates a new DB instance with context.
// WithContext 创建带有context的新DB实例
func WithContext(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}
