package service

import (
	"context"
	"services/webapi/models"
	"services/webapi/word"

	"github.com/google/uuid"
)

type WordService struct {
	r word.Repository
}

func NewWordService(r word.Repository) WordService {
	return WordService{r: r}
}

func validateWord(w models.Word) error {
	if len(w.EngText) == 0 || len(w.NativeText) == 0 || len(w.Difficulty) == 0 {
		return word.ErrValidation
	}
	return nil
}

func (w *WordService) GetAll(ctx context.Context) []*models.Word {
	return w.r.GetAll(ctx)
}

func (w *WordService) GetById(ctx context.Context, id string) *models.Word {
	return w.r.GetById(ctx, id)
}

func (w *WordService) Create(ctx context.Context, engText, nativeText, transcription, difficulty, categoryId string) (*models.Word, error) {

	categoryUuid, err := uuid.Parse(categoryId)
	if err != nil {
		categoryUuid = uuid.Nil
	}
	word := &models.Word{
		Id:            uuid.New(),
		EngText:       engText,
		NativeText:    nativeText,
		Transcription: transcription,
		Difficulty:    difficulty,
		CategoryId:    categoryUuid,
	}
	err = validateWord(*word)
	if err != nil {
		return nil, err
	}
	return w.r.Create(ctx, word)
}

func (w *WordService) UpdateById(ctx context.Context, id string, newWord *models.Word) (*models.Word, error) {
	err := validateWord(*newWord)
	if err != nil {
		return nil, err
	}
	if id == "" {
		return nil, word.ErrValidation
	}
	return w.r.UpdateById(ctx, id, newWord)
}

func (w *WordService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return word.ErrValidation
	}
	return w.r.Delete(ctx, id)
}
