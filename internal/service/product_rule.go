package service

import (
	"errors"
)

func (s *ProductService) checkDuplicateSKU(sku string) error {

	duplicate, err := s.repo.IsSKUDuplicate(sku)
	if err != nil {
		return err
	}

	if duplicate {
		return errors.New("SKU Code already exists")
	}

	return nil
}

