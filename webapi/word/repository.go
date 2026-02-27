package word

import (
	"context"
	"webapi/models"
)

type Repository interface {
	GetAll(ctx context.Context, filter WordFilter) []*models.Word
	GetById(ctx context.Context, id string) *models.Word
	Create(ctx context.Context, word *models.Word) (*models.Word, error)
	UpdateById(ctx context.Context, id string, word *models.Word) (*models.Word, error)
	Delete(ctx context.Context, id string) error
}
