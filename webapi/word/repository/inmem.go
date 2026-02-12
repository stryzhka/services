package repository

import (
	"context"
	"log"
	"services/webapi/models"
	word2 "services/webapi/word"

	"github.com/hashicorp/go-memdb"
)

type InmemRepository struct {
	db *memdb.MemDB
}

func NewInmemRepository(db *memdb.MemDB) *InmemRepository {
	return &InmemRepository{db: db}
}

func (m *InmemRepository) GetAll(ctx context.Context) []*models.Word {
	var words []*models.Word
	txn := m.db.Txn(false)
	defer txn.Abort()
	it, err := txn.Get("word", "id")
	if err != nil {
		return nil
	}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		w := obj.(*models.Word)
		words = append(words, w)
	}
	txn.Commit()
	return words
}

func (m *InmemRepository) GetById(ctx context.Context, id string) *models.Word {
	txn := m.db.Txn(false)
	defer txn.Abort()
	found, err := txn.First("word", "id", id)
	log.Println(found)
	if err != nil {
		log.Println(err)
		return nil
	}
	return found.(*models.Word)
}

func (m *InmemRepository) Create(ctx context.Context, word *models.Word) (*models.Word, error) {
	txn := m.db.Txn(true)
	defer txn.Abort()

	// Проверяем, существует ли уже запись с таким id
	existing, err := txn.First("word", "id", word.Id)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, word2.ErrWordAlreadyExists
	}

	err = txn.Insert("word", word)
	if err != nil {
		return nil, err
	}

	txn.Commit()
	return word, nil
}

func (m *InmemRepository) UpdateById(ctx context.Context, id string, word *models.Word) (*models.Word, error) {
	txn := m.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("word", "id", id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, word2.ErrWordNotFound
	}
	err = txn.Insert("word", word)
	if err != nil {
		return nil, err
	}
	txn.Commit()
	return word, nil
}

func (m *InmemRepository) Delete(ctx context.Context, id string) error {
	txn := m.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("word", "id", id)
	if err != nil {
		return err
	}
	if existing == nil {
		return word2.ErrWordNotFound
	}
	err = txn.Delete("word", existing)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}
