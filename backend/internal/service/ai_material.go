package service

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm"

	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

// Sentinel-ошибки сервиса. Хендлер мапит их на HTTP-коды:
// ErrNotFound → 404, ErrForbidden → 403, всё остальное (валидация и пр.) → 400.
var (
	ErrAIMaterialNotFound  = errors.New("материал не найден")
	ErrAIMaterialForbidden = errors.New("недостаточно прав")
	ErrAIMaterialCommentNotFound = errors.New("комментарий не найден")
)

type AIMaterialService struct {
	repo *repository.AIMaterialRepository
}

func NewAIMaterialService() *AIMaterialService {
	return &AIMaterialService{repo: repository.NewAIMaterialRepository()}
}

func (s *AIMaterialService) Search(f repository.AIMaterialFilter) ([]models.AIMaterial, int64, error) {
	return s.repo.Search(f)
}

// GetByID — единая точка проверки видимости материала. Скрытые материалы
// возвращаются автору и админу, остальным отдаётся ErrNotFound (без утечки
// факта существования). Все интерактивные эндпоинты должны идти через этот
// метод, чтобы лайки/закладки/комменты не лились на скрытое.
func (s *AIMaterialService) GetByID(id, viewerID int64, isAdmin bool) (*models.AIMaterial, error) {
	m, err := s.repo.GetByID(id, viewerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAIMaterialNotFound
		}
		return nil, err
	}
	if m.IsHidden && !isAdmin && m.AuthorId != viewerID {
		return nil, ErrAIMaterialNotFound
	}
	return m, nil
}

func (s *AIMaterialService) Create(req *models.CreateAIMaterialRequest, authorID int64) (*models.AIMaterial, error) {
	normalized, tags, err := s.validateAndNormalize(req)
	if err != nil {
		return nil, err
	}

	item := &models.AIMaterial{
		AuthorId:     authorID,
		Title:        normalized.Title,
		Summary:      normalized.Summary,
		ContentType:  normalized.ContentType,
		MaterialKind: normalized.MaterialKind,
		PromptBody:   normalized.PromptBody,
		ExternalURL:  normalized.ExternalURL,
		AgentConfig:  normalized.AgentConfig,
	}

	created, err := s.repo.Create(item, tags)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *AIMaterialService) Update(id int64, req *models.UpdateAIMaterialRequest, memberID int64, isAdmin bool) (*models.AIMaterial, error) {
	existing, err := s.repo.GetByID(id, 0)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAIMaterialNotFound
		}
		return nil, err
	}
	if !isAdmin && existing.AuthorId != memberID {
		return nil, ErrAIMaterialForbidden
	}

	normalized, tags, err := s.validateAndNormalize(req)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"title":         normalized.Title,
		"summary":       normalized.Summary,
		"content_type":  normalized.ContentType,
		"material_kind": normalized.MaterialKind,
		"prompt_body":   normalized.PromptBody,
		"external_url":  normalized.ExternalURL,
		"agent_config":  normalized.AgentConfig,
	}

	if err := s.repo.Update(id, updates, &tags); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id, memberID)
}

func (s *AIMaterialService) Delete(id, memberID int64, isAdmin bool) error {
	existing, err := s.repo.GetByID(id, 0)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAIMaterialNotFound
		}
		return err
	}
	if !isAdmin && existing.AuthorId != memberID {
		return ErrAIMaterialForbidden
	}
	return s.repo.Delete(id)
}

// SetHidden — admin-only мягкое скрытие материала из листинга.
func (s *AIMaterialService) SetHidden(id int64, hidden, isAdmin bool) error {
	if !isAdmin {
		return ErrAIMaterialForbidden
	}
	if _, err := s.repo.GetByID(id, 0); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAIMaterialNotFound
		}
		return err
	}
	return s.repo.SetHidden(id, hidden)
}

func (s *AIMaterialService) TopTags(q string, limit int) ([]string, error) {
	return s.repo.TopTags(strings.ToLower(strings.TrimSpace(q)), limit)
}

func (s *AIMaterialService) ToggleLike(materialID, memberID int64, isAdmin bool) (bool, int, error) {
	if _, err := s.GetByID(materialID, memberID, isAdmin); err != nil {
		return false, 0, err
	}
	return s.repo.ToggleLike(materialID, memberID)
}

func (s *AIMaterialService) ToggleBookmark(materialID, memberID int64, isAdmin bool) (bool, int, error) {
	if _, err := s.GetByID(materialID, memberID, isAdmin); err != nil {
		return false, 0, err
	}
	return s.repo.ToggleBookmark(materialID, memberID)
}

const (
	AIMaterialCommentMinLen = 1
	AIMaterialCommentMaxLen = 4_000
)

func (s *AIMaterialService) ListComments(materialID, viewerID int64, isAdmin bool) ([]models.AIMaterialComment, error) {
	if _, err := s.GetByID(materialID, viewerID, isAdmin); err != nil {
		return nil, err
	}
	// Скрытые комментарии видит только админ.
	return s.repo.ListComments(materialID, viewerID, isAdmin)
}

// ToggleCommentLike — лайк/анлайк комментария. Доступно только если родительский
// материал виден пользователю (скрытое нельзя «накрутить» по прямому ID комментария).
func (s *AIMaterialService) ToggleCommentLike(commentID, memberID int64, isAdmin bool) (bool, int, error) {
	c, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, 0, ErrAIMaterialCommentNotFound
		}
		return false, 0, err
	}
	if _, err := s.GetByID(c.MaterialId, memberID, isAdmin); err != nil {
		return false, 0, err
	}
	return s.repo.ToggleCommentLike(commentID, memberID)
}

func (s *AIMaterialService) CreateComment(materialID, authorID int64, body string, isAdmin bool) (*models.AIMaterialComment, error) {
	if _, err := s.GetByID(materialID, authorID, isAdmin); err != nil {
		return nil, err
	}
	body = strings.TrimSpace(body)
	if l := utf8.RuneCountInString(body); l < AIMaterialCommentMinLen || l > AIMaterialCommentMaxLen {
		return nil, fmt.Errorf("длина комментария должна быть от %d до %d символов",
			AIMaterialCommentMinLen, AIMaterialCommentMaxLen)
	}
	return s.repo.CreateComment(&models.AIMaterialComment{
		MaterialId: materialID,
		AuthorId:   authorID,
		Body:       body,
	})
}

func (s *AIMaterialService) UpdateComment(commentID, memberID int64, body string, isAdmin bool) (*models.AIMaterialComment, error) {
	existing, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAIMaterialCommentNotFound
		}
		return nil, err
	}
	if !isAdmin && existing.AuthorId != memberID {
		return nil, ErrAIMaterialForbidden
	}
	body = strings.TrimSpace(body)
	if l := utf8.RuneCountInString(body); l < AIMaterialCommentMinLen || l > AIMaterialCommentMaxLen {
		return nil, fmt.Errorf("длина комментария должна быть от %d до %d символов",
			AIMaterialCommentMinLen, AIMaterialCommentMaxLen)
	}
	if err := s.repo.UpdateComment(commentID, body); err != nil {
		return nil, err
	}
	return s.repo.GetCommentByID(commentID)
}

func (s *AIMaterialService) DeleteComment(commentID, memberID int64, isAdmin bool) error {
	existing, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAIMaterialCommentNotFound
		}
		return err
	}
	if !isAdmin && existing.AuthorId != memberID {
		return ErrAIMaterialForbidden
	}
	return s.repo.DeleteComment(commentID)
}

// SetCommentHidden — admin-only.
func (s *AIMaterialService) SetCommentHidden(commentID int64, hidden, isAdmin bool) error {
	if !isAdmin {
		return ErrAIMaterialForbidden
	}
	if _, err := s.repo.GetCommentByID(commentID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAIMaterialCommentNotFound
		}
		return err
	}
	return s.repo.SetCommentHidden(commentID, hidden)
}

func (s *AIMaterialService) validateAndNormalize(req *models.CreateAIMaterialRequest) (*models.CreateAIMaterialRequest, []string, error) {
	out := &models.CreateAIMaterialRequest{
		Title:        strings.TrimSpace(req.Title),
		Summary:      strings.TrimSpace(req.Summary),
		ContentType:  models.AIMaterialContentType(strings.ToLower(strings.TrimSpace(string(req.ContentType)))),
		MaterialKind: models.AIMaterialKind(strings.ToLower(strings.TrimSpace(string(req.MaterialKind)))),
		PromptBody:   strings.TrimSpace(req.PromptBody),
		ExternalURL:  strings.TrimSpace(req.ExternalURL),
		AgentConfig:  strings.TrimSpace(req.AgentConfig),
	}

	titleLen := utf8.RuneCountInString(out.Title)
	if titleLen < models.AIMaterialMinTitleLen || titleLen > models.AIMaterialMaxTitleLen {
		return nil, nil, fmt.Errorf("длина названия должна быть от %d до %d символов", models.AIMaterialMinTitleLen, models.AIMaterialMaxTitleLen)
	}
	summaryLen := utf8.RuneCountInString(out.Summary)
	if summaryLen < models.AIMaterialMinSummaryLen || summaryLen > models.AIMaterialMaxSummaryLen {
		return nil, nil, fmt.Errorf("длина описания должна быть от %d до %d символов", models.AIMaterialMinSummaryLen, models.AIMaterialMaxSummaryLen)
	}

	if !models.IsValidAIMaterialContentType(out.ContentType) {
		return nil, nil, errors.New("некорректный тип контента")
	}
	if !models.IsValidAIMaterialKind(out.MaterialKind) {
		return nil, nil, errors.New("некорректная категория материала")
	}

	switch out.ContentType {
	case models.AIMaterialContentTypePrompt:
		if out.PromptBody == "" {
			return nil, nil, errors.New("содержимое промта обязательно")
		}
		if utf8.RuneCountInString(out.PromptBody) > models.AIMaterialMaxPromptBody {
			return nil, nil, fmt.Errorf("содержимое промта не должно превышать %d символов", models.AIMaterialMaxPromptBody)
		}
		out.ExternalURL = ""
		out.AgentConfig = ""
	case models.AIMaterialContentTypeLink:
		if out.ExternalURL == "" {
			return nil, nil, errors.New("ссылка обязательна")
		}
		if len(out.ExternalURL) > models.AIMaterialMaxURLLen {
			return nil, nil, fmt.Errorf("длина ссылки не должна превышать %d символов", models.AIMaterialMaxURLLen)
		}
		if !isValidHTTPURL(out.ExternalURL) {
			return nil, nil, errors.New("ссылка должна начинаться с http:// или https://")
		}
		out.PromptBody = ""
		out.AgentConfig = ""
	case models.AIMaterialContentTypeAgent:
		if out.AgentConfig == "" {
			return nil, nil, errors.New("конфиг агента обязателен")
		}
		if utf8.RuneCountInString(out.AgentConfig) > models.AIMaterialMaxAgentConfig {
			return nil, nil, fmt.Errorf("конфиг агента не должен превышать %d символов", models.AIMaterialMaxAgentConfig)
		}
		out.PromptBody = ""
		out.ExternalURL = ""
	}

	tags, err := normalizeTags(req.Tags)
	if err != nil {
		return nil, nil, err
	}

	return out, tags, nil
}

func normalizeTags(in []string) ([]string, error) {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, raw := range in {
		t := strings.ToLower(strings.TrimSpace(raw))
		if t == "" {
			continue
		}
		if utf8.RuneCountInString(t) > models.AIMaterialMaxTagLen {
			return nil, fmt.Errorf("длина тега не должна превышать %d символов", models.AIMaterialMaxTagLen)
		}
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
		if len(out) >= models.AIMaterialMaxTags {
			break
		}
	}
	return out, nil
}

func isValidHTTPURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return u.Host != ""
}
