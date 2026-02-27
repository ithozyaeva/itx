package handler

import (
	"ithozyeva/internal/service"

	"github.com/gofiber/fiber/v2"
)

type DictionaryHandler struct {
	dictionaryService *service.DictionaryService
}

func NewDictionaryHandler() *DictionaryHandler {
	return &DictionaryHandler{
		dictionaryService: service.NewDictionaryService(),
	}
}

func (h *DictionaryHandler) GetDictionaries(c *fiber.Ctx) error {
	dictionaries := h.dictionaryService.GetAllDictionaries()
	return c.JSON(dictionaries)
}
