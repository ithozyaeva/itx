package repository

import (
	"errors"
	"fmt"

	"ithozyeva/database"
	"ithozyeva/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
)

// ErrUsernameTaken — попытка записать username, который уже занят другим
// участником. Хендлер должен превратить это в 409 Conflict.
var ErrUsernameTaken = errors.New("username already taken")

// usernameUniqueIndex — имя partial UNIQUE индекса из миграции
// 20260506000000_dedupe_and_unique_username.sql. Сверяемся именно по
// constraint_name, чтобы 23505 от другой колонки случайно не превратился
// в «никнейм занят».
const usernameUniqueIndex = "members_username_unique"

func isUsernameUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	return pgErr.Code == "23505" && pgErr.ConstraintName == usernameUniqueIndex
}

// Изменяем с type alias на новый тип
type MemberRepositoryInterface interface {
	BaseRepository[models.Member]

	GetByTelegramID(telegramID int64) (*models.Member, error)
	HasRole(memberID int64, role models.Role) bool
	HasPermission(memberID int64, permission models.Permission) bool
	GetMemberPermissions(memberID int64) ([]models.Permission, error)
	GetAllPermissions() ([]models.Permission, error)
}
type MemberRepository struct {
	BaseRepository[models.Member]
}

func NewMemberRepository() *MemberRepository {
	return &MemberRepository{
		BaseRepository: NewBaseRepository(database.DB, &models.Member{}),
	}
}

func (r *MemberRepository) GetMemberByTelegram(telegram string) (*models.Member, error) {
	if telegram == "" {
		return nil, fmt.Errorf("empty username")
	}
	var member models.Member
	result := database.DB.Preload("MemberRoles").
		Where("LOWER(username) = LOWER(?) AND username <> ''", telegram).
		Order("id DESC").
		First(&member)
	if result.Error != nil {
		return nil, result.Error
	}
	return &member, nil
}

// ReleaseUsername освобождает поле username у всех записей с таким значением,
// кроме указанной (keepID). Используется когда новый владелец логинится через
// Telegram с username, который в БД ещё «висит» за чужим аккаунтом.
func (r *MemberRepository) ReleaseUsername(username string, keepID int64) error {
	if username == "" {
		return nil
	}
	return database.DB.Model(&models.Member{}).
		Where("LOWER(username) = LOWER(?) AND id <> ?", username, keepID).
		Update("username", "").Error
}

func (r *MemberRepository) GetByTelegramID(telegramID int64) (*models.Member, error) {
	entity := new(models.Member)
	if err := database.DB.Preload("MemberRoles").Where("telegram_id = ?", telegramID).First(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

// GetById получает участника по ID с проверкой на статус ментора
func (r *MemberRepository) GetById(id int64) (*models.Member, error) {
	var member models.Member
	if err := database.DB.Preload("MemberRoles").First(&member, id).Error; err != nil {
		return nil, err
	}
	result := &models.Member{
		Id:         member.Id,
		Username:   member.Username,
		FirstName:  member.FirstName,
		TelegramID: member.TelegramID,
		LastName:   member.LastName,
		Bio:        member.Bio,
		AvatarURL:  member.AvatarURL,
		Roles:      member.GetRoleStrings(),
		Birthday:   member.Birthday,
		CreatedAt:  member.CreatedAt,
	}

	return result, nil
}

func (r *MemberRepository) Create(member *models.Member) (*models.Member, error) {
	result := database.DB.Model(&models.Member{}).
		Create(&member)

	if result.Error != nil {
		if isUsernameUniqueViolation(result.Error) {
			return nil, ErrUsernameTaken
		}
		return nil, result.Error
	}

	member.SetRoleStrings(member.Roles, member.Id)
	database.DB.Model(member).Association("MemberRoles").Replace(member.MemberRoles)

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("member not found")
	}

	return member, nil
}

func (r *MemberRepository) Update(member *models.Member) (*models.Member, error) {
	result := database.DB.Model(&models.Member{}).
		Where("id = ?", member.Id).
		Updates(map[string]interface{}{
			"birthday":   member.Birthday,
			"first_name": member.FirstName,
			"last_name":  member.LastName,
			"bio":        member.Bio,
			"grade":      member.Grade,
			"company":    member.Company,
			"avatar_url": member.AvatarURL,
			"username":   member.Username,
		})

	if result.Error != nil {
		if isUsernameUniqueViolation(result.Error) {
			return nil, ErrUsernameTaken
		}
		return nil, result.Error
	}

	member.SetRoleStrings(member.Roles, member.Id)
	database.DB.Where("member_id = ? AND role NOT IN ?", member.Id, member.Roles).Delete(&models.MemberRole{})

	database.DB.Model(member).Association("MemberRoles").Replace(member.MemberRoles)

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("member not found")
	}

	return member, nil
}

func (r *MemberRepository) GetTodayBirthdays() ([]string, error) {
	// Сравниваем по МСК-дате, а не по UTC: бизнес живёт в Москве,
	// а DSN устанавливает session TZ='UTC'. Раньше SQL использовал
	// CURRENT_DATE (UTC) — для пользователя с днём рождения 1 января
	// в 00:00–03:00 MSK 1 января UTC ещё показывал 31 декабря, и поздравление
	// проскакивало. Plain conversion: AT TIME ZONE сначала в UTC (если
	// колонка timestamp without time zone, считается UTC), потом в Moscow.
	query := `
		SELECT
			username
		FROM members
		WHERE
			role = ?
			AND
			DATE_PART('day', birthday) = DATE_PART('day', (CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow')::date)
			AND
			DATE_PART('month', birthday) = DATE_PART('month', (CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow')::date)
	`

	rows, err := database.DB.Raw(query, models.MemberRoleSubscriber).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}
	return usernames, nil
}

func (r *MemberRepository) Search(limit *int, offset *int, filter *SearchFilter, order *Order) ([]models.Member, int64, error) {
	var members []models.Member
	var count int64

	query := database.DB.Model(&models.Member{})

	if filter != nil {
		for key, value := range *filter {
			query = query.Where(key, value)
		}
	}

	// Count the filtered results
	countQuery := database.DB.Model(&models.Member{})
	if filter != nil {
		for key, value := range *filter {
			countQuery = countQuery.Where(key, value)
		}
	}

	if err := countQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if order != nil {
		query = query.Order(fmt.Sprintf("\"%s\" %s", order.ColumnBy, order.Order))
	} else {
		query = query.Order("id ASC")
	}

	if limit != nil {
		query = query.Limit(*limit)
	}

	if offset != nil {
		query = query.Offset(*offset)
	}

	if err := query.Preload("MemberRoles").Find(&members).Error; err != nil {
		return nil, 0, err
	}

	return members, count, nil
}

func (r *MemberRepository) HasRole(memberID int64, role models.Role) bool {
	var member models.Member
	if err := database.DB.Preload("Roles").First(&member, memberID).Error; err != nil {
		return false
	}

	for _, r := range member.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (r *MemberRepository) HasPermission(memberID int64, permission models.Permission) bool {
	// Get member roles using the member_roles table
	var roleNames []string
	err := database.DB.Table("member_roles").
		Select("role").
		Where("member_id = ?", memberID).
		Pluck("member_roles.role", &roleNames).Error

	if err != nil || len(roleNames) == 0 {
		return false
	}

	var count int64
	err = database.DB.Table("role_permissions").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role IN ? AND permissions.name = ?", roleNames, string(permission)).
		Count(&count).Error

	if err != nil {
		return false
	}

	return count > 0
}

func (r *MemberRepository) GetMemberPermissions(memberID int64) ([]models.Permission, error) {
	var permissions []models.Permission

	// Get member roles using the member_roles table
	var roleNames []string
	err := database.DB.Table("member_roles").
		Select("role").
		Where("member_id = ?", memberID).
		Pluck("member_roles.role", &roleNames).Error

	if err != nil {
		return nil, err
	}

	// Get permissions for these roles
	err = database.DB.Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role IN ?", roleNames).
		Pluck("permissions.name", &permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *MemberRepository) GetAllPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	err := database.DB.Table("permissions").Pluck("permissions.name", &permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetMembersByRole returns all members that have the given role.
func (r *MemberRepository) GetMembersByRole(role models.Role) ([]models.Member, error) {
	var members []models.Member
	err := database.DB.
		Table("members").
		Joins("INNER JOIN member_roles ON members.id = member_roles.member_id").
		Where("member_roles.role = ?", role).
		Preload("MemberRoles").
		Find(&members).Error
	return members, err
}

// GetSubscribedMembersWithTelegram получает всех подписанных пользователей (SUBSCRIBER) с telegram_id
func (r *MemberRepository) GetSubscribedMembersWithTelegram() ([]models.Member, error) {
	var members []models.Member
	
	err := database.DB.
		Table("members").
		Joins("INNER JOIN member_roles ON members.id = member_roles.member_id").
		Where("member_roles.role = ?", models.MemberRoleSubscriber).
		Where("members.telegram_id IS NOT NULL AND members.telegram_id != 0").
		Preload("MemberRoles").
		Find(&members).Error
	
	return members, err
}
