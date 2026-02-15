package service

import (
	"errors"
	"strings"
	"wearlab_backend/internal/domain"
)

func validateProduct(product *domain.Product) error {

	if strings.TrimSpace(product.Name) == "" {
		return errors.New("Name is required")
	}

	if len(product.SKU) < 3 {
		return errors.New("SKU Code must be at least 3 characters")
	}
	
	if product.Price <= 0 {
		return errors.New("Price must be greater than 0")
	}

	if product.Stock < 0 {
		return errors.New("Stock must not be negative")
	}

	return nil
}


