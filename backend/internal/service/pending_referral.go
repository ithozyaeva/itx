package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// pendingReferralTTL — окно, в котором /start ref_<id> в боте даёт
// атрибуцию при последующем создании members-записи. 30 дней покрывают
// типичный кейс «увидел ссылку → дотянулся до бота через неделю»,
// но не вечно занимают Redis ключами от ушедших юзеров.
const pendingReferralTTL = 30 * 24 * time.Hour

func pendingReferralKey(telegramUserID int64) string {
	return fmt.Sprintf("referee:pending:%d", telegramUserID)
}

// PendingReferralService хранит «пришёл по реф-ссылке X, ещё не создан
// members-record» в Redis. Постоянное хранилище не нужно: Redis у нас
// уже в инфре, для атрибуции 30-дневный TTL достаточен. Когда auth-handler
// создаёт members-запись, он вычитывает значение и переносит в БД.
type PendingReferralService struct {
	redis *redis.Client
}

func NewPendingReferralService(redisClient *redis.Client) *PendingReferralService {
	return &PendingReferralService{redis: redisClient}
}

// Set записывает ассоциацию telegramUserID → linkID. Перезаписывает
// существующее значение: если бот вызывают повторно с другим payload'ом,
// побеждает последний — это намеренно, авторитет за свежим действием.
func (s *PendingReferralService) Set(ctx context.Context, telegramUserID int64, linkID int64) error {
	if s.redis == nil || telegramUserID == 0 || linkID == 0 {
		return nil
	}
	return s.redis.Set(ctx, pendingReferralKey(telegramUserID), strconv.FormatInt(linkID, 10), pendingReferralTTL).Err()
}

// GetAndDelete атомарно возвращает linkID и удаляет ключ. Используется
// auth-handler'ом ровно один раз — после выполнения linkID должен попасть
// в members.referred_by_link_id и больше не нужен. Возвращает (0, nil)
// если ключа нет.
func (s *PendingReferralService) GetAndDelete(ctx context.Context, telegramUserID int64) (int64, error) {
	if s.redis == nil || telegramUserID == 0 {
		return 0, nil
	}
	key := pendingReferralKey(telegramUserID)
	val, err := s.redis.GetDel(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}
