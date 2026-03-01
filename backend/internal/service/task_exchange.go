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
	task.Status = models.TaskStatusOpen
	return s.repo.Create(task)
}

func (s *TaskExchangeService) Assign(id int64, memberId int64) (*models.TaskExchange, error) {
	task, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("задание не найдено")
	}
	if task.CreatorId == memberId {
		return nil, errors.New("нельзя взять своё задание")
	}
	rows, err := s.repo.Assign(id, memberId)
	if err != nil {
		return nil, errors.New("не удалось взять задание")
	}
	if rows == 0 {
		return nil, errors.New("задание недоступно для взятия")
	}
	return s.repo.GetById(id)
}

func (s *TaskExchangeService) Unassign(id int64, memberId int64) (*models.TaskExchange, error) {
	task, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("задание не найдено")
	}
	if task.AssigneeId == nil || *task.AssigneeId != memberId {
		return nil, errors.New("вы не являетесь исполнителем")
	}
	rows, err := s.repo.Unassign(id)
	if err != nil {
		return nil, errors.New("не удалось отказаться от задания")
	}
	if rows == 0 {
		return nil, errors.New("задание недоступно для отказа")
	}
	return s.repo.GetById(id)
}

func (s *TaskExchangeService) MarkDone(id int64, memberId int64) (*models.TaskExchange, error) {
	task, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("задание не найдено")
	}
	if task.AssigneeId == nil || *task.AssigneeId != memberId {
		return nil, errors.New("вы не являетесь исполнителем")
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
