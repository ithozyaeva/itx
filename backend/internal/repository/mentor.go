package repository

import (
	"errors"
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"reflect"

	"gorm.io/gorm"
)

// MentorRepository интерфейс для репозитория менторов
type MentorRepositoryInterface interface {
	BaseRepository[models.MentorDbShortModel]
	GetByIdFull(id int64) (*models.MentorDbModel, error)
	GetServices(id int64) ([]models.Service, error)
	AddReviewToService(review *models.ReviewOnService) (*models.ReviewOnService, error)
	CreateWithRelations(mentor *models.MentorDbModel) (*models.MentorDbModel, error)
	UpdateWithRelations(mentor *models.MentorDbModel) (*models.MentorDbModel, error)
	GetAllWithRelations(limit *int, offset *int) ([]models.MentorDbModel, int64, error)
	GetByMemberID(memberId int64) (*models.MentorDbModel, error)
}

// MentorRepository реализует интерфейс MentorRepositoryInterface
type MentorRepository struct {
	BaseRepository[models.MentorDbShortModel]
}

// NewMentorRepository создает новый экземпляр репозитория менторов
func NewMentorRepository() *MentorRepository {
	return &MentorRepository{
		BaseRepository: NewBaseRepository(database.DB, &models.MentorDbShortModel{}),
	}
}

// GetByIdFull получает ментора по ID с полной информацией
func (r *MentorRepository) GetByIdFull(id int64) (*models.MentorDbModel, error) {
	// Начинаем транзакцию для согласованного чтения
	tx := database.DB.Begin()
	defer tx.Rollback()

	// Получаем базовую информацию о менторе
	var mentor models.MentorDbModel
	if err := tx.First(&mentor, id).Error; err != nil {
		return nil, err
	}

	// Загружаем информацию о пользователе
	if err := tx.First(&mentor.Member, mentor.MemberId).Error; err != nil {
		return nil, err
	}

	// Загружаем профессиональные теги
	var mentorTags []models.MentorsTag
	if err := tx.Where("mentor_id = ?", id).Find(&mentorTags).Error; err != nil {
		return nil, err
	}

	// Загружаем теги по их ID
	var tagIds []int64
	for _, mt := range mentorTags {
		tagIds = append(tagIds, mt.TagId)
	}

	if len(tagIds) > 0 {
		if err := tx.Where("id IN ?", tagIds).Find(&mentor.ProfTags).Error; err != nil {
			return nil, err
		}
	}

	// Загружаем контакты
	if err := tx.Where("\"ownerId\" = ?", id).Find(&mentor.Contacts).Error; err != nil {
		return nil, err
	}

	// Загружаем услуги
	if err := tx.Where("\"ownerId\" = ?", id).Find(&mentor.Services).Error; err != nil {
		return nil, err
	}

	tx.Commit()
	return &mentor, nil
}

func (r *MentorRepository) GetServices(id int64) ([]models.Service, error) {
	var services []models.Service

	if err := database.DB.Model(&models.Service{}).Where("ownerId = ?", id).Find(&services).Error; err != nil {
		return nil, err
	}

	return services, nil
}

// AddReviewToService добавляет отзыв к услуге ментора
func (r *MentorRepository) AddReviewToService(review *models.ReviewOnService) (*models.ReviewOnService, error) {
	if err := database.DB.Create(review).Error; err != nil {
		return nil, err
	}
	return review, nil
}

// CreateWithRelations создает нового ментора со всеми связанными сущностями
func (r *MentorRepository) CreateWithRelations(mentor *models.MentorDbModel) (*models.MentorDbModel, error) {
	// Начинаем транзакцию
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Создаем базовую модель ментора
	mentorDb := &models.MentorDbShortModel{
		MemberId:   mentor.MemberId,
		Occupation: mentor.Occupation,
		Experience: mentor.Experience,
		Order:      mentor.Order,
	}

	if err := tx.Create(&mentorDb).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Важно: выполняем явный commit и начинаем новую транзакцию
	// чтобы гарантировать, что ID ментора сохранен в базе
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Начинаем новую транзакцию для связанных сущностей
	tx = database.DB.Begin()

	// Создаем полную модель для возврата
	result := &models.MentorDbModel{
		Id:         mentorDb.Id,
		MemberId:   mentorDb.MemberId,
		Occupation: mentorDb.Occupation,
		Experience: mentorDb.Experience,
		Order:      mentorDb.Order,
	}

	// Обрабатываем связанные сущности
	var err error
	result.ProfTags, err = r.handleProfTags(tx, mentorDb.Id, mentor.ProfTags)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Проверяем, что ID ментора не равен 0
	if mentorDb.Id == 0 {
		tx.Rollback()
		return nil, errors.New("ID ментора не может быть 0")
	}

	contactsInterface, err := r.handleRelatedEntities(tx, mentorDb.Id, mentor.Contacts, "contacts", "ownerId")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	result.Contacts = contactsInterface.([]models.Contact)

	servicesInterface, err := r.handleRelatedEntities(tx, mentorDb.Id, mentor.Services, "services", "ownerId")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	result.Services = servicesInterface.([]models.Service)

	// Загружаем информацию о пользователе
	var member models.Member
	if err := tx.First(&member, mentorDb.MemberId).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	result.Member = member

	// Фиксируем транзакцию
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateWithRelations обновляет ментора со всеми связанными сущностями
func (r *MentorRepository) UpdateWithRelations(mentor *models.MentorDbModel) (*models.MentorDbModel, error) {
	// Начинаем транзакцию
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Проверяем существование ментора
	var existingMentor, memeberExistenceError = r.GetByMemberID(mentor.MemberId)
	if memeberExistenceError != nil {
		return nil, memeberExistenceError
	}

	// Обновляем базовую модель ментора
	existingMentor.MemberId = mentor.MemberId
	existingMentor.Occupation = mentor.Occupation
	existingMentor.Experience = mentor.Experience
	existingMentor.Order = mentor.Order

	if err := tx.Save(&existingMentor).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Создаем полную модель для возврата
	result := &models.MentorDbModel{
		Id:         existingMentor.Id,
		MemberId:   existingMentor.MemberId,
		Occupation: existingMentor.Occupation,
		Experience: existingMentor.Experience,
		Order:      existingMentor.Order,
		Member:     existingMentor.Member,
	}

	// Обрабатываем связанные сущности
	var err error

	result.ProfTags, err = r.handleProfTags(tx, existingMentor.Id, mentor.ProfTags)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	contactsInterface, err := r.handleRelatedEntities(tx, existingMentor.Id, mentor.Contacts, "contacts", "ownerId")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	result.Contacts = contactsInterface.([]models.Contact)

	servicesInterface, err := r.handleRelatedEntities(tx, existingMentor.Id, mentor.Services, "services", "ownerId")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	result.Services = servicesInterface.([]models.Service)

	// Фиксируем транзакцию
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return result, nil
}

// handleProfTags обрабатывает профессиональные теги
func (r *MentorRepository) handleProfTags(tx *gorm.DB, mentorId int64, tagRequests []models.ProfTag) ([]models.ProfTag, error) {
	var result []models.ProfTag

	// Получаем существующие связи
	var existingLinks []models.MentorsTag
	if err := tx.Where("mentor_id = ?", mentorId).Find(&existingLinks).Error; err != nil {
		return nil, err
	}

	// Создаем карту существующих связей для быстрого поиска
	existingTagsMap := make(map[int64]bool)
	for _, link := range existingLinks {
		existingTagsMap[link.TagId] = true
	}

	// Создаем карту запрошенных тегов
	requestedTagsMap := make(map[int64]bool)

	// Обрабатываем запросы тегов
	for _, tagReq := range tagRequests {
		var tag models.ProfTag

		// Проверяем, существует ли тег
		if tagReq.Id > 0 {
			// Ищем тег по ID
			if err := tx.First(&tag, tagReq.Id).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, err
				}
				// Если тег не найден, создаем новый
				tag = models.ProfTag{Title: tagReq.Title}
				if err := tx.Create(&tag).Error; err != nil {
					return nil, err
				}
			} else {
				// Обновляем существующий тег
				tag.Title = tagReq.Title
				if err := tx.Save(&tag).Error; err != nil {
					return nil, err
				}
			}
		} else {
			// Ищем тег по названию
			if err := tx.Where("title = ?", tagReq.Title).First(&tag).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Создаем новый тег
					tag = models.ProfTag{Title: tagReq.Title}
					if err := tx.Create(&tag).Error; err != nil {
						return nil, err
					}
				} else {
					return nil, err
				}
			}
		}

		// Добавляем тег в результат
		result = append(result, tag)

		// Отмечаем тег как запрошенный
		requestedTagsMap[tag.Id] = true

		// Если связи нет, создаем её
		if !existingTagsMap[tag.Id] {
			mentorTag := models.MentorsTag{
				MentorId: mentorId,
				TagId:    tag.Id,
			}
			if err := tx.Table("mentors_tags").Create(&mentorTag).Error; err != nil {
				return nil, err
			}
		}

		// Удаляем тег из карты существующих, чтобы потом удалить оставшиеся
		delete(existingTagsMap, tag.Id)
	}

	// Удаляем связи, которых нет в запросе
	for tagId := range existingTagsMap {
		// Используем Raw SQL с правильными именами столбцов в кавычках
		if err := tx.Exec(`DELETE FROM "mentors_tags" WHERE "mentor_id" = ? AND "tag_id" = ?`, mentorId, tagId).Error; err != nil {
			return nil, err
		}
	}

	return result, nil
}

// handleRelatedEntities обрабатывает связанные сущности (контакты, услуги)
func (r *MentorRepository) handleRelatedEntities(tx *gorm.DB, mentorId int64, requests interface{}, tableName string, foreignKeyName string) (interface{}, error) {
	// Получаем тип запросов
	reqValue := reflect.ValueOf(requests)
	if reqValue.Kind() != reflect.Slice {
		return nil, errors.New("requests должен быть срезом")
	}

	// Создаем пустой срез для результата
	var resultSlice reflect.Value
	var modelType reflect.Type

	// Определяем тип модели на основе типа запроса
	switch tableName {
	case "contacts":
		resultSlice = reflect.MakeSlice(reflect.TypeOf([]models.Contact{}), 0, reqValue.Len())
		modelType = reflect.TypeOf(models.Contact{})
	case "services":
		resultSlice = reflect.MakeSlice(reflect.TypeOf([]models.Service{}), 0, reqValue.Len())
		modelType = reflect.TypeOf(models.Service{})
	default:
		return nil, errors.New("неизвестный тип таблицы")
	}

	// Получаем существующие записи
	query := tx.Table(tableName).Where("\""+foreignKeyName+"\" = ?", mentorId)
	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Создаем карту существующих записей
	existingMap := make(map[int64]bool)
	for rows.Next() {
		entity := reflect.New(modelType).Interface()
		if err := tx.ScanRows(rows, entity); err != nil {
			return nil, err
		}

		// Получаем ID записи
		idField := reflect.ValueOf(entity).Elem().FieldByName("Id")
		if idField.IsValid() {
			existingMap[idField.Int()] = true
		}
	}

	// Обрабатываем запросы
	for i := 0; i < reqValue.Len(); i++ {
		reqItem := reqValue.Index(i)

		// Получаем ID из запроса
		idField := reqItem.FieldByName("Id")
		id := idField.Int()

		if id > 0 && existingMap[id] {
			// Обновляем существующую запись
			entity := reflect.New(modelType).Interface()
			if err := tx.Table(tableName).First(entity, id).Error; err != nil {
				return nil, err
			}

			// Копируем поля из запроса в сущность
			entityValue := reflect.ValueOf(entity).Elem()
			for j := 0; j < reqItem.NumField(); j++ {
				fieldName := reqItem.Type().Field(j).Name
				if fieldName != "Id" {
					entityField := entityValue.FieldByName(fieldName)
					if entityField.IsValid() && entityField.CanSet() {
						entityField.Set(reqItem.Field(j))
					}
				}
			}

			// Сохраняем обновленную запись
			if err := tx.Save(entity).Error; err != nil {
				return nil, err
			}

			resultSlice = reflect.Append(resultSlice, entityValue)
			delete(existingMap, id)
		} else {
			// Создаем новую запись с явным указанием внешнего ключа
			switch tableName {
			case "contacts":
				contact := models.Contact{
					Type:    int16(reqItem.FieldByName("Type").Int()),
					Link:    reqItem.FieldByName("Link").String(),
					OwnerId: mentorId,
				}
				if err := tx.Create(&contact).Error; err != nil {
					return nil, err
				}
				resultSlice = reflect.Append(resultSlice, reflect.ValueOf(contact))

			case "services":
				service := models.Service{
					Name:    reqItem.FieldByName("Name").String(),
					Price:   int(reqItem.FieldByName("Price").Int()),
					OwnerId: mentorId,
				}
				if err := tx.Create(&service).Error; err != nil {
					return nil, err
				}
				resultSlice = reflect.Append(resultSlice, reflect.ValueOf(service))

			default:
				return nil, errors.New("неизвестный тип таблицы")
			}
		}
	}

	// Удаляем записи, которые не были обновлены
	for id := range existingMap {
		if err := tx.Table(tableName).Delete(reflect.New(modelType).Interface(), id).Error; err != nil {
			return nil, err
		}
	}

	// Возвращаем результат
	return resultSlice.Interface(), nil
}

// GetAllWithRelations получает всех менторов с полной информацией о связях
func (r *MentorRepository) GetAllWithRelations(limit *int, offset *int) ([]models.MentorDbModel, int64, error) {
	// Начинаем транзакцию для согласованного чтения
	tx := database.DB.Begin()
	defer tx.Rollback()

	var mentors []models.MentorDbModel
	var count int64

	// Сначала считаем общее количество всех записей
	if err := tx.Model(&models.MentorDbModel{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Создаем базовый запрос
	query := tx.Model(&models.MentorDbModel{})

	// Применяем limit только если он передан
	if limit != nil {
		query = query.Limit(*limit)
	}

	// Применяем offset только если он передан
	if offset != nil {
		query = query.Offset(*offset)
	}

	query = query.Order("\"order\" ASC")

	// Получаем базовую информацию о менторах
	if err := query.Find(&mentors).Error; err != nil {
		return nil, 0, err
	}

	// Для каждого ментора загружаем связанные данные
	for i := range mentors {
		// Загружаем информацию о пользователе
		if err := tx.First(&mentors[i].Member, mentors[i].MemberId).Error; err != nil {
			return nil, 0, err
		}

		// Загружаем профессиональные теги
		var mentorTags []models.MentorsTag
		if err := tx.Where("mentor_id = ?", mentors[i].Id).Find(&mentorTags).Error; err != nil {
			return nil, 0, err
		}

		// Загружаем теги по их ID
		var tagIds []int64
		for _, mt := range mentorTags {
			tagIds = append(tagIds, mt.TagId)
		}

		if len(tagIds) > 0 {
			if err := tx.Where("id IN ?", tagIds).Find(&mentors[i].ProfTags).Error; err != nil {
				return nil, 0, err
			}
		}

		// Загружаем контакты
		if err := tx.Where("\"ownerId\" = ?", mentors[i].Id).Find(&mentors[i].Contacts).Error; err != nil {
			return nil, 0, err
		}

		// Загружаем услуги
		if err := tx.Where("\"ownerId\" = ?", mentors[i].Id).Find(&mentors[i].Services).Error; err != nil {
			return nil, 0, err
		}
	}

	tx.Commit()
	return mentors, count, nil
}

func (r *MentorRepository) GetByMemberID(memberId int64) (*models.MentorDbModel, error) {
	var entity models.MentorDbModel
	err := database.DB.Model(&models.MentorDbModel{}).
		Where("\"memberId\" = ?", memberId).
		Preload("Member").
		Preload("Member.MemberRoles").
		Preload("ProfTags").
		Preload("Contacts").
		Preload("Services").
		First(&entity).Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}
