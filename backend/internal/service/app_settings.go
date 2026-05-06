package service

import (
	"encoding/json"
	"ithozyeva/internal/repository"
	"sync"
	"time"
)

const settingsCacheTTL = 60 * time.Second

type appSettingsCacheEntry struct {
	raw     []byte
	fetched time.Time
}

// AppSettingsService — тонкий wrapper над app_settings с in-memory кэшем
// (TTL 60 сек). Кэш на уровне процесса; переживает один тикер фоновой
// задачи и достаточно агрессивно инвалидируется, чтобы изменение в
// админке стало видимым в течение минуты без перезапуска backend.
type AppSettingsService struct {
	repo  *repository.AppSettingsRepository
	cache map[string]appSettingsCacheEntry
	mu    sync.RWMutex
}

func NewAppSettingsService() *AppSettingsService {
	return &AppSettingsService{
		repo:  repository.NewAppSettingsRepository(),
		cache: make(map[string]appSettingsCacheEntry),
	}
}

func (s *AppSettingsService) getRaw(key string) []byte {
	s.mu.RLock()
	if e, ok := s.cache[key]; ok && time.Since(e.fetched) < settingsCacheTTL {
		raw := e.raw
		s.mu.RUnlock()
		return raw
	}
	s.mu.RUnlock()

	raw, err := s.repo.GetRaw(key)
	if err != nil {
		return nil
	}
	s.mu.Lock()
	s.cache[key] = appSettingsCacheEntry{raw: raw, fetched: time.Now()}
	s.mu.Unlock()
	return raw
}

// GetInt возвращает int-значение по ключу. Если ключа нет, значение null
// или невалидный JSON — возвращает defaultVal.
func (s *AppSettingsService) GetInt(key string, defaultVal int) int {
	raw := s.getRaw(key)
	if len(raw) == 0 || string(raw) == "null" {
		return defaultVal
	}
	var v int
	if err := json.Unmarshal(raw, &v); err != nil {
		return defaultVal
	}
	return v
}

// GetFloat возвращает float64 по ключу с тем же fallback'ом.
func (s *AppSettingsService) GetFloat(key string, defaultVal float64) float64 {
	raw := s.getRaw(key)
	if len(raw) == 0 || string(raw) == "null" {
		return defaultVal
	}
	var v float64
	if err := json.Unmarshal(raw, &v); err != nil {
		return defaultVal
	}
	return v
}
