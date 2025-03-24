package repository

import (
	"fmt"
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

// Изменяем с type alias на новый тип
type MemberRepository struct {
	BaseRepository[models.Member]
}

func NewMemberRepository() *MemberRepository {
	return &MemberRepository{
		BaseRepository: NewBaseRepository(database.DB, &models.Member{}),
	}
}

// GetMemberByTelegram находит пользователя по его Telegram
func (r *MemberRepository) GetMemberByTelegram(telegram string) (*models.Member, error) {
	var member models.Member
	result := database.DB.Where("tg = ?", telegram).First(&member)
	if result.Error != nil {
		return nil, result.Error
	}
	return &member, nil
}

// Search выполняет поиск участников с пагинацией и проверкой на статус ментора
func (r *MemberRepository) Search(limit *int, offset *int) ([]models.MemberModel, int, error) {
	var count int64

	// Сначала считаем общее количество всех записей
	if err := database.DB.Model(&models.Member{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Создаем SQL-запрос с LEFT JOIN для получения всех участников и информации о том, являются ли они менторами
	query := `
		SELECT m.id, m.tg, m.name, 
		       CASE WHEN mt.id IS NOT NULL THEN true ELSE false END as is_mentor
		FROM members m
		LEFT JOIN mentors mt ON m.id = mt."memberId"
	`

	// Добавляем LIMIT и OFFSET, если они переданы
	if limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *limit)
	}
	if offset != nil {
		query += fmt.Sprintf(" OFFSET %d", *offset)
	}

	// Выполняем запрос
	rows, err := database.DB.Raw(query).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Обрабатываем результаты
	var result []models.MemberModel
	for rows.Next() {
		var id int64
		var tg, name string
		var isMentor bool
		if err := rows.Scan(&id, &tg, &name, &isMentor); err != nil {
			return nil, 0, err
		}
		result = append(result, models.MemberModel{
			Id:       id,
			Tg:       tg,
			Name:     name,
			IsMentor: isMentor,
		})
	}

	return result, int(count), nil
}

// GetById получает участника по ID с проверкой на статус ментора
func (r *MemberRepository) GetById(id int64) (*models.MemberModel, error) {
	var member models.Member
	if err := database.DB.First(&member, id).Error; err != nil {
		return nil, err
	}

	// Проверяем, является ли участник ментором
	var count int64
	if err := database.DB.Model(&models.MentorDbShortModel{}).Where("memberId = ?", id).Count(&count).Error; err != nil {
		return nil, err
	}

	result := &models.MemberModel{
		Id:       member.Id,
		Tg:       member.Tg,
		Name:     member.Name,
		IsMentor: count > 0, // Если есть хотя бы одна запись, то участник является ментором
	}

	return result, nil
}
