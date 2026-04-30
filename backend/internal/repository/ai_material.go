package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"

	"gorm.io/gorm"
)

type AIMaterialRepository struct{}

func NewAIMaterialRepository() *AIMaterialRepository {
	return &AIMaterialRepository{}
}

type AIMaterialFilter struct {
	Kind        string
	Tag         string
	Query       string
	AuthorID    int64
	Bookmarked  bool // если true — фильтрует только материалы, на которые есть закладка ViewerID
	IncludeHidden bool
	ViewerID    int64 // member, для которого считаем liked/bookmarked
	Sort        string // "new" | "popular"
	Limit       int
	Offset      int
}

func (r *AIMaterialRepository) Search(f AIMaterialFilter) ([]models.AIMaterial, int64, error) {
	q := database.DB.Model(&models.AIMaterial{})

	if !f.IncludeHidden {
		q = q.Where("is_hidden = ?", false)
	}
	if f.Kind != "" {
		q = q.Where("material_kind = ?", f.Kind)
	}
	if f.AuthorID > 0 {
		q = q.Where("author_id = ?", f.AuthorID)
	}
	if f.Query != "" {
		like := "%" + f.Query + "%"
		q = q.Where("title ILIKE ? OR summary ILIKE ?", like, like)
	}
	if f.Tag != "" {
		q = q.Where("EXISTS (SELECT 1 FROM ai_material_tags t WHERE t.material_id = ai_materials.id AND t.tag = ?)", f.Tag)
	}
	if f.Bookmarked && f.ViewerID > 0 {
		q = q.Where("EXISTS (SELECT 1 FROM ai_material_bookmarks b WHERE b.material_id = ai_materials.id AND b.member_id = ?)", f.ViewerID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := f.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	findQ := q.Preload("Author")
	switch f.Sort {
	case "popular":
		findQ = findQ.Order("likes_count DESC").Order("created_at DESC")
	default:
		findQ = findQ.Order("created_at DESC")
	}

	var items []models.AIMaterial
	if err := findQ.Limit(limit).Offset(f.Offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	if len(items) == 0 {
		return items, total, nil
	}

	ids := make([]int64, len(items))
	for i, it := range items {
		ids[i] = it.Id
	}

	tagsByMaterial, err := r.fetchTagsForMaterials(ids)
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		items[i].Tags = tagsByMaterial[items[i].Id]
		if items[i].Tags == nil {
			items[i].Tags = []string{}
		}
	}

	if f.ViewerID > 0 {
		liked, err := r.fetchLikedByViewer(ids, f.ViewerID)
		if err != nil {
			return nil, 0, err
		}
		bookmarked, err := r.fetchBookmarkedByViewer(ids, f.ViewerID)
		if err != nil {
			return nil, 0, err
		}
		for i := range items {
			items[i].Liked = liked[items[i].Id]
			items[i].Bookmarked = bookmarked[items[i].Id]
		}
	}

	return items, total, nil
}

func (r *AIMaterialRepository) GetByID(id int64, viewerID int64) (*models.AIMaterial, error) {
	var item models.AIMaterial
	if err := database.DB.Preload("Author").First(&item, id).Error; err != nil {
		return nil, err
	}

	tags, err := r.fetchTagsForMaterials([]int64{id})
	if err != nil {
		return nil, err
	}
	item.Tags = tags[id]
	if item.Tags == nil {
		item.Tags = []string{}
	}

	if viewerID > 0 {
		liked, err := r.fetchLikedByViewer([]int64{id}, viewerID)
		if err != nil {
			return nil, err
		}
		bookmarked, err := r.fetchBookmarkedByViewer([]int64{id}, viewerID)
		if err != nil {
			return nil, err
		}
		item.Liked = liked[id]
		item.Bookmarked = bookmarked[id]
	}

	return &item, nil
}

func (r *AIMaterialRepository) Create(item *models.AIMaterial, tags []string) (*models.AIMaterial, error) {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		if len(tags) > 0 {
			rows := make([]models.AIMaterialTag, 0, len(tags))
			for _, t := range tags {
				rows = append(rows, models.AIMaterialTag{MaterialId: item.Id, Tag: t})
			}
			if err := tx.Create(&rows).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetByID(item.Id, 0)
}

func (r *AIMaterialRepository) Update(id int64, updates map[string]interface{}, tags *[]string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if len(updates) > 0 {
			if err := tx.Model(&models.AIMaterial{}).Where("id = ?", id).Updates(updates).Error; err != nil {
				return err
			}
		}
		if tags != nil {
			if err := tx.Where("material_id = ?", id).Delete(&models.AIMaterialTag{}).Error; err != nil {
				return err
			}
			if len(*tags) > 0 {
				rows := make([]models.AIMaterialTag, 0, len(*tags))
				for _, t := range *tags {
					rows = append(rows, models.AIMaterialTag{MaterialId: id, Tag: t})
				}
				if err := tx.Create(&rows).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *AIMaterialRepository) Delete(id int64) error {
	// CASCADE на FK уберёт теги/лайки/закладки/комменты автоматически.
	return database.DB.Delete(&models.AIMaterial{}, id).Error
}

func (r *AIMaterialRepository) SetHidden(id int64, hidden bool) error {
	return database.DB.Model(&models.AIMaterial{}).Where("id = ?", id).Update("is_hidden", hidden).Error
}

// TopTags возвращает популярные теги для autocomplete; q — необязательный
// фильтр-префикс. Сортировка по частоте использования.
func (r *AIMaterialRepository) TopTags(q string, limit int) ([]string, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	type row struct {
		Tag string
	}
	var rows []row
	query := database.DB.
		Table("ai_material_tags").
		Select("tag").
		Group("tag").
		Order("COUNT(*) DESC")
	if q != "" {
		query = query.Where("tag ILIKE ?", q+"%")
	}
	if err := query.Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	tags := make([]string, len(rows))
	for i, r := range rows {
		tags[i] = r.Tag
	}
	return tags, nil
}

func (r *AIMaterialRepository) fetchTagsForMaterials(ids []int64) (map[int64][]string, error) {
	if len(ids) == 0 {
		return map[int64][]string{}, nil
	}
	var tags []models.AIMaterialTag
	if err := database.DB.Where("material_id IN ?", ids).Order("tag").Find(&tags).Error; err != nil {
		return nil, err
	}
	m := make(map[int64][]string, len(ids))
	for _, t := range tags {
		m[t.MaterialId] = append(m[t.MaterialId], t.Tag)
	}
	return m, nil
}

func (r *AIMaterialRepository) fetchLikedByViewer(ids []int64, viewerID int64) (map[int64]bool, error) {
	m := make(map[int64]bool, len(ids))
	if len(ids) == 0 {
		return m, nil
	}
	var rows []models.AIMaterialLike
	if err := database.DB.Where("material_id IN ? AND member_id = ?", ids, viewerID).Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, r := range rows {
		m[r.MaterialId] = true
	}
	return m, nil
}

func (r *AIMaterialRepository) fetchBookmarkedByViewer(ids []int64, viewerID int64) (map[int64]bool, error) {
	m := make(map[int64]bool, len(ids))
	if len(ids) == 0 {
		return m, nil
	}
	var rows []models.AIMaterialBookmark
	if err := database.DB.Where("material_id IN ? AND member_id = ?", ids, viewerID).Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, r := range rows {
		m[r.MaterialId] = true
	}
	return m, nil
}
