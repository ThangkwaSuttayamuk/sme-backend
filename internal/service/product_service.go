package service

import (
	"errors"
	"math"
	"wearlab_backend/internal/domain"
	"wearlab_backend/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(product *domain.Product) error {

	if err := validateProduct(product); err != nil {
		return err
	}

	if err := s.checkDuplicateSKU(product.SKU); err != nil {
		return err
	}

	// Round price to 2 decimal places
	product.Price = math.Round(product.Price*100) / 100

	return s.repo.Create(product)
}

func (s *ProductService) GetProductsWithFilter(category, keyword string) ([]domain.Product, int, float64, error) {
	return s.repo.GetWithFilter(category, keyword)
}

func (s *ProductService) SellProduct(productID int, quantity int) error {

	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	return s.repo.UpdateStock(productID, quantity)
}

func (s *ProductService) BulkPriceUpdate(items []domain.BulkPriceUpdateItem) (int, int, error) {

	total := len(items)
	updated := 0

	for _, item := range items {

		if item.NewPrice <= 0 {
			continue // skip invalid price
		}

		// Round price to 2 decimal places
		roundedPrice := math.Round(item.NewPrice*100) / 100

		rows, err := s.repo.UpdatePrice(item.ProductID, roundedPrice)
		if err != nil {
			return total, updated, err
		}

		if rows > 0 {
			updated++
		}
	}

	return total, updated, nil
}
