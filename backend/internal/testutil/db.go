// Package testutil содержит хелперы для интеграционных тестов.
//
// Используется shared-модель: один Postgres-контейнер на тестовый пакет.
// EnsureTestDB поднимает контейнер при первом вызове в пакете, прогоняет
// миграции, подменяет database.DB и возвращает *gorm.DB. Контейнер живёт
// до конца жизни тестового процесса — его очистка ложится на ОС вместе
// с runtime'ом.
//
// Между тестами и подтестами используется TruncateAll(t, db, tables...)
// для изоляции данных. Контейнер не пересоздаётся.
//
// Если Docker недоступен или TEST_SKIP_DB=1 — EnsureTestDB вызывает
// t.Skip и тест проходит как пропущенный.
package testutil

import (
	"context"
	"fmt"
	"ithozyeva/database"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dbImage = "postgres:15-alpine"

var (
	once     sync.Once
	sharedDB *gorm.DB
	sharedErr error
	skipReason string
)

// EnsureTestDB поднимает один Postgres-контейнер на тестовый процесс
// (через sync.Once) и возвращает готовый *gorm.DB. Подменяет
// database.DB. Если Docker недоступен — t.Skip.
func EnsureTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	if os.Getenv("TEST_SKIP_DB") == "1" {
		t.Skip("TEST_SKIP_DB=1, пропускаю интеграционный тест")
	}

	once.Do(initSharedDB)

	if skipReason != "" {
		t.Skip(skipReason)
	}
	if sharedErr != nil {
		t.Fatalf("init shared db: %v", sharedErr)
	}

	// Привязываем глобальный database.DB. Восстанавливать в Cleanup
	// не нужно: shared instance живёт до конца процесса, никто другой
	// не пишет в database.DB параллельно (тесты в одном пакете идут
	// последовательно по умолчанию).
	database.DB = sharedDB
	return sharedDB
}

func initSharedDB() {
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx, dbImage,
		tcpostgres.WithDatabase("itx_test"),
		tcpostgres.WithUsername("itx"),
		tcpostgres.WithPassword("itx"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		skipReason = fmt.Sprintf("не удалось поднять postgres-контейнер (Docker недоступен?): %v", err)
		return
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable", "TimeZone=UTC")
	if err != nil {
		sharedErr = fmt.Errorf("connection string: %w", err)
		return
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		sharedErr = fmt.Errorf("gorm open: %w", err)
		return
	}

	if err := applyMigrations(db); err != nil {
		sharedErr = fmt.Errorf("apply migrations: %w", err)
		return
	}

	sharedDB = db
}

// TruncateAll очищает заданные таблицы — рекомендуется вызывать в
// начале каждого теста, чтобы между тестами не текли данные. CASCADE
// снимает FK, RESTART IDENTITY сбрасывает sequence-ы.
func TruncateAll(t *testing.T, db *gorm.DB, tables ...string) {
	t.Helper()
	if len(tables) == 0 {
		return
	}
	q := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(tables, ", "))
	if err := db.Exec(q).Error; err != nil {
		t.Fatalf("truncate %v: %v", tables, err)
	}
}

// SetupTestDB — устаревший alias на EnsureTestDB для обратной
// совместимости. Новый код должен вызывать EnsureTestDB напрямую.
func SetupTestDB(t *testing.T) *gorm.DB {
	return EnsureTestDB(t)
}

func applyMigrations(db *gorm.DB) error {
	dir := migrationsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read migrations dir %s: %w", dir, err)
	}

	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		raw, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		if err := db.Exec(string(raw)).Error; err != nil {
			return fmt.Errorf("apply %s: %w", name, err)
		}
	}
	return nil
}

func migrationsDir() string {
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "database", "migrations")
}
