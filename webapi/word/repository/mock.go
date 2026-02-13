package repository

import (
	"context"
	"webapi/models"

	"github.com/stretchr/testify/mock"
)

type MockWordRepository struct {
	mock.Mock
}

func (r *MockWordRepository) GetAll(ctx context.Context) []*models.Word {
	args := r.Called()
	return args.Get(0).([]*models.Word)
}

func (r *MockWordRepository) GetById(ctx context.Context, id string) *models.Word {
	args := r.Called(id)
	return args.Get(0).(*models.Word)
}

func (r *MockWordRepository) Create(ctx context.Context, word *models.Word) (*models.Word, error) {
	args := r.Called(word)
	return args.Get(0).(*models.Word), args.Error(1)
}

func (r *MockWordRepository) UpdateById(ctx context.Context, id string, word *models.Word) (*models.Word, error) {
	args := r.Called(id, word)
	return args.Get(0).(*models.Word), args.Error(1)
}

func (r *MockWordRepository) Delete(ctx context.Context, id string) error {
	args := r.Called(id)
	return args.Error(0)
}
