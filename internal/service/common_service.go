package service

import (
	"wearlab_backend/internal/domain"
	"wearlab_backend/internal/repository"
)

type CommonService struct {
	repo *repository.CommonRepository
}

func NewCommonService(repo *repository.CommonRepository) *CommonService {
	return &CommonService{repo: repo}
}

func (s *CommonService) GetTypes() ([]domain.Type, error) {
	return s.repo.GetTypes()
}