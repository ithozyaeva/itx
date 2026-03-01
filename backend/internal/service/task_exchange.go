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
		return nil, err
	}
	if task.Status != models.TaskStatusOpen {
		return nil, errors.New("task is not open")
	}
	if task.CreatorId == memberId {
		return nil, errors.New("cannot assign own task")
	}
	if err := s.repo.Assign(id, memberId); err != nil {
		return nil, err
	}
	return s.repo.GetById(id)
}

func (s *TaskExchangeService) Unassign(id int64, memberId int64) (*models.TaskExchange, error) {
	task, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	if task.Status != models.TaskStatusInProgress {
		return nil, errors.New("task is not in progress")
	}
	if task.AssigneeId == nil || *task.AssigneeId != memberId {
		return nil, errors.New("you are not the assignee")
	}
	if err := s.repo.Unassign(id); err != nil {
		return nil, err
	}
	return s.repo.GetById(id)
}

func (s *TaskExchangeService) MarkDone(id int64, memberId int64) (*models.TaskExchange, error) {
	task, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	if task.Status != models.TaskStatusInProgress {
		return nil, errors.New("task is not in progress")
	}
	if task.AssigneeId == nil || *task.AssigneeId != memberId {
		return nil, errors.New("you are not the assignee")
	}
	if err := s.repo.MarkDone(id); err != nil {
		return nil, err
	}
	return s.repo.GetById(id)
}

func (s *TaskExchangeService) Approve(id int64) (*models.TaskExchange, error) {
	task, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	if task.Status != models.TaskStatusDone {
		return nil, errors.New("task is not done")
	}
	if err := s.repo.Approve(id); err != nil {
		return nil, err
	}
	return s.repo.GetById(id)
}

func (s *TaskExchangeService) Reject(id int64) (*models.TaskExchange, error) {
	task, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	if task.Status != models.TaskStatusDone {
		return nil, errors.New("task is not done")
	}
	if err := s.repo.Reject(id); err != nil {
		return nil, err
	}
	return s.repo.GetById(id)
}

func (s *TaskExchangeService) Delete(id int64, memberId int64, isAdmin bool) error {
	task, err := s.repo.GetById(id)
	if err != nil {
		return err
	}
	if isAdmin {
		return s.repo.Delete(id)
	}
	if task.CreatorId != memberId {
		return errors.New("only creator can delete")
	}
	if task.Status != models.TaskStatusOpen {
		return errors.New("can only delete open tasks")
	}
	return s.repo.Delete(id)
}
