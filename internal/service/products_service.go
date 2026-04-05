package service

import (
	"restapi/internal/model"
	"restapi/internal/repository"
)

type ProductService struct {
	repo *repository.Repository
}

func NewProductService(repo *repository.Repository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

// GetAll retrieves all products
func (s *ProductService) GetAll() ([]model.Product, error) {
	return s.repo.GetAllProducts()
}
