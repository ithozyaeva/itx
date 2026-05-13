package service

import (
	"fmt"
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"ithozyeva/internal/utils"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm/clause"
)

// ChatQuestGenerator выбирает «тихие» tracked-чаты (без якорных) и создаёт
// для них однодневный авто-квест на 5 содержательных сообщений. Идемпотентен:
// уникальный индекс (chat_id, DATE(starts_at)) WHERE auto_generated не даст
// задвоить квест при повторном запуске cron в тот же день.
type ChatQuestGenerator struct{}

func NewChatQuestGenerator() *ChatQuestGenerator {
	return &ChatQuestGenerator{}
}

// quietChatCandidate — кандидат на буст: чат, активность которого вчера
// заметно ниже своей же медианы за 14 дней.
type quietChatCandidate struct {
	ChatID       int64   `gorm:"column:chat_id"`
	Title        string  `gorm:"column:title"`
	MedianDaily  float64 `gorm:"column:median_daily"`
	YesterdayCnt int     `gorm:"column:yesterday_cnt"`
}

// quietChatThresholds — параметры отбора тихих чатов. Вынесены в одно
// место, чтобы крутить без поиска по коду. medianMin=3 отсекает «мёртвые»
// чаты — буст бессмысленен, если активность никогда не была. dropRatio=0.5
// = вчера было меньше половины медианы.
const (
	medianLookbackDays = 14
	medianMinFloor     = 3.0
	dropRatio          = 0.5
	severeDropRatio    = 0.2 // если вчера < 20% медианы → повышенная награда

	questTargetCount   = 5
	questBaseReward    = 10
	questSevereBonus   = 5
	questDurationHours = 24

	cooldownSuccess = 3 * 24 * time.Hour
	cooldownFail    = 7 * 24 * time.Hour
)

// questTemplates — короткие тексты для авто-квестов. {{title}} подставится
// в description. Без LLM: детерминировано, бесплатно, ревьюится глазами.
var questTemplates = []struct {
	title       string
	description string
}{
	{"Оживи чат", "В чате «%s» давно тихо. Расскажи, чем занят, или задай вопрос."},
	{"Подкинь тему", "Что нового изучаешь? Поделись с участниками чата «%s»."},
	{"Разогрей чат", "Поделись свежей мыслью, кейсом или ссылкой в «%s»."},
	{"Не дай чату уснуть", "Расскажи в «%s», над чем работаешь."},
	{"Запусти обсуждение", "Задай вопрос или поделись опытом в чате «%s»."},
}

// GenerateDailyChatQuests — главная точка входа из cron. Возвращает число
// созданных квестов и первую ошибку, если что-то пошло не так (но продолжает
// обработку остальных кандидатов, чтобы одна сбойная запись не валила весь
// прогон).
func (g *ChatQuestGenerator) GenerateDailyChatQuests() (int, error) {
	candidates, err := g.selectQuietChats()
	if err != nil {
		return 0, fmt.Errorf("select quiet chats: %w", err)
	}
	if len(candidates) == 0 {
		log.Printf("[chat-quest-gen] no quiet chats found")
		return 0, nil
	}

	cooldowns, err := g.loadAutoQuestCooldowns()
	if err != nil {
		return 0, fmt.Errorf("load cooldowns: %w", err)
	}

	now := time.Now()
	created := 0

	for _, c := range candidates {
		if blocked, reason := g.isOnCooldown(cooldowns[c.ChatID], now); blocked {
			log.Printf("[chat-quest-gen] skip chat=%d (%s): %s", c.ChatID, c.Title, reason)
			continue
		}

		quest := g.buildQuest(c, now)
		if err := g.insertIfNotExists(quest); err != nil {
			log.Printf("[chat-quest-gen] insert error chat=%d: %v", c.ChatID, err)
			continue
		}
		if quest.Id == 0 {
			// ON CONFLICT DO NOTHING — запись уже есть (рестарт cron в тот же день).
			continue
		}
		created++
		log.Printf("[chat-quest-gen] created quest=%d chat=%d (%s) target=%d reward=%d median=%.1f yesterday=%d",
			quest.Id, c.ChatID, c.Title, quest.TargetCount, quest.PointsReward, c.MedianDaily, c.YesterdayCnt)
	}

	log.Printf("[chat-quest-gen] done: %d quests created from %d candidates", created, len(candidates))
	return created, nil
}

// selectQuietChats считает медиану сообщений/день за 14 предыдущих дней
// (включая дни-нули через generate_series, иначе медиана искажена в пользу
// активных дней) и вытаскивает чаты, где вчерашняя активность ниже половины
// медианы. Якорные чаты исключаются: пользователь сказал считать их только
// точкой входа в другие чаты, а не объектом буста.
func (g *ChatQuestGenerator) selectQuietChats() ([]quietChatCandidate, error) {
	mskNow := time.Now().In(utils.MSKLocation())
	mskToday := time.Date(mskNow.Year(), mskNow.Month(), mskNow.Day(), 0, 0, 0, 0, utils.MSKLocation())
	yesterday := mskToday.AddDate(0, 0, -1).Format("2006-01-02")
	windowStart := mskToday.AddDate(0, 0, -medianLookbackDays).Format("2006-01-02")

	var rows []quietChatCandidate
	err := database.DB.Raw(`
		WITH chats AS (
			SELECT tc.chat_id, tc.title
			FROM tracked_chats tc
			WHERE tc.is_active = TRUE
			  AND tc.chat_id NOT IN (
			    SELECT id FROM subscription_chats WHERE anchor_for_tier_id IS NOT NULL
			  )
		),
		days AS (
			SELECT generate_series(?::date, (?::date - INTERVAL '1 day')::date, INTERVAL '1 day')::date AS day
		),
		chat_day_counts AS (
			SELECT c.chat_id, c.title, d.day,
			       COALESCE(COUNT(cm.id), 0)::int AS cnt
			FROM chats c
			CROSS JOIN days d
			LEFT JOIN chat_messages cm
			  ON cm.chat_id = c.chat_id AND DATE(cm.sent_at) = d.day
			GROUP BY c.chat_id, c.title, d.day
		),
		chat_stats AS (
			SELECT chat_id, title,
			       percentile_cont(0.5) WITHIN GROUP (ORDER BY cnt) AS median_daily,
			       COALESCE(SUM(CASE WHEN day = ?::date THEN cnt ELSE 0 END), 0)::int AS yesterday_cnt
			FROM chat_day_counts
			GROUP BY chat_id, title
		)
		SELECT chat_id, title, median_daily, yesterday_cnt
		FROM chat_stats
		WHERE median_daily >= ?
		  AND yesterday_cnt::float < median_daily * ?
		ORDER BY median_daily DESC
	`, windowStart, yesterday, yesterday, medianMinFloor, dropRatio).Scan(&rows).Error
	return rows, err
}

// autoQuestCooldown — последний авто-квест по чату и был ли он кем-то выполнен.
type autoQuestCooldown struct {
	ChatID     int64     `gorm:"column:chat_id"`
	LastEnd    time.Time `gorm:"column:last_end"`
	HadSuccess bool      `gorm:"column:had_success"`
}

// loadAutoQuestCooldowns собирает все авто-квесты прошлых дней с признаком
// «был ли хоть один участник, который его выполнил». Возвращает мапу для
// O(1) lookup в основном цикле.
func (g *ChatQuestGenerator) loadAutoQuestCooldowns() (map[int64]autoQuestCooldown, error) {
	var rows []autoQuestCooldown
	err := database.DB.Raw(`
		SELECT q.chat_id,
		       MAX(q.ends_at) AS last_end,
		       BOOL_OR(EXISTS (
		         SELECT 1 FROM chat_quest_progress p
		         WHERE p.quest_id = q.id AND p.completed = TRUE
		       )) AS had_success
		FROM chat_quests q
		WHERE q.auto_generated = TRUE AND q.chat_id IS NOT NULL
		GROUP BY q.chat_id
	`).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[int64]autoQuestCooldown, len(rows))
	for _, r := range rows {
		out[r.ChatID] = r
	}
	return out, nil
}

// isOnCooldown решает, можно ли создать новый авто-квест для чата прямо
// сейчас. Логика: после успешного квеста ждём 3 дня (не долбим уже-оживший
// чат), после полного провала — 7 дней (не реанимируем труп ежедневно).
func (g *ChatQuestGenerator) isOnCooldown(cd autoQuestCooldown, now time.Time) (bool, string) {
	if cd.ChatID == 0 {
		return false, "" // никогда не было авто-квеста
	}
	elapsed := now.Sub(cd.LastEnd)
	if cd.HadSuccess && elapsed < cooldownSuccess {
		return true, fmt.Sprintf("cooldown after success (%.1fh elapsed of %.0fh)", elapsed.Hours(), cooldownSuccess.Hours())
	}
	if !cd.HadSuccess && elapsed < cooldownFail {
		return true, fmt.Sprintf("cooldown after no-success (%.1fh elapsed of %.0fh)", elapsed.Hours(), cooldownFail.Hours())
	}
	return false, ""
}

// buildQuest формирует ChatQuest по кандидату. Шаблон выбирается случайно,
// награда — базовая или с бонусом, если просадка серьёзная (см. severeDropRatio).
func (g *ChatQuestGenerator) buildQuest(c quietChatCandidate, now time.Time) *models.ChatQuest {
	tpl := questTemplates[rand.Intn(len(questTemplates))]

	reward := questBaseReward
	if c.MedianDaily > 0 && float64(c.YesterdayCnt)/c.MedianDaily < severeDropRatio {
		reward += questSevereBonus
	}

	chatID := c.ChatID
	return &models.ChatQuest{
		Title:         tpl.title,
		Description:   fmt.Sprintf(tpl.description, c.Title),
		QuestType:     models.QuestTypeMessageCount,
		ChatID:        &chatID,
		TargetCount:   questTargetCount,
		PointsReward:  reward,
		StartsAt:      now,
		EndsAt:        now.Add(questDurationHours * time.Hour),
		IsActive:      true,
		AutoGenerated: true,
	}
}

// insertIfNotExists вставляет квест с ON CONFLICT DO NOTHING по партиальному
// индексу (chat_id, DATE(starts_at)) WHERE auto_generated. Если запись уже
// была (повторный запуск cron в тот же день) — id останется 0, caller это
// поймёт и не будет считать как создание.
func (g *ChatQuestGenerator) insertIfNotExists(q *models.ChatQuest) error {
	return database.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(q).Error
}
