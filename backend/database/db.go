package database

import (
	"database/sql"
	"fmt"
	"ithozyeva/config"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Migration struct {
	Filename  string `gorm:"primaryKey"`
	AppliedAt int64
}

func SetupDatabase() error {
	baseDSN := fmt.Sprintf(
		"host=%s user=%s password=%s port=%s dbname=%s sslmode=disable",
		config.CFG.Database.Host,
		config.CFG.Database.User,
		config.CFG.Database.Password,
		config.CFG.Database.Port,
		config.CFG.Database.Name,
	)

	baseDB, err := sql.Open("postgres", baseDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	defer baseDB.Close()
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", config.CFG.Database.Name)
	err = baseDB.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		fmt.Printf("Database '%s' does not exist, creating...\n", config.CFG.Database.Name)
		_, err = baseDB.Exec(fmt.Sprintf("CREATE DATABASE %s", config.CFG.Database.Name))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		fmt.Printf("Database '%s' created successfully\n", config.CFG.Database.Name)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.CFG.Database.Host,
		config.CFG.Database.User,
		config.CFG.Database.Password,
		config.CFG.Database.Name,
		config.CFG.Database.Port)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB.AutoMigrate(&Migration{})

	// Применение миграций
	migrationsDir := filepath.Join("database", "migrations")
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			var count int64
			if err := DB.Model(&Migration{}).Where("filename = ?", file.Name()).Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check migration status: %w", err)
			}

			if count > 0 {
				continue
			}

			filePath := filepath.Join(migrationsDir, file.Name())
			migrationSQL, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
			}

			tx := DB.Begin()
			if err := tx.Exec(string(migrationSQL)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
			}

			if err := tx.Create(&Migration{Filename: file.Name(), AppliedAt: time.Now().Unix()}).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to record migration %s: %w", file.Name(), err)
			}

			tx.Commit()
			fmt.Printf("Applied migration: %s\n", file.Name())
		}
	}

	return nil
}
