// Package testutil содержит хелперы для интеграционных тестов.
//
// SetupTestDB поднимает PostgreSQL в Docker через testcontainers-go,
// прогоняет на нём все миграции из database/migrations/ и возвращает
// готовый *gorm.DB. Контейнер автоматически останавливается через
// t.Cleanup. На пакет рекомендуется поднимать одну БД и truncate-ить
// между тестами через TruncateAll — старт контейнера дороже самих
// тестов.
//
// Тесты, использующие этот пакет, требуют запущенного Docker. Если
// Docker недоступен (TEST_SKIP_DB=1 или ошибка соединения), тесты
// должны вызывать t.Skip — за это отвечает SetupTestDB.
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
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dbImage = "postgres:15-alpine"

// SetupTestDB поднимает свежую PostgreSQL под тест и применяет все
// миграции. Возвращает *gorm.DB; останов и очистку регистрирует через
// t.Cleanup.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	if os.Getenv("TEST_SKIP_DB") == "1" {
		t.Skip("TEST_SKIP_DB=1, пропускаю интеграционный тест")
	}

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
		t.Skipf("не удалось поднять postgres-контейнер (Docker недоступен?): %v", err)
	}

	t.Cleanup(func() {
		// Используем отдельный контекст с таймаутом — t.Cleanup иногда
		// вызывается уже после отмены родительского.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := container.Terminate(ctx); err != nil {
			t.Logf("не удалось остановить контейнер: %v", err)
		}
	})

	dsn, err := container.ConnectionString(ctx, "sslmode=disable", "TimeZone=UTC")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm open: %v", err)
	}

	// Привязываем глобальный database.DB — много кода использует его
	// напрямую (PointsRepository.AwardPoints и т.п.). Сохраняем старое
	// значение и восстанавливаем в Cleanup, чтобы соседи не моргнули.
	prevDB := database.DB
	database.DB = db
	t.Cleanup(func() { database.DB = prevDB })

	if err := applyMigrations(db); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	return db
}

// TruncateAll очищает заданные таблицы — обычно вызывается перед
// каждым тестом одной горутины, чтобы между тестами не текли данные.
// CASCADE снимает FK, RESTART IDENTITY сбрасывает sequence-ы, чтобы
// id-ы были предсказуемыми.
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

// applyMigrations выполняет все .sql в database/migrations/ в порядке
// имени файла. Для тестов проще, чем поднимать database.SetupDatabase()
// — не нужно ни viper, ни env-переменных.
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

// migrationsDir возвращает абсолютный путь к database/migrations/.
// Тесты могут запускаться из любого каталога пакета — берём корень
// модуля относительно положения этого файла.
func migrationsDir() string {
	_, thisFile, _, _ := runtime.Caller(0)
	// .../backend/internal/testutil/db.go → .../backend/database/migrations
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "database", "migrations")
}
