package service

import (
	"context"
	"webapi/models"
	"webapi/word"

	"github.com/google/uuid"
)

type WordService struct {
	r              word.Repository
	eventPublisher word.EventPublisher
}

func NewWordService(r word.Repository, e word.EventPublisher) *WordService {
	return &WordService{
		r:              r,
		eventPublisher: e,
	}
}

func validateWord(w models.Word) error {
	if len(w.EngText) == 0 || len(w.NativeText) == 0 || len(w.Difficulty) == 0 || len(w.ConfirmedUserId) == 0 {
		return word.ErrValidation
	}
	return nil
}

func (w *WordService) GetAll(ctx context.Context, filter word.WordFilter) []*models.Word {
	return w.r.GetAll(ctx, filter)
}

func (w *WordService) GetById(ctx context.Context, id string) *models.Word {
	return w.r.GetById(ctx, id)
}

func (w *WordService) Create(ctx context.Context, engText, nativeText, transcription, difficulty, categoryId, confirmedUserId string) (*models.Word, error) {

	categoryUuid, err := uuid.Parse(categoryId)
	if err != nil {
		categoryUuid = uuid.Nil
	}
	_word := &models.Word{
		Id:              uuid.New().String(),
		EngText:         engText,
		NativeText:      nativeText,
		Transcription:   transcription,
		Difficulty:      difficulty,
		CategoryId:      categoryUuid.String(),
		ConfirmedUserId: confirmedUserId,
		ConfirmedStatus: "pending",
	}
	err = validateWord(*_word)
	if err != nil {
		return nil, err
	}
	e := &word.WordConfirmMessage{UserId: _word.ConfirmedUserId, ObjectId: _word.Id}
	err = w.eventPublisher.Publish(ctx, "obj-to-users", _word.Id, e)
	if err != nil {
		return nil, err
	}
	return w.r.Create(ctx, _word)
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
