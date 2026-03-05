package service

import (
	"fmt"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

type KudosService struct {
	repo     *repository.KudosRepository
	pointSvc *PointsService
}

func NewKudosService() *KudosService {
	return &KudosService{
		repo:     repository.NewKudosRepository(),
		pointSvc: NewPointsService(),
	}
}

func (s *KudosService) Send(fromId, toId int64, message string) (*models.Kudos, error) {
	if fromId == toId {
		return nil, fmt.Errorf("нельзя отправить благодарность самому себе")
	}
	if message == "" {
		return nil, fmt.Errorf("сообщение не может быть пустым")
	}

	count, err := s.repo.CountTodayByFrom(fromId)
	if err != nil {
		return nil, err
	}
	if count >= 3 {
		return nil, fmt.Errorf("можно отправить не более 3 благодарностей в день")
	}

	kudos := &models.Kudos{
		FromId:  fromId,
		ToId:    toId,
		Message: message,
	}

	if err := s.repo.Create(kudos); err != nil {
		return nil, err
	}

	// Award points to the receiver
	go s.pointSvc.GiveForAction(toId, models.PointReasonKudosReceived, "kudos", kudos.Id, "Благодарность от участника")

	return kudos, nil
}

func (s *KudosService) GetRecent(limit, offset int) ([]models.KudosPublic, int64, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	return s.repo.GetRecent(limit, offset)
}
