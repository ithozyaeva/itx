package handler

import (
	"errors"
	"fmt"
	"ithozyeva/internal/bot"
	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/service"
	"ithozyeva/internal/utils"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type EventsHandler struct {
	BaseHandler[models.Event]
	svc      *service.EventsService
	auditSvc *service.AuditService
}

func NewEventsHandler() *EventsHandler {
	svc := service.NewEventsService()
	return &EventsHandler{
		BaseHandler: *NewBaseHandler(svc),
		svc:         svc,
		auditSvc:    service.NewAuditService(),
	}
}

var EventsSearchFields = map[string]string{
	"dateFrom": "date >= ?",
	"dateTo":   "date < ?",
}

type EventsSearchRequest struct {
	Limit     *int    `query:"limit"`
	Offset    *int    `query:"offset"`
	DateFrom  *string `query:"dateFrom"`
	DateTo    *string `query:"dateTo"`
	Title     *string `query:"title"`
	PlaceType *string `query:"placeType"`
}

func (h *EventsHandler) Search(c *fiber.Ctx) error {
	req := new(EventsSearchRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	filter := make(repository.SearchFilter)
	// *filter = make(map[string]interface{})

	if req.DateFrom != nil {
		// Предстоящие: обычные события с датой >= now ИЛИ повторяющиеся с активными будущими вхождениями
		filter["(date >= ? OR (is_repeating = TRUE AND (repeat_end_date IS NULL OR repeat_end_date >= ?)))"] =
			[]interface{}{*req.DateFrom, *req.DateFrom}
	}
	if req.DateTo != nil {
		// Архив: обычные прошедшие события И повторяющиеся, у которых больше нет будущих вхождений
		filter["(date < ? AND NOT (is_repeating = TRUE AND (repeat_end_date IS NULL OR repeat_end_date >= ?)))"] =
			[]interface{}{*req.DateTo, *req.DateTo}
	}
	if req.Title != nil && *req.Title != "" {
		filter["title ILIKE ?"] = "%" + *req.Title + "%"
	}
	if req.PlaceType != nil && *req.PlaceType != "" {
		filter["place_type = ?"] = *req.PlaceType
	}

	result, err := h.service.Search(req.Limit, req.Offset, &filter, &repository.Order{
		ColumnBy: "date",
		Order:    "DESC",
	})
	if err != nil {
		log.Printf("events search error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка поиска событий"})
	}

	return c.JSON(result)
}

func (h *EventsHandler) GetOld(c *fiber.Ctx) error {

	result, err := h.service.Search(nil, nil, &repository.SearchFilter{
		"date < ?": gorm.Expr("CURRENT_TIMESTAMP"),
	}, &repository.Order{
		ColumnBy: "date",
		Order:    "DESC",
	})
	if err != nil {
		log.Printf("get old events error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки событий"})
	}

	return c.JSON(result)
}

func (h *EventsHandler) GetNext(c *fiber.Ctx) error {
	result, err := h.service.Search(nil, nil, &repository.SearchFilter{
		"date >= ?": gorm.Expr("CURRENT_TIMESTAMP"),
	}, &repository.Order{
		ColumnBy: "date",
		Order:    "ASC",
	})
	if err != nil {
		log.Printf("get next events error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки событий"})
	}

	return c.JSON(result)
}

func (h *EventsHandler) AddMember(c *fiber.Ctx) error {
	req := new(WorkWithEventRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	member := c.Locals("member").(*models.Member)

	result, err := h.svc.AddMember(req.EventId, int(member.Id))
	if err != nil {
		if errors.Is(err, service.ErrParticipantLimitReached) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Достигнут лимит участников"})
		}
		log.Printf("add member to event error (event=%d, member=%d): %v", req.EventId, member.Id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка регистрации на событие"})
	}

	return c.JSON(result)
}

type WorkWithEventRequest struct {
	EventId int `json:"eventId" query:"eventId"`
}

func (h *EventsHandler) RemoveMember(c *fiber.Ctx) error {
	req := new(WorkWithEventRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	member := c.Locals("member").(*models.Member)

	result, err := h.svc.RemoveMember(req.EventId, int(member.Id))
	if err != nil {
		log.Printf("remove member from event error (event=%d, member=%d): %v", req.EventId, member.Id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка отмены регистрации"})
	}

	return c.JSON(result)
}

func (h *EventsHandler) GetICSFile(c *fiber.Ctx) error {
	req := new(WorkWithEventRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	event, err := h.svc.GetById(int64(req.EventId))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Событие не найдено"})
	}

	ics := utils.GenerateICS(event)

	c.Set("Content-Type", "text/calendar")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=event_%d.ics", event.Id))
	return c.SendString(ics)
}

// Create переопределяет базовый метод Create для отправки алертов при создании события
func (h *EventsHandler) Create(c *fiber.Ctx) error {
	event := new(models.Event)
	if err := c.BodyParser(event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	// Подставляем название эксклюзивного чата
	h.resolveExclusiveChatTitle(event)

	result, err := h.service.Create(event)
	if err != nil {
		log.Printf("create event error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка создания события"})
	}

	// Отправляем инициализирующие алерты в фоне
	go func() {
		telegramBot := bot.GetGlobalBot()
		if telegramBot == nil {
			log.Printf("Telegram bot is not initialized, skipping alerts for event %d", result.Id)
			return
		}
		if err := telegramBot.SendInitialEventAlerts(result); err != nil {
			log.Printf("Error sending initial event alerts: %v", err)
		}
	}()

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionCreate, "event", result.Id, result.Title)

	return c.Status(fiber.StatusCreated).JSON(result)
}

// Update переопределяет базовый метод Update для отправки уведомлений об изменении события
func (h *EventsHandler) Update(c *fiber.Ctx) error {
	event := new(models.Event)
	if err := c.BodyParser(event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный запрос"})
	}

	// Подставляем название эксклюзивного чата
	h.resolveExclusiveChatTitle(event)

	result, err := h.service.Update(event)
	if err != nil {
		log.Printf("update event error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка обновления события"})
	}

	// Отправляем уведомления об изменении события в фоне
	go func() {
		telegramBot := bot.GetGlobalBot()
		if telegramBot == nil {
			log.Printf("Telegram bot is not initialized, skipping update alerts for event %d", result.Id)
			return
		}
		if err := telegramBot.SendEventUpdateAlert(result); err != nil {
			log.Printf("Error sending event update alerts: %v", err)
		} else {
			log.Printf("Successfully sent update alerts for event %d to all subscribed members", result.Id)
		}
	}()

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionUpdate, "event", result.Id, result.Title)

	return c.JSON(result)
}

// Delete переопределяет базовый метод Delete для добавления аудит-лога и уведомлений об отмене
func (h *EventsHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ID"})
	}

	entity, err := h.service.GetById(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Событие не найдено"})
	}

	// Snapshot member IDs before cascade delete removes event_members
	memberIds := GetEventMemberIds(int64(id))

	if err := h.service.Delete(entity); err != nil {
		log.Printf("delete event error (id=%d): %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка удаления события"})
	}

	// Send notifications after delete using pre-snapshotted member IDs
	go func() {
		telegramBot := bot.GetGlobalBot()
		if telegramBot == nil {
			log.Printf("Telegram bot is not initialized, skipping cancel alerts for event %d", entity.Id)
			return
		}
		if err := telegramBot.SendEventCancelAlert(entity); err != nil {
			log.Printf("Error sending event cancel alerts: %v", err)
		}
	}()

	notifBody := fmt.Sprintf("Событие \"%s\" было отменено.", entity.Title)
	go CreateNotificationsForMembers(memberIds, "event_cancel", "Событие отменено", notifBody)

	go h.auditSvc.Log(getActorId(c), getActorName(c), getActorType(c), models.AuditActionDelete, "event", int64(id), entity.Title)

	return c.SendStatus(fiber.StatusNoContent)
}

// resolveExclusiveChatTitle подставляет название чата по его ID
func (h *EventsHandler) resolveExclusiveChatTitle(event *models.Event) {
	if event.ExclusiveChatID == nil || *event.ExclusiveChatID == 0 {
		event.ExclusiveChatTitle = ""
		return
	}
	chatActivitySvc := service.NewChatActivityService()
	chats, err := chatActivitySvc.GetTrackedChats()
	if err != nil {
		return
	}
	for _, chat := range chats {
		if chat.ChatID == *event.ExclusiveChatID {
			event.ExclusiveChatTitle = chat.Title
			return
		}
	}
}
