package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// pendingReferralTTL — окно, в котором /start ref_<code> в боте даёт
// атрибуцию при последующем создании members-записи. 30 дней покрывают
// типичный кейс «увидел ссылку → дотянулся до бота через неделю».
const pendingReferralTTL = 30 * 24 * time.Hour

func pendingReferralKey(telegramUserID int64) string {
	return fmt.Sprintf("referee:pending:%d", telegramUserID)
}

// PendingReferralService хранит «пришёл по реф-ссылке инвайтера X, ещё
// не создан members-record» в Redis. Значение — referrer_member_id.
//
// Не путать с referal_links.id (вакансии): здесь сохраняется привязка
// именно к юзеру-инвайтеру в community-программе (PR #350+).
type PendingReferralService struct {
	redis *redis.Client
}

func NewPendingReferralService(redisClient *redis.Client) *PendingReferralService {
	return &PendingReferralService{redis: redisClient}
}

// Set записывает ассоциацию telegramUserID → referrerMemberID. Перезаписывает
// существующее значение: если бот вызвали повторно с другим payload'ом,
// побеждает последний — намеренно, авторитет за свежим действием.
func (s *PendingReferralService) Set(ctx context.Context, telegramUserID int64, referrerMemberID int64) error {
	if s.redis == nil || telegramUserID == 0 || referrerMemberID == 0 {
		return nil
	}
	return s.redis.Set(ctx, pendingReferralKey(telegramUserID), strconv.FormatInt(referrerMemberID, 10), pendingReferralTTL).Err()
}

// GetAndDelete атомарно возвращает referrerMemberID и удаляет ключ.
// Используется auth-handler'ом ровно один раз — после выполнения
// referrerID должен попасть в members.referred_by_member_id.
// Возвращает (0, nil) если ключа нет.
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
