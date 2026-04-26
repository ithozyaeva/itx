package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"strconv"
	"strings"
	"time"
)

type ModerationService struct {
	repo *repository.ModerationRepository
}

func NewModerationService() *ModerationService {
	return &ModerationService{repo: repository.NewModerationRepository()}
}

// ParseHumanDuration принимает "30m", "1h", "1d", "7d", "12h30m" и пр.
// Поддерживает суффиксы s/m/h/d (день = 24h). Пустая строка → (0, nil).
func ParseHumanDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return 0, nil
	}

	// Расширяем "d" в "24h" и доверяем стандартному парсеру для остального.
	// Пример: "1d12h" → "36h0m".
	var sb strings.Builder
	i := 0
	for i < len(s) {
		j := i
		for j < len(s) && (s[j] == '.' || (s[j] >= '0' && s[j] <= '9')) {
			j++
		}
		if j == i || j >= len(s) {
			return 0, fmt.Errorf("неверный формат длительности: %q", s)
		}
		num := s[i:j]
		unit := s[j]
		i = j + 1
		switch unit {
		case 'd':
			n, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return 0, fmt.Errorf("неверное число: %q", num)
			}
			sb.WriteString(strconv.FormatFloat(n*24, 'f', -1, 64))
			sb.WriteByte('h')
		case 's', 'm', 'h':
			sb.WriteString(num)
			sb.WriteByte(unit)
		default:
			return 0, fmt.Errorf("неизвестная единица %q", string(unit))
		}
	}

	return time.ParseDuration(sb.String())
}

// FormatDurationHuman — короткая русская форма ("1ч", "30м", "2д").
func FormatDurationHuman(d time.Duration) string {
	if d <= 0 {
		return "навсегда"
	}
	if d%(24*time.Hour) == 0 {
		return fmt.Sprintf("%dд", int(d/(24*time.Hour)))
	}
	if d%time.Hour == 0 {
		return fmt.Sprintf("%dч", int(d/time.Hour))
	}
	if d%time.Minute == 0 {
		return fmt.Sprintf("%dм", int(d/time.Minute))
	}
	return d.String()
}

// LogAction сохраняет запись в журнал модерации.
func (s *ModerationService) LogAction(action *models.ModerationAction) error {
	return s.repo.LogAction(action)
}

// LogActionWithMeta — то же, но с произвольной meta-payload.
func (s *ModerationService) LogActionWithMeta(action *models.ModerationAction, meta map[string]interface{}) error {
	if meta != nil {
		raw, err := json.Marshal(meta)
		if err == nil {
			action.Meta = string(raw)
		}
	}
	return s.repo.LogAction(action)
}

// MessagesForCleanup — id телеграм-сообщений юзера в чате за период.
func (s *ModerationService) MessagesForCleanup(chatID, userID int64, since time.Time) ([]int, error) {
	return s.repo.MessagesForCleanup(chatID, userID, since)
}

// DeleteCleanedMessages удаляет записи из chat_messages после успешного удаления в Telegram.
func (s *ModerationService) DeleteCleanedMessages(chatID int64, messageIDs []int) (int64, error) {
	return s.repo.DeleteCleanedMessages(chatID, messageIDs)
}

// VotebanStartParams описывает входные данные старта голосования.
type VotebanStartParams struct {
	ChatID           int64
	ChatTitle        string
	TargetUserID     int64
	TargetUsername   string
	TargetFirstName  string
	InitiatorUserID  int64
	TriggerMessageID *int
	PollMessageID    int
	RequiredVotes    int
	MuteSeconds      int
	WindowSeconds    int
}

// StartVoteban создаёт запись об открытом голосовании. Возвращает
// ErrVotebanAlreadyOpen, если на (chat, target) уже есть открытое.
func (s *ModerationService) StartVoteban(p VotebanStartParams) (*models.Voteban, error) {
	existing, err := s.repo.FindOpenVoteban(p.ChatID, p.TargetUserID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, ErrVotebanAlreadyOpen
	}
	v := &models.Voteban{
		ChatID:           p.ChatID,
		ChatTitle:        p.ChatTitle,
		TargetUserID:     p.TargetUserID,
		TargetUsername:   p.TargetUsername,
		TargetFirstName:  p.TargetFirstName,
		InitiatorUserID:  p.InitiatorUserID,
		TriggerMessageID: p.TriggerMessageID,
		PollMessageID:    p.PollMessageID,
		RequiredVotes:    p.RequiredVotes,
		MuteSeconds:      p.MuteSeconds,
		ExpiresAt:        time.Now().Add(time.Duration(p.WindowSeconds) * time.Second),
		Status:           models.VotebanStatusOpen,
	}
	if err := s.repo.CreateVoteban(v); err != nil {
		return nil, err
	}
	return v, nil
}

// GetVoteban возвращает запись по id.
func (s *ModerationService) GetVoteban(id int64) (*models.Voteban, error) {
	return s.repo.GetVoteban(id)
}

// CastVoteResult — результат голосования.
type CastVoteResult struct {
	Tally    models.VotebanTally
	Voteban  *models.Voteban
	Threshold bool // достигнут ли порог "за"
	Changed   bool // изменился ли голос (false — повторный тот же)
}

// CastVote ставит/обновляет голос. Если порог достигнут, возвращает Threshold=true
// (финализация — на стороне вызывающего, чтобы он мог сделать Telegram-действие).
func (s *ModerationService) CastVote(votebanID, voterID int64, vote int16) (*CastVoteResult, error) {
	if vote != models.VotebanVoteFor && vote != models.VotebanVoteAgainst {
		return nil, fmt.Errorf("неверное значение голоса")
	}
	v, err := s.repo.GetVoteban(votebanID)
	if err != nil {
		return nil, err
	}
	if v.Status != models.VotebanStatusOpen {
		return &CastVoteResult{Voteban: v}, ErrVotebanClosed
	}
	if voterID == v.TargetUserID {
		return &CastVoteResult{Voteban: v}, ErrVoteSelfTarget
	}

	prev, err := s.repo.GetVote(votebanID, voterID)
	if err != nil {
		return nil, err
	}
	changed := prev == nil || *prev != vote

	if err := s.repo.UpsertVote(votebanID, voterID, vote); err != nil {
		return nil, err
	}
	tally, err := s.repo.CountVotes(votebanID)
	if err != nil {
		return nil, err
	}
	return &CastVoteResult{
		Tally:     tally,
		Voteban:   v,
		Threshold: tally.For >= v.RequiredVotes,
		Changed:   changed,
	}, nil
}

// FinalizeVoteban переводит запись в финальный статус. Идемпотентно:
// повторный вызов после успешной финализации ничего не делает (запись
// уже не "open"). Возвращает true, если перевод произошёл.
func (s *ModerationService) FinalizeVoteban(id int64, status string) (bool, error) {
	v, err := s.repo.GetVoteban(id)
	if err != nil {
		return false, err
	}
	if v.Status != models.VotebanStatusOpen {
		return false, nil
	}
	if err := s.repo.FinalizeVoteban(id, status); err != nil {
		return false, err
	}
	return true, nil
}

// ListExpiredOpenVotebans возвращает голосования, у которых истекло окно.
func (s *ModerationService) ListExpiredOpenVotebans(now time.Time) ([]models.Voteban, error) {
	return s.repo.ListExpiredOpenVotebans(now)
}

// CountVotes — текущая раскладка голосов.
func (s *ModerationService) CountVotes(votebanID int64) (models.VotebanTally, error) {
	return s.repo.CountVotes(votebanID)
}

var (
	ErrVotebanAlreadyOpen = errors.New("voteban: на этого пользователя уже идёт голосование")
	ErrVotebanClosed      = errors.New("voteban: голосование закрыто")
	ErrVoteSelfTarget     = errors.New("voteban: цель голосования не может голосовать")
)

// --- Global bans ---

// UpsertGlobalBan создаёт/обновляет запись глобального бана.
func (s *ModerationService) UpsertGlobalBan(userID, bannedBy int64, reason *string, duration time.Duration) (*models.GlobalBan, error) {
	b := &models.GlobalBan{
		UserID:   userID,
		BannedBy: bannedBy,
		Reason:   reason,
	}
	if duration > 0 {
		t := time.Now().Add(duration)
		b.ExpiresAt = &t
	}
	if err := s.repo.UpsertGlobalBan(b); err != nil {
		return nil, err
	}
	return b, nil
}

// GetGlobalBan возвращает запись или nil. Запись с истёкшим expires_at
// возвращается, чтобы вызывающий мог сам решить (например, считать неактивной).
func (s *ModerationService) GetGlobalBan(userID int64) (*models.GlobalBan, error) {
	return s.repo.GetGlobalBan(userID)
}

// IsGloballyBanned — true, если запись существует и активна на now.
func (s *ModerationService) IsGloballyBanned(userID int64) (bool, *models.GlobalBan, error) {
	b, err := s.repo.GetGlobalBan(userID)
	if err != nil || b == nil {
		return false, nil, err
	}
	return b.IsActive(time.Now()), b, nil
}

// DeleteGlobalBan снимает глобальный бан.
func (s *ModerationService) DeleteGlobalBan(userID int64) error {
	return s.repo.DeleteGlobalBan(userID)
}

// ListActiveGlobalBans — список действующих банов для /globalbans.
func (s *ModerationService) ListActiveGlobalBans() ([]models.GlobalBan, error) {
	return s.repo.ListActiveGlobalBans(time.Now())
}

// KnownChatIDs возвращает все известные боту чаты (subscription + tracked).
func (s *ModerationService) KnownChatIDs() ([]int64, error) {
	return s.repo.KnownChatIDs()
}
