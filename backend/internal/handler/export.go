package handler

import (
	"fmt"
	"strings"
	"time"

	"ithozyeva/database"
	"ithozyeva/internal/models"
	"ithozyeva/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type ExportHandler struct{}

func NewExportHandler() *ExportHandler {
	return &ExportHandler{}
}

func (h *ExportHandler) ExportMembers(c *fiber.Ctx) error {
	var members []models.Member
	if err := database.DB.Preload("MemberRoles").Find(&members).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	headers := []string{"ID", "Имя", "Фамилия", "Telegram", "Роли", "День рождения"}
	rows := make([][]string, 0, len(members))

	for _, m := range members {
		roles := make([]string, len(m.MemberRoles))
		for i, r := range m.MemberRoles {
			roles[i] = string(r.Role)
		}

		birthday := ""
		if m.Birthday != nil {
			birthday = m.Birthday.String()
		}

		rows = append(rows, []string{
			fmt.Sprintf("%d", m.Id),
			m.FirstName,
			m.LastName,
			m.Username,
			strings.Join(roles, ", "),
			birthday,
		})
	}

	data, err := utils.GenerateCSV(headers, rows)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=members_%s.csv", time.Now().Format("2006-01-02")))
	return c.Send(data)
}

func (h *ExportHandler) ExportMentors(c *fiber.Ctx) error {
	var mentors []models.MentorDbModel
	if err := database.DB.Preload("Member").Preload("ProfTags").Find(&mentors).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	headers := []string{"ID", "Имя", "Фамилия", "Telegram", "Должность", "Опыт", "Теги"}
	rows := make([][]string, 0, len(mentors))

	for _, m := range mentors {
		tags := make([]string, len(m.ProfTags))
		for i, t := range m.ProfTags {
			tags[i] = t.Title
		}

		rows = append(rows, []string{
			fmt.Sprintf("%d", m.Id),
			m.Member.FirstName,
			m.Member.LastName,
			m.Member.Username,
			m.Occupation,
			m.Experience,
			strings.Join(tags, ", "),
		})
	}

	data, err := utils.GenerateCSV(headers, rows)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=mentors_%s.csv", time.Now().Format("2006-01-02")))
	return c.Send(data)
}

func (h *ExportHandler) ExportEvents(c *fiber.Ctx) error {
	var events []models.Event
	if err := database.DB.Preload("EventTags").Preload("Members").Find(&events).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	headers := []string{"ID", "Название", "Дата", "Тип места", "Место", "Тип события", "Открытое", "Теги", "Участников"}
	rows := make([][]string, 0, len(events))

	for _, e := range events {
		tags := make([]string, len(e.EventTags))
		for i, t := range e.EventTags {
			tags[i] = t.Name
		}

		open := "Нет"
		if e.Open {
			open = "Да"
		}

		rows = append(rows, []string{
			fmt.Sprintf("%d", e.Id),
			e.Title,
			e.Date.Format("2006-01-02 15:04"),
			string(e.PlaceType),
			e.Place,
			e.EventType,
			open,
			strings.Join(tags, ", "),
			fmt.Sprintf("%d", len(e.Members)),
		})
	}

	data, err := utils.GenerateCSV(headers, rows)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=events_%s.csv", time.Now().Format("2006-01-02")))
	return c.Send(data)
}
