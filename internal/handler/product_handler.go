package handler

import (
	"wearlab_backend/internal/domain"
	"wearlab_backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	product := new(domain.Product)

	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := h.service.CreateProduct(product)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.SendString("Create Product Successfully.")
}

func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	category := c.Query("category")
	keyword := c.Query("keyword")

	products, total, totalValue, err := h.service.GetProductsWithFilter(category, keyword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get products")
	}

	return c.JSON(fiber.Map{
		"products":   products,
		"total":      total,
		"totalValue": totalValue,
	})
}

func (h *ProductHandler) GetProductsWithFilter(c *fiber.Ctx) error {
	category := c.Query("category")
	keyword := c.Query("keyword")

	products, total, totalValue, err := h.service.GetProductsWithFilter(category, keyword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get products")
	}

	return c.JSON(fiber.Map{
		"products":   products,
		"total":      total,
		"totalValue": totalValue,
	})
}

func (h *ProductHandler) SellProduct(c *fiber.Ctx) error {

	var req domain.SellProductRequest

	// parse body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// call service
	err := h.service.SellProduct(req.ProductID, req.Quantity)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "product sold successfully",
	})
}

func (h *ProductHandler) BulkPriceUpdate(c *fiber.Ctx) error {

	var items []domain.BulkPriceUpdateItem

	if err := c.BodyParser(&items); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	total, updated, err := h.service.BulkPriceUpdate(items)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"total":   total,
		"updated": updated,
	})
}
