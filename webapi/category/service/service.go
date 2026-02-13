package service

import (
	"context"
	"services/webapi/category"
	"services/webapi/models"
	"strings"

	"github.com/google/uuid"
)

type CategoryService struct {
	r category.Repository
}

func NewCategoryService(r category.Repository) *CategoryService {
	return &CategoryService{r: r}
}

func (s *CategoryService) GetAll(ctx context.Context) []*models.Category {
	return s.r.GetAll(ctx)
}

func (s *CategoryService) GetById(ctx context.Context, id string) *models.Category {
	return s.r.GetById(ctx, id)
}

func (s *CategoryService) GetAllWords(ctx context.Context, categoryId string) []*models.Word {
	return s.r.GetAllWords(ctx, categoryId)
}

func (s *CategoryService) Create(ctx context.Context, name string) (*models.Category, error) {
	newCategory := &models.Category{
		Id:   uuid.New().String(),
		Name: name,
	}
	if strings.TrimSpace(newCategory.Name) == "" {
		return nil, category.ErrValidation
	}
	return s.r.Create(ctx, newCategory)
}

func (s *CategoryService) UpdateById(ctx context.Context, id string, newCategory *models.Category) (*models.Category, error) {
	if strings.TrimSpace(newCategory.Name) == "" || strings.TrimSpace(id) == "" {
		return nil, category.ErrValidation
	}
	return s.r.UpdateById(ctx, id, newCategory)
}

func (s *CategoryService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return category.ErrValidation
	}
	return s.r.Delete(ctx, id)
}
