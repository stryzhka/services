package category

import (
	"context"
	"services/webapi/models"
)

type Service interface {
	GetAll(ctx context.Context) []*models.Category
	GetById(ctx context.Context, id string) *models.Category
	GetAllWords(ctx context.Context, categoryId string) []*models.Word
	Create(ctx context.Context, name string) (*models.Category, error)
	UpdateById(ctx context.Context, id string, newCategory *models.Category) (*models.Category, error)
	Delete(ctx context.Context, id string) error
}
