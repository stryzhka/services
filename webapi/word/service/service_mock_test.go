package service

import (
	"context"
	"testing"
	"webapi/models"
	"webapi/word"
	"webapi/word/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSuccessCreate(t *testing.T) {
	r := new(repository.MockWordRepository)
	s := NewWordService(r)

	expectedWord := &models.Word{
		EngText:       "test",
		NativeText:    "test",
		Transcription: "test",
		Difficulty:    "test",
		CategoryId:    uuid.Nil,
	}
	r.On("Create", mock.MatchedBy(func(w *models.Word) bool {
		return w.EngText == "test" &&
			w.NativeText == "test" &&
			w.Transcription == "test" &&
			w.Difficulty == "test" &&
			w.CategoryId == uuid.Nil
	})).Return(expectedWord, nil)
	expWord, err := s.Create(context.Background(), "test", "test", "test", "test", "")
	assert.NoError(t, err)
	assert.Equal(t, expectedWord, expWord)
}

func TestFailValidationCreate(t *testing.T) {
	r := new(repository.MockWordRepository)
	s := NewWordService(r)
	_, err := s.Create(context.Background(), "", "test", "test", "test", "")
	assert.Error(t, err)
	assert.IsType(t, word.ErrValidation, err)
}

func TestFailValidationDelete(t *testing.T) {
	r := new(repository.MockWordRepository)
	s := NewWordService(r)
	err := s.Delete(context.Background(), "")
	assert.Error(t, err)
	assert.IsType(t, word.ErrValidation, err)
}

func TestFailValidationUpdate(t *testing.T) {
	r := new(repository.MockWordRepository)
	s := NewWordService(r)
	w := &models.Word{
		EngText:       "test",
		NativeText:    "test",
		Transcription: "test",
		Difficulty:    "test",
		CategoryId:    uuid.Nil,
	}
	_, err := s.UpdateById(context.Background(), "", w)
	assert.Error(t, err)
	assert.IsType(t, word.ErrValidation, err)
}
