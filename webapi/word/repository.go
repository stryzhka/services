package word

import (
	"context"
	"services/webapi/models"

	"github.com/google/uuid"
)

type Repository interface {
	GetAll(ctx context.Context) []models.Word
	GetById(ctx context.Context, id uuid.UUID) (*models.Word, error)
	Create(ctx context.Context, word *models.Word) (*models.Word, error)
	UpdateById(ctx context.Context, id uuid.UUID, word *models.Word) (*models.Word, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
