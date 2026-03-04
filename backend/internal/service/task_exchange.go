package service

import (
	"errors"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

type TaskExchangeService struct {
	repo *repository.TaskExchangeRepository
}

func NewTaskExchangeService() *TaskExchangeService {
	return &TaskExchangeService{
		repo: repository.NewTaskExchangeRepository(),
	}
}

func (s *TaskExchangeService) Search(status *string, limit, offset int) ([]models.TaskExchange, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.repo.Search(status, limit, offset)
}

func (s *TaskExchangeService) GetById(id int64) (*models.TaskExchange, error) {
	return s.repo.GetById(id)
}

func (s *TaskExchangeService) Create(task *models.TaskExchange) (*models.TaskExchange, error) {
	if task.Title == "" {
		return nil, errors.New("title is required")
	}
	if task.MaxAssignees <= 0 {
		task.MaxAssignees = 1
	}
	task.Status = models.TaskStatusOpen
	return s.repo.Create(task)
}

func (s *TaskExchangeService) Update(id int64, memberId int64, isAdmin bool, req models.UpdateTaskExchangeRequest) (*models.TaskExchange, error) {
	task, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("задание не найдено")
	}
	if !isAdmin && task.CreatorId != memberId {
		return nil, errors.New("только автор или админ может редактировать задание")
	}
	if task.Status == models.TaskStatusDone || task.Status == models.TaskStatusApproved {
		return nil, errors.New("нельзя редактировать задание в этом статусе")
	}

	updates := map[string]interface{}{}
	if req.Title != nil && *req.Title != "" {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.MaxAssignees != nil && *req.MaxAssignees >= 1 {
		count, err := s.repo.GetAssigneesCount(id)
		if err != nil {
			return nil, errors.New("не удалось проверить количество исполнителей")
		}
		if int64(*req.MaxAssignees) < count {
			return nil, errors.New("нельзя уменьшить лимит ниже текущего числа исполнителей")
		}
		updates["max_assignees"] = *req.MaxAssignees
	}

	if len(updates) == 0 {
		return task, nil
	}

	if err := s.repo.Update(id, updates); err != nil {
		return nil, errors.New("не удалось обновить задание")
	}

	return s.repo.GetById(id)
}

func (s *TaskExchangeService) Assign(id int64, memberId int64) (*models.TaskExchange, error) {
	task, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("задание не найдено")
	}
	if task.CreatorId == memberId {
		return nil, errors.New("нельзя взять своё задание")
	}
	if task.Status != models.TaskStatusOpen {
		return nil, errors.New("задание недоступно для взятия")
	}

	isAlready, err := s.repo.IsAssignee(id, memberId)
	if err != nil {
		return nil, errors.New("не удалось проверить назначение")
	}
	if isAlready {
		return nil, errors.New("вы уже являетесь исполнителем")
	}

	count, err := s.repo.GetAssigneesCount(id)
	if err != nil {
		return nil, errors.New("не удалось проверить количество исполнителей")
	}
	if count >= int64(task.MaxAssignees) {
		return nil, errors.New("все слоты исполнителей заняты")
	}

	if err := s.repo.AddAssignee(id, memberId); err != nil {
		return nil, errors.New("не удалось взять задание")
	}

	// Check if all slots are filled now
	if count+1 >= int64(task.MaxAssignees) {
		if err := s.repo.UpdateStatus(id, models.TaskStatusInProgress); err != nil {
			return nil, errors.New("не удалось обновить статус задания")
		}
	}

	return s.repo.GetById(id)
}

func (s *TaskExchangeService) Unassign(id int64, memberId int64) (*models.TaskExchange, error) {
	task, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("задание не найдено")
	}

	isAssigned, err := s.repo.IsAssignee(id, memberId)
	if err != nil {
		return nil, errors.New("не удалось проверить назначение")
	}
	if !isAssigned {
		return nil, errors.New("вы не являетесь исполнителем")
	}

	if err := s.repo.RemoveAssignee(id, memberId); err != nil {
		return nil, errors.New("не удалось отказаться от задания")
	}

	// If was IN_PROGRESS and now has free slots, revert to OPEN
	if task.Status == models.TaskStatusInProgress {
		if err := s.repo.UpdateStatus(id, models.TaskStatusOpen); err != nil {
			return nil, errors.New("не удалось обновить статус задания")
		}
	}

	return s.repo.GetById(id)
}

func (s *TaskExchangeService) RemoveAssignee(id int64, assigneeId int64, requesterId int64, isAdmin bool) (*models.TaskExchange, error) {
	task, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("задание не найдено")
	}
	if !isAdmin && task.CreatorId != requesterId {
		return nil, errors.New("только автор или админ может удалять исполнителей")
	}

	isAssigned, err := s.repo.IsAssignee(id, assigneeId)
	if err != nil {
		return nil, errors.New("не удалось проверить назначение")
	}
	if !isAssigned {
		return nil, errors.New("пользователь не является исполнителем")
	}

	if err := s.repo.RemoveAssignee(id, assigneeId); err != nil {
		return nil, errors.New("не удалось удалить исполнителя")
	}

	// If was IN_PROGRESS and now has free slots, revert to OPEN
	if task.Status == models.TaskStatusInProgress {
		if err := s.repo.UpdateStatus(id, models.TaskStatusOpen); err != nil {
			return nil, errors.New("не удалось обновить статус задания")
		}
	}

	return s.repo.GetById(id)
}

func (s *TaskExchangeService) MarkDone(id int64, memberId int64, isAdmin bool) (*models.TaskExchange, error) {
	task, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("задание не найдено")
	}
	if !isAdmin && task.CreatorId != memberId {
		return nil, errors.New("только автор или админ может отметить выполнение")
	}
	rows, err := s.repo.MarkDone(id)
	if err != nil {
		return nil, errors.New("не удалось отметить задание выполненным")
	}
	if rows == 0 {
		return nil, errors.New("задание не в работе")
	}
	return s.repo.GetById(id)
}

func (s *TaskExchangeService) Approve(id int64) (*models.TaskExchange, error) {
	rows, err := s.repo.Approve(id)
	if err != nil {
		return nil, errors.New("не удалось одобрить задание")
	}
	if rows == 0 {
		return nil, errors.New("задание не на проверке")
	}
	return s.repo.GetById(id)
}

func (s *TaskExchangeService) Reject(id int64) (*models.TaskExchange, error) {
	rows, err := s.repo.Reject(id)
	if err != nil {
		return nil, errors.New("не удалось отклонить задание")
	}
	if rows == 0 {
		return nil, errors.New("задание не на проверке")
	}
	return s.repo.GetById(id)
}

func (s *TaskExchangeService) Delete(id int64, memberId int64, isAdmin bool) error {
	task, err := s.repo.GetById(id)
	if err != nil {
		return errors.New("задание не найдено")
	}
	if isAdmin {
		return s.repo.Delete(id)
	}
	if task.CreatorId != memberId {
		return errors.New("только автор может удалить задание")
	}
	if task.Status != models.TaskStatusOpen {
		return errors.New("можно удалить только открытые задания")
	}
	return s.repo.Delete(id)
}
