package word

import (
	"context"
	"webapi/models"
)

type Service interface {
	GetAll(ctx context.Context) []*models.Word
	GetById(ctx context.Context, id string) *models.Word
	Create(ctx context.Context, engText, nativeText, transcription, difficulty, categoryId string) (*models.Word, error)
	UpdateById(ctx context.Context, id string, newWord *models.Word) (*models.Word, error)
	Delete(ctx context.Context, id string) error
}
