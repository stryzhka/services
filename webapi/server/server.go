package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
	"webapi/category"
	categoryrepo "webapi/category/repository"
	categorysvc "webapi/category/service"
	categoryhttp "webapi/category/transport/http"
	"webapi/models"
	wordpkg "webapi/word"
	wordrepo "webapi/word/repository"
	wordsvc "webapi/word/service"
	wordhttp "webapi/word/transport/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-memdb"
	httpSwagger "github.com/swaggo/http-swagger"
)

func initWordDb() *memdb.MemDB {

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
			"category": &memdb.TableSchema{
				Name: "category",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
					"name": &memdb.IndexSchema{
						Name:    "name",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Name"},
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

func insertTestValues(db *memdb.MemDB, count int) {
	str := "test"
	txn := db.Txn(true)
	defer txn.Abort()
	for i := 0; i < count; i++ {
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

type App struct {
	wordService     wordpkg.Service
	categoryService category.Service
	server          *http.Server
}

func NewApp() *App {
	db := initWordDb()
	insertTestValues(db, 3)
	wordRepository := wordrepo.NewInmemRepository(db)
	wordService := wordsvc.NewWordService(wordRepository)

	categoryRepository := categoryrepo.NewInmemRepository(db)
	categoryService := categorysvc.NewCategoryService(categoryRepository)
	return &App{
		wordService:     wordService,
		categoryService: categoryService,
	}
}

func (a *App) Run(port string) error {
	wordHandler := wordhttp.NewHandler(a.wordService)
	categoryHandler := categoryhttp.NewHandler(a.categoryService)
	router := mux.NewRouter()

	// word routes
	router.HandleFunc("/api/words/", wordHandler.GetAll).Methods("GET")
	router.HandleFunc("/api/words/{id}", wordHandler.GetById).Methods("GET")
	router.HandleFunc("/api/words/", wordHandler.Create).Methods("POST")
	router.HandleFunc("/api/words/{id}", wordHandler.Delete).Methods("DELETE")
	router.HandleFunc("/api/words/{id}", wordHandler.Update).Methods("PUT")

	// category routes
	router.HandleFunc("/api/categories/", categoryHandler.GetAll).Methods("GET")
	router.HandleFunc("/api/categories/{id}", categoryHandler.GetById).Methods("GET")
	router.HandleFunc("/api/categories/", categoryHandler.Create).Methods("POST")
	router.HandleFunc("/api/categories/{id}", categoryHandler.Delete).Methods("DELETE")
	router.HandleFunc("/api/categories/{id}", categoryHandler.Update).Methods("PUT")
	router.HandleFunc("/api/categories/{id}/words", categoryHandler.GetAllWords).Methods("GET")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	a.server = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			log.Fatalf("Server error: %s", err.Error())
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()
	return a.server.Shutdown(ctx)

}
