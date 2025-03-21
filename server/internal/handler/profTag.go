package handler

import (
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
)

type ProfTagHandler = BaseHandler[models.ProfTag]

func NewProfTagsHandler() ProfTagHandler {
	return *NewBaseHandler(service.NewProfTagService())
}
