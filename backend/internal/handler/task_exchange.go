package handler

import (
	"fmt"
	"log"
	"strconv"

	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type TaskExchangeHandler struct {
	svc      *service.TaskExchangeService
	pointSvc *service.PointsService
}

func NewTaskExchangeHandler() *TaskExchangeHandler {
	return &TaskExchangeHandler{
		svc:      service.NewTaskExchangeService(),
		pointSvc: service.NewPointsService(),
	}
}

func (h *TaskExchangeHandler) Search(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "20")
	offsetStr := c.Query("offset", "0")
	status := c.Query("status")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	items, total, err := h.svc.Search(statusPtr, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить задания"})
	}

	return c.JSON(fiber.Map{"items": items, "total": total})
}

func (h *TaskExchangeHandler) GetById(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	task, err := h.svc.GetById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Задание не найдено"})
	}

	return c.JSON(task)
}

func (h *TaskExchangeHandler) Create(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)

	var req models.CreateTaskExchangeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Название обязательно"})
	}

	task := &models.TaskExchange{
		Title:       req.Title,
		Description: req.Description,
		CreatorId:   member.Id,
	}

	created, err := h.svc.Create(task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось создать задание"})
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *TaskExchangeHandler) Assign(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	task, err := h.svc.Assign(id, member.Id)
	if err != nil {
		log.Printf("Task assign error (task=%d, member=%d): %v", id, member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Не удалось взять задание"})
	}

	return c.JSON(task)
}

func (h *TaskExchangeHandler) Unassign(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	task, err := h.svc.Unassign(id, member.Id)
	if err != nil {
		log.Printf("Task unassign error (task=%d, member=%d): %v", id, member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Не удалось отказаться от задания"})
	}

	return c.JSON(task)
}

func (h *TaskExchangeHandler) MarkDone(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	task, err := h.svc.MarkDone(id, member.Id)
	if err != nil {
		log.Printf("Task markDone error (task=%d, member=%d): %v", id, member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Не удалось отметить задание выполненным"})
	}

	// Notify creator that task is done
	go func() {
		if err := CreateNotification(task.CreatorId, "task", "Задание выполнено",
			fmt.Sprintf("Задание «%s» помечено как выполненное и ожидает проверки", task.Title)); err != nil {
			log.Printf("Error creating notification: %v", err)
		}
	}()

	return c.JSON(task)
}

func (h *TaskExchangeHandler) Approve(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	task, err := h.svc.Approve(id)
	if err != nil {
		log.Printf("Task approve error (task=%d): %v", id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Не удалось одобрить задание"})
	}

	// Award points
	go func() {
		h.pointSvc.GiveForAction(task.CreatorId, models.PointReasonTaskCreate, "task_exchange", task.Id,
			fmt.Sprintf("Создание задания: %s", task.Title))

		if task.AssigneeId != nil {
			h.pointSvc.GiveForAction(*task.AssigneeId, models.PointReasonTaskExecute, "task_exchange", task.Id,
				fmt.Sprintf("Выполнение задания: %s", task.Title))
		}

		// Notify creator
		if err := CreateNotification(task.CreatorId, "task", "Задание одобрено",
			fmt.Sprintf("Задание «%s» одобрено! Вам начислено %d баллов", task.Title, models.PointValues[models.PointReasonTaskCreate])); err != nil {
			log.Printf("Error creating notification: %v", err)
		}

		// Notify assignee
		if task.AssigneeId != nil {
			if err := CreateNotification(*task.AssigneeId, "task", "Задание одобрено",
				fmt.Sprintf("Выполнение задания «%s» одобрено! Вам начислено %d баллов", task.Title, models.PointValues[models.PointReasonTaskExecute])); err != nil {
				log.Printf("Error creating notification: %v", err)
			}
		}
	}()

	return c.JSON(task)
}

func (h *TaskExchangeHandler) Reject(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	// Get task before rejection to know the assignee
	taskBefore, err := h.svc.GetById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Задание не найдено"})
	}
	assigneeId := taskBefore.AssigneeId

	task, err := h.svc.Reject(id)
	if err != nil {
		log.Printf("Task reject error (task=%d): %v", id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Не удалось отклонить задание"})
	}

	// Notify assignee about rejection
	go func() {
		if assigneeId != nil {
			if err := CreateNotification(*assigneeId, "task", "Задание отклонено",
				fmt.Sprintf("Выполнение задания «%s» отклонено. Задание возвращено в открытые", task.Title)); err != nil {
				log.Printf("Error creating notification: %v", err)
			}
		}
	}()

	return c.JSON(task)
}

func (h *TaskExchangeHandler) Delete(c *fiber.Ctx) error {
	member := c.Locals("member").(*models.Member)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	// Check if admin
	isAdmin := false
	for _, role := range member.Roles {
		if role == models.MemberRoleAdmin {
			isAdmin = true
			break
		}
	}

	if err := h.svc.Delete(id, member.Id, isAdmin); err != nil {
		log.Printf("Task delete error (task=%d, member=%d): %v", id, member.Id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Не удалось удалить задание"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
