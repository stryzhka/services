package category

import (
	"context"
	"services/webapi/models"
)

type Repository interface {
	GetAll(ctx context.Context) []*models.Category
	GetById(ctx context.Context, id string) *models.Category
	GetAllWords(ctx context.Context, categoryId string) []*models.Word
	Create(ctx context.Context, category *models.Category) (*models.Category, error)
	UpdateById(ctx context.Context, id string, category *models.Category) (*models.Category, error)
	Delete(ctx context.Context, id string) error
}
