package service

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm"

	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

// Sentinel-ошибки сервиса комментариев.
var (
	ErrCommentNotFound  = errors.New("комментарий не найден")
	ErrCommentForbidden = errors.New("недостаточно прав")
	// ErrEntityNotFound возвращается visibility-checker'ом, когда родительская
	// сущность не существует или скрыта от viewer'а. Сервис прокидывает
	// её наверх как is — handler мапит на 404.
	ErrEntityNotFound = errors.New("сущность не найдена")
)

// EntityVisibilityChecker проверяет, что viewer имеет право работать с
// комментариями к конкретной сущности (entity видима ему, и подписочные
// требования соблюдены). Возвращает ErrEntityNotFound для скрытого /
// недоступного entity, либо ошибку «нет доступа». Принимает полного
// *models.Member, потому что разным сущностям нужна разная информация:
// AI-материалу — TelegramID для проверки тира и роли, event'у — только
// факт наличия. Регистрируется отдельно для каждого entity_type.
type EntityVisibilityChecker func(entityID int64, member *models.Member) error

type CommentService struct {
	repo     *repository.CommentRepository
	checkers map[models.CommentEntityType]EntityVisibilityChecker
}

func NewCommentService(checkers map[models.CommentEntityType]EntityVisibilityChecker) *CommentService {
	return &CommentService{
		repo:     repository.NewCommentRepository(),
		checkers: checkers,
	}
}

func (s *CommentService) ensureVisible(entityType models.CommentEntityType, entityID int64, member *models.Member) error {
	checker, ok := s.checkers[entityType]
	if !ok {
		return fmt.Errorf("неизвестный тип сущности: %s", entityType)
	}
	return checker(entityID, member)
}

func (s *CommentService) List(entityType models.CommentEntityType, entityID int64, member *models.Member, isAdmin bool, limit, offset int) ([]models.Comment, int64, error) {
	if !models.IsValidCommentEntityType(entityType) {
		return nil, 0, fmt.Errorf("неизвестный тип сущности")
	}
	if err := s.ensureVisible(entityType, entityID, member); err != nil {
		return nil, 0, err
	}
	return s.repo.List(entityType, entityID, member.Id, isAdmin, limit, offset)
}

func (s *CommentService) Create(entityType models.CommentEntityType, entityID int64, member *models.Member, body string) (*models.Comment, error) {
	if !models.IsValidCommentEntityType(entityType) {
		return nil, fmt.Errorf("неизвестный тип сущности")
	}
	if err := s.ensureVisible(entityType, entityID, member); err != nil {
		return nil, err
	}
	body = strings.TrimSpace(body)
	if l := utf8.RuneCountInString(body); l < models.CommentMinLen || l > models.CommentMaxLen {
		return nil, fmt.Errorf("длина комментария должна быть от %d до %d символов",
			models.CommentMinLen, models.CommentMaxLen)
	}
	return s.repo.Create(&models.Comment{
		EntityType: entityType,
		EntityId:   entityID,
		AuthorId:   member.Id,
		Body:       body,
	})
}

func (s *CommentService) Update(commentID, memberID int64, body string, isAdmin bool) (*models.Comment, error) {
	existing, err := s.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}
	if !isAdmin && existing.AuthorId != memberID {
		return nil, ErrCommentForbidden
	}
	body = strings.TrimSpace(body)
	if l := utf8.RuneCountInString(body); l < models.CommentMinLen || l > models.CommentMaxLen {
		return nil, fmt.Errorf("длина комментария должна быть от %d до %d символов",
			models.CommentMinLen, models.CommentMaxLen)
	}
	if err := s.repo.Update(commentID, body); err != nil {
		return nil, err
	}
	return s.repo.GetByID(commentID)
}

func (s *CommentService) Delete(commentID, memberID int64, isAdmin bool) error {
	existing, err := s.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotFound
		}
		return err
	}
	if !isAdmin && existing.AuthorId != memberID {
		return ErrCommentForbidden
	}
	return s.repo.Delete(commentID)
}

// SetHidden — admin-only.
func (s *CommentService) SetHidden(commentID int64, hidden, isAdmin bool) error {
	if !isAdmin {
		return ErrCommentForbidden
	}
	if _, err := s.repo.GetByID(commentID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotFound
		}
		return err
	}
	return s.repo.SetHidden(commentID, hidden)
}

// ToggleLike — лайк/анлайк коммента. Доступно только если родительская
// сущность видна viewer'у (нельзя «накрутить» лайк на коммент скрытой /
// недоступной сущности по прямому commentId).
func (s *CommentService) ToggleLike(commentID int64, member *models.Member) (bool, int, error) {
	c, err := s.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, 0, ErrCommentNotFound
		}
		return false, 0, err
	}
	if err := s.ensureVisible(c.EntityType, c.EntityId, member); err != nil {
		return false, 0, err
	}
	return s.repo.ToggleLike(commentID, member.Id)
}
