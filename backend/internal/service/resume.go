package service

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/utils"
)

type ResumeService struct {
	repo *repository.ResumeRepository
}

func NewResumeService() *ResumeService {
	return &ResumeService{
		repo: repository.NewResumeRepository(),
	}
}

func (s *ResumeService) UploadResume(member *models.Member, fileName string, contentType string, content []byte, req *models.CreateResumeRequest) (*models.Resume, *utils.ParsedResumeData, error) {
	client, err := utils.NewS3Client()
	if err != nil {
		return nil, nil, err
	}

	ext := strings.ToLower(filepath.Ext(fileName))
	if ext != ".pdf" && ext != ".docx" && ext != ".doc" {
		return nil, nil, fmt.Errorf("unsupported file format %s", ext)
	}

	key := fmt.Sprintf("resumes/%d/%s%s", member.TelegramID, uuid.NewString(), ext)
	if err := client.Upload(context.Background(), key, content, contentType); err != nil {
		return nil, nil, err
	}

	var parsed *utils.ParsedResumeData
	if data, err := utils.ParseResume(fileName, content); err == nil {
		parsed = data
	} else {
		log.Printf("resume parse skipped: %v", err)
		parsed = &utils.ParsedResumeData{}
	}

	workExperience := strings.TrimSpace(req.WorkExperience)
	if workExperience == "" {
		workExperience = parsed.WorkExperience
	}

	desiredPosition := strings.TrimSpace(req.DesiredPosition)
	if desiredPosition == "" {
		desiredPosition = parsed.DesiredPosition
	}

	workFormat := req.WorkFormat
	if (workFormat == "" || !workFormat.IsValid()) && parsed.WorkFormat.IsValid() {
		workFormat = parsed.WorkFormat
	}

	resume := &models.Resume{
		TgID:            member.TelegramID,
		FilePath:        key,
		FileName:        fileName,
		WorkExperience:  workExperience,
		DesiredPosition: desiredPosition,
		WorkFormat:      workFormat,
	}

	created, err := s.repo.Create(resume)
	if err != nil {
		_ = client.Delete(context.Background(), key)
		return nil, nil, err
	}

	created.ParsedConfidence = parsed.Confidence
	return created, parsed, nil
}

func (s *ResumeService) ListByTelegramID(tgID int64) ([]models.Resume, error) {
	return s.repo.ListByTelegramID(tgID)
}

func (s *ResumeService) UpdateResume(id, tgID int64, payload *models.UpdateResumeRequest) (*models.Resume, error) {
	resume, err := s.repo.GetByIDAndTelegram(id, tgID)
	if err != nil {
		return nil, err
	}

	if payload.WorkExperience != nil {
		resume.WorkExperience = strings.TrimSpace(*payload.WorkExperience)
	}
	if payload.DesiredPosition != nil {
		resume.DesiredPosition = strings.TrimSpace(*payload.DesiredPosition)
	}
	if payload.WorkFormat != nil && payload.WorkFormat.IsValid() {
		resume.WorkFormat = *payload.WorkFormat
	}

	return s.repo.Update(resume)
}

func (s *ResumeService) DeleteResume(id, tgID int64) error {
	resume, err := s.repo.GetByIDAndTelegram(id, tgID)
	if err != nil {
		return err
	}

	client, err := utils.NewS3Client()
	if err != nil {
		return err
	}

	if err := client.Delete(context.Background(), resume.FilePath); err != nil {
		return err
	}

	return s.repo.Delete(resume)
}

func (s *ResumeService) SearchForAdmin(limit *int, offset *int, filter *models.ResumeFilter) (*models.RegistrySearch[models.Resume], error) {
	items, total, err := s.repo.SearchForAdmin(limit, offset, filter)
	if err != nil {
		return nil, err
	}
	return &models.RegistrySearch[models.Resume]{
		Items: items,
		Total: int(total),
	}, nil
}

func (s *ResumeService) GenerateArchive(filter *models.ResumeFilter) ([]byte, error) {
	items, _, err := s.repo.SearchForAdmin(nil, nil, filter)
	if err != nil {
		return nil, err
	}

	client, err := utils.NewS3Client()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	for _, resume := range items {
		data, err := client.Download(context.Background(), resume.FilePath)
		if err != nil {
			return nil, err
		}
		filename := fmt.Sprintf("%d_%s", resume.TgID, sanitizeFileName(resume.FileName))
		w, err := zipWriter.Create(filename)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(data); err != nil {
			return nil, err
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func sanitizeFileName(name string) string {
	replacer := strings.NewReplacer(" ", "_")
	return replacer.Replace(name)
}

func (s *ResumeService) GetByIdWithMember(id int64) (*models.Resume, error) {
	return s.repo.GetByIdWithMember(id)
}
