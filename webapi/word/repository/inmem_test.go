package repository

import (
	"context"
	"log"
	"services/webapi/models"
	"services/webapi/word"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/assert"
)

// кто писал эту бд
// тесты сломаны яне буду их писать
func initDb() *memdb.MemDB {

	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"word": &memdb.TableSchema{
				Name: "word",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
					"eng_text": &memdb.IndexSchema{
						Name:    "eng_text",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "EngText"},
					},
					"native_text": &memdb.IndexSchema{
						Name:    "native_text",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "NativeText"},
					},
					"transcription": &memdb.IndexSchema{
						Name:    "transcription",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Transcription"},
					},
					"difficulty": &memdb.IndexSchema{
						Name:    "difficulty",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Difficulty"},
					},
					"category_id": &memdb.IndexSchema{
						Name:    "category_id",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "CategoryId"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		log.Println(err)
	}
	return db
}

func insertTestValues(db *memdb.MemDB) {
	str := "test"
	txn := db.Txn(true)
	defer txn.Abort()
	for i := 0; i < 100000; i++ {
		testWord := &models.Word{
			Id:            uuid.Nil.String(),
			EngText:       str + strconv.Itoa(i),
			NativeText:    str + strconv.Itoa(i),
			Transcription: str + strconv.Itoa(i),
			Difficulty:    str + strconv.Itoa(i),
			CategoryId:    uuid.Nil.String(),
		}
		err := txn.Insert("word", testWord)
		if err != nil {
			panic(err)
		}
	}
	txn.Commit()
}

func insertUniqueTestValues(db *memdb.MemDB) {
	str := "test"
	txn := db.Txn(true)
	defer txn.Abort()
	for i := 0; i < 100000; i++ {
		testWord := &models.Word{
			Id:            uuid.New().String(),
			EngText:       str + strconv.Itoa(i),
			NativeText:    str + strconv.Itoa(i),
			Transcription: str + strconv.Itoa(i),
			Difficulty:    str + strconv.Itoa(i),
			CategoryId:    uuid.Nil.String(),
		}
		err := txn.Insert("word", testWord)
		if err != nil {
			panic(err)
		}
	}
	txn.Commit()
}

func TestGetAll(t *testing.T) {
	db := initDb()
	insertUniqueTestValues(db)
	db.Txn(false)
	r := InmemRepository{db: db}
	words := r.GetAll(context.Background())
	assert.Equal(t, 100000, len(words))
}

func TestGetById(t *testing.T) {
	db := initDb()
	//insertUniqueTestValues(db)
	r := InmemRepository{db: db}
	id := uuid.New().String()
	word := &models.Word{
		Id:            id,
		EngText:       "test",
		NativeText:    "тест",
		Transcription: "|əˈkleɪm|",
		Difficulty:    "easy",
		CategoryId:    id,
	}
	word, err := r.Create(context.Background(), word)
	assert.NoError(t, err)
	all := r.GetAll(context.Background())
	log.Println(len(all))
	foundWord := r.GetById(context.Background(), id)
	assert.NotNil(t, foundWord)
}

func TestFailGetById(t *testing.T) {
	db := initDb()
	//insertUniqueTestValues(db)
	r := InmemRepository{db: db}
	id := uuid.New().String()
	//word := &models.Word{
	//	Id:            id,
	//	EngText:       "test",
	//	NativeText:    "тест",
	//	Transcription: "|əˈkleɪm|",
	//	Difficulty:    "easy",
	//	CategoryId:    id,
	//}
	//word, err := r.Create(context.Background(), word)
	//assert.NoError(t, err)
	//all := r.GetAll(context.Background())
	//log.Println(len(all))
	foundWord := r.GetById(context.Background(), id)
	assert.Nil(t, foundWord)
}

func TestSuccessCreate(t *testing.T) {
	db := initDb()
	r := InmemRepository{db: db}
	word := &models.Word{
		Id:            uuid.New().String(),
		EngText:       "test",
		NativeText:    "тест",
		Transcription: "|əˈkleɪm|",
		Difficulty:    "easy",
		CategoryId:    uuid.New().String(),
	}
	word, err := r.Create(context.Background(), word)
	assert.Nil(t, err)
}

func TestFailCreate(t *testing.T) {
	db := initDb()
	r := InmemRepository{db: db}
	id := uuid.New().String()
	word := &models.Word{
		Id:            id,
		EngText:       "test",
		NativeText:    "тест",
		Transcription: "|əˈkleɪm|",
		Difficulty:    "easy",
		CategoryId:    id,
	}
	word2 := &models.Word{
		Id:            id,
		EngText:       "test",
		NativeText:    "тест",
		Transcription: "|əˈkleɪm|",
		Difficulty:    "easy",
		CategoryId:    id,
	}
	word, err := r.Create(context.Background(), word)
	assert.Nil(t, err)
	word, err2 := r.Create(context.Background(), word2)
	assert.Error(t, err2)
	log.Println(err2)
}

func TestSuccessUpdate(t *testing.T) {
	db := initDb()
	r := InmemRepository{db: db}
	id := uuid.New().String()
	newWord := &models.Word{
		Id:            id,
		EngText:       "test",
		NativeText:    "тест",
		Transcription: "test",
		Difficulty:    "easy",
		CategoryId:    id,
	}
	_, err := r.Create(context.Background(), newWord)
	assert.NoError(t, err)
	words := r.GetAll(context.Background())
	assert.Equal(t, 1, len(words))
	updatedWord := &models.Word{
		Id:            id,
		EngText:       "test2",
		NativeText:    "тест2",
		Transcription: "test2",
		Difficulty:    "easy2",
		CategoryId:    id,
	}
	_, err = r.UpdateById(context.Background(), id, updatedWord)
	assert.NoError(t, err)
	words = r.GetAll(context.Background())
	assert.Equal(t, 1, len(words))
	assert.Equal(t, "test2", words[0].EngText)
	log.Println(words[0])
}

func TestFailUpdate(t *testing.T) {
	db := initDb()
	r := InmemRepository{db: db}
	id := uuid.New().String()
	newWord := &models.Word{
		Id:            id,
		EngText:       "test",
		NativeText:    "тест",
		Transcription: "test",
		Difficulty:    "easy",
		CategoryId:    id,
	}
	_, err := r.Create(context.Background(), newWord)
	assert.NoError(t, err)
	words := r.GetAll(context.Background())
	assert.Equal(t, 1, len(words))
	updatedWord := &models.Word{
		Id:            id,
		EngText:       "test2",
		NativeText:    "тест2",
		Transcription: "test2",
		Difficulty:    "easy2",
		CategoryId:    id,
	}
	_, err = r.UpdateById(context.Background(), "879324789324789", updatedWord)
	assert.Error(t, err)
	assert.IsType(t, word.ErrWordNotFound, err)
}

func TestSuccessDelete(t *testing.T) {
	db := initDb()
	r := InmemRepository{db: db}
	id := uuid.New().String()
	newWord := &models.Word{
		Id:            id,
		EngText:       "test",
		NativeText:    "тест",
		Transcription: "test",
		Difficulty:    "easy",
		CategoryId:    id,
	}
	_, err := r.Create(context.Background(), newWord)
	assert.NoError(t, err)
	words := r.GetAll(context.Background())
	assert.Equal(t, 1, len(words))

	err = r.Delete(context.Background(), id)
	assert.NoError(t, err)
	words = r.GetAll(context.Background())
	assert.Equal(t, 0, len(words))
	//assert.IsType(t, word.ErrWordNotFound, err)
}

func TestFailDelete(t *testing.T) {
	db := initDb()
	r := InmemRepository{db: db}
	id := uuid.New().String()
	newWord := &models.Word{
		Id:            id,
		EngText:       "test",
		NativeText:    "тест",
		Transcription: "test",
		Difficulty:    "easy",
		CategoryId:    id,
	}
	_, err := r.Create(context.Background(), newWord)
	assert.NoError(t, err)
	words := r.GetAll(context.Background())
	assert.Equal(t, 1, len(words))

	err = r.Delete(context.Background(), "21124214")
	assert.Error(t, err)
	words = r.GetAll(context.Background())
	assert.Equal(t, 1, len(words))
	//assert.IsType(t, word.ErrWordNotFound, err)
}
