package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"services/webapi/word"
	"services/webapi/word/repository"
	"services/webapi/word/service"
	http2 "services/webapi/word/transport/http"
	"time"

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
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		log.Println(err)
	}
	return db
}

type App struct {
	wordService word.Service
	server      *http.Server
}

func NewApp() *App {
	db := initWordDb()
	wordRepository := repository.NewInmemRepository(db)
	wordService := service.NewWordService(wordRepository)
	return &App{
		wordService: wordService,
	}
}

func (a *App) Run(port string) error {
	wordHandler := http2.NewHandler(a.wordService)
	wordRouter := mux.NewRouter()
	wordRouter.HandleFunc("/api/words/", wordHandler.GetAll).Methods("GET")
	wordRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	a.server = &http.Server{
		Addr:           ":" + port,
		Handler:        wordRouter,
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
