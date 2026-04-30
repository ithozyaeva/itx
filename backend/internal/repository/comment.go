package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentRepository struct{}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{}
}

func (r *CommentRepository) List(entityType models.CommentEntityType, entityID, viewerID int64, includeHidden bool, limit, offset int) ([]models.Comment, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	countQ := database.DB.Model(&models.Comment{}).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID)
	if !includeHidden {
		countQ = countQ.Where("is_hidden = ?", false)
	}
	var total int64
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.Comment
	q := database.DB.Preload("Author").
		Where("entity_type = ? AND entity_id = ?", entityType, entityID)
	if !includeHidden {
		q = q.Where("is_hidden = ?", false)
	}
	if err := q.Order("created_at ASC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	if viewerID > 0 && len(items) > 0 {
		ids := make([]int64, len(items))
		for i, c := range items {
			ids[i] = c.Id
		}
		liked, err := r.fetchLikedByViewer(ids, viewerID)
		if err != nil {
			return nil, 0, err
		}
		for i := range items {
			items[i].Liked = liked[items[i].Id]
		}
	}
	return items, total, nil
}

func (r *CommentRepository) GetByID(id int64) (*models.Comment, error) {
	var c models.Comment
	if err := database.DB.Preload("Author").First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CommentRepository) Create(c *models.Comment) (*models.Comment, error) {
	if err := database.DB.Create(c).Error; err != nil {
		return nil, err
	}
	return r.GetByID(c.Id)
}

func (r *CommentRepository) Update(id int64, body string) error {
	return database.DB.Model(&models.Comment{}).
		Where("id = ?", id).
		Update("body", body).Error
}

func (r *CommentRepository) Delete(id int64) error {
	return database.DB.Delete(&models.Comment{}, id).Error
}

func (r *CommentRepository) SetHidden(id int64, hidden bool) error {
	return database.DB.Model(&models.Comment{}).
		Where("id = ?", id).
		Update("is_hidden", hidden).Error
}

// ToggleLike — атомарный переключатель лайка коммента. ON CONFLICT
// DO NOTHING защищает от гонки параллельных POST'ов; ветка по
// RowsAffected определяет финальное состояние. Триггеры в миграции
// держат comments.likes_count в синхроне.
func (r *CommentRepository) ToggleLike(commentID, memberID int64) (bool, int, error) {
	var liked bool
	var count int
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Clauses(clause.OnConflict{DoNothing: true}).
			Create(&models.CommentLike{CommentId: commentID, MemberId: memberID})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 1 {
			liked = true
		} else {
			if dErr := tx.Where("comment_id = ? AND member_id = ?", commentID, memberID).
				Delete(&models.CommentLike{}).Error; dErr != nil {
				return dErr
			}
			liked = false
		}

		var c models.Comment
		if err := tx.Select("likes_count").First(&c, commentID).Error; err != nil {
			return err
		}
		count = c.LikesCount
		return nil
	})
	return liked, count, err
}

func (r *CommentRepository) fetchLikedByViewer(ids []int64, viewerID int64) (map[int64]bool, error) {
	m := make(map[int64]bool, len(ids))
	if len(ids) == 0 {
		return m, nil
	}
	var rows []models.CommentLike
	if err := database.DB.Where("comment_id IN ? AND member_id = ?", ids, viewerID).Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, r := range rows {
		m[r.CommentId] = true
	}
	return m, nil
}
