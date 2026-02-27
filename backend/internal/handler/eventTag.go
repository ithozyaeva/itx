package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
)

type EventTagHandler = BaseHandler[models.EventTag]

func NewEventTagHandler() EventTagHandler {
	return *NewBaseHandler(service.NewEventTagService())
}