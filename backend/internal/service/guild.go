package service

import (
	"fmt"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

type GuildService struct {
	repo *repository.GuildRepository
}

func NewGuildService() *GuildService {
	return &GuildService{
		repo: repository.NewGuildRepository(),
	}
}

func (s *GuildService) GetAll(memberId int64) ([]models.GuildPublic, error) {
	return s.repo.GetAll(memberId)
}

func (s *GuildService) Create(ownerId int64, req *models.CreateGuildRequest) (*models.Guild, error) {
	// Check if member is already in a guild
	existingId, _ := s.repo.GetMemberGuildId(ownerId)
	if existingId > 0 {
		return nil, fmt.Errorf("вы уже состоите в гильдии")
	}

	guild := &models.Guild{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
		OwnerId:     ownerId,
	}
	if guild.Icon == "" {
		guild.Icon = "users"
	}
	if guild.Color == "" {
		guild.Color = "#6366f1"
	}

	if err := s.repo.Create(guild); err != nil {
		return nil, err
	}

	// Owner automatically joins
	s.repo.Join(guild.Id, ownerId)

	return guild, nil
}

func (s *GuildService) Join(guildId, memberId int64) error {
	existingId, _ := s.repo.GetMemberGuildId(memberId)
	if existingId > 0 {
		return fmt.Errorf("вы уже состоите в гильдии")
	}
	return s.repo.Join(guildId, memberId)
}

func (s *GuildService) Leave(guildId, memberId int64) error {
	guild, err := s.repo.GetById(guildId)
	if err != nil {
		return err
	}
	if guild.OwnerId == memberId {
		return fmt.Errorf("создатель не может покинуть гильдию")
	}
	return s.repo.Leave(guildId, memberId)
}

func (s *GuildService) Update(guildId, memberId int64, req *models.CreateGuildRequest) error {
	guild, err := s.repo.GetById(guildId)
	if err != nil {
		return err
	}
	if guild.OwnerId != memberId {
		return fmt.Errorf("только создатель может редактировать гильдию")
	}
	guild.Name = req.Name
	guild.Description = req.Description
	guild.Icon = req.Icon
	guild.Color = req.Color
	return s.repo.Update(guild)
}

func (s *GuildService) Delete(guildId, memberId int64) error {
	guild, err := s.repo.GetById(guildId)
	if err != nil {
		return err
	}
	if guild.OwnerId != memberId {
		return fmt.Errorf("только создатель может удалить гильдию")
	}
	return s.repo.Delete(guildId)
}

func (s *GuildService) GetMembers(guildId int64) ([]models.MemberPointsBalance, error) {
	return s.repo.GetGuildMembers(guildId)
}
