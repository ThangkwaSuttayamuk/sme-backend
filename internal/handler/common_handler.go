package handler

import (
	"wearlab_backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type CommonHandler struct {
	service *service.CommonService
}

func NewCommonHandler(service *service.CommonService) *CommonHandler {
	return &CommonHandler{service: service}
}

func (h *CommonHandler) GetTypes(c *fiber.Ctx) error {
	types, err := h.service.GetTypes()

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(types)
}