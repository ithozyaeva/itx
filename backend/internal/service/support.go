package service

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// SupportService держит состояние «юзер ждёт, пока мы перешлём его
// следующее сообщение супер-админу» в Redis с авто-истечением. Это
// простой ticket-flow из welcome-меню бота: нажал «Написать админу»,
// одно сообщение уходит — и state сбрасывается.
type SupportService struct {
	redis *redis.Client
}

const (
	supportAwaitingTTL    = 10 * time.Minute
	supportRateLimitTTL   = 1 * time.Minute
	supportAwaitingPrefix = "support:awaiting:"
	supportRatePrefix     = "support:rate:"
)

func NewSupportService(redisClient *redis.Client) *SupportService {
	return &SupportService{redis: redisClient}
}

// BeginTicket помечает, что от userID ждём одно сообщение. Возвращает
// ошибку ErrSupportRateLimited, если ticket уже начинался менее
// supportRateLimitTTL назад — чтобы избежать спама/повторных нажатий.
func (s *SupportService) BeginTicket(ctx context.Context, userID int64) error {
	rateKey := fmt.Sprintf("%s%d", supportRatePrefix, userID)
	set, err := s.redis.SetNX(ctx, rateKey, "1", supportRateLimitTTL).Result()
	if err != nil {
		return err
	}
	if !set {
		return ErrSupportRateLimited
	}
	awaitKey := fmt.Sprintf("%s%d", supportAwaitingPrefix, userID)
	return s.redis.Set(ctx, awaitKey, "1", supportAwaitingTTL).Err()
}

// IsAwaiting возвращает true, если от userID ещё ждём ticket-сообщение.
func (s *SupportService) IsAwaiting(ctx context.Context, userID int64) bool {
	n, err := s.redis.Exists(ctx, fmt.Sprintf("%s%d", supportAwaitingPrefix, userID)).Result()
	if err != nil {
		return false
	}
	return n > 0
}

// EndTicket сбрасывает awaiting (но не rate-limit — он своим TTL отстоится).
func (s *SupportService) EndTicket(ctx context.Context, userID int64) error {
	return s.redis.Del(ctx, fmt.Sprintf("%s%d", supportAwaitingPrefix, userID)).Err()
}

// ErrSupportRateLimited — BeginTicket нельзя выполнить, юзер слишком
// часто открывает ticket.
var ErrSupportRateLimited = fmt.Errorf("support: rate limited")
