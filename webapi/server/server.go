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
	"webapi/category/service"
	http2 "webapi/category/transport/http"
	"webapi/models"
	wordpkg "webapi/word"
	"webapi/word/cache"
	wordrepo "webapi/word/repository"
	wordsvc "webapi/word/service"
	wordhttp "webapi/word/transport/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-memdb"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func insertTestValues(db *memdb.MemDB, wordsCount int) {
	str := "test"
	txn := db.Txn(true)
	defer txn.Abort()
	categoryNames := []string{"food", "it", "places"}
	for _, v := range categoryNames {
		testCategory := &models.Category{
			Id:   uuid.New().String(),
			Name: v,
		}
		err := txn.Insert("category", testCategory)
		if err != nil {
			panic(err)
		}
		for i := 0; i < wordsCount; i++ {
			testWord := &models.Word{
				Id:            uuid.New().String(),
				EngText:       str + strconv.Itoa(i),
				NativeText:    str + strconv.Itoa(i),
				Transcription: str + strconv.Itoa(i),
				Difficulty:    str + strconv.Itoa(i),
				CategoryId:    testCategory.Id,
			}
			err := txn.Insert("word", testWord)
			if err != nil {
				panic(err)
			}
		}
	}
	txn.Commit()
}

func initMongo() (*mongo.Client, error) {
	connectionString := os.Getenv("MONGODB_URI")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// Defines the options for the MongoDB client
	opts := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI)
	// Creates a new client and connects to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	// Sends a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return nil, err
	}
	log.Println("mongo connected")
	return client, nil
}

func initRedis() (*redis.Client, error) {
	log.Println("redis on", os.Getenv("redis_uri"))
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_uri"),
		Password: os.Getenv("redis_password"),
		DB:       0,
	})
	err := client.Set(context.Background(), "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	} else {
		log.Println("redis connected")
	}
	return client, nil
}

type App struct {
	wordService     wordpkg.Service
	categoryService category.Service
	server          *http.Server
	mongoClient     *mongo.Client
	redisClient     *redis.Client
}

func NewApp() *App {
	//db := initWordDb()
	client, err := initMongo()
	redisClient, _ := initRedis()
	if err != nil {
		panic(err)
	}
	//insertTestValues(db, 3)
	//wordRepository := wordrepo.NewInmemRepository(db)
	redis := cache.NewRedisWordCache(redisClient)
	wordRepository := wordrepo.NewMongoWordRepository(client, redis)
	wordService := wordsvc.NewWordService(wordRepository)

	categoryRepository := categoryrepo.NewMongoCategoryRepository(client)
	categoryService := service.NewCategoryService(categoryRepository)
	return &App{
		wordService:     wordService,
		mongoClient:     client,
		categoryService: categoryService,
	}
}

func (a *App) Run(port string) error {
	wordHandler := wordhttp.NewHandler(a.wordService)
	categoryHandler := http2.NewHandler(a.categoryService)
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
	return a.Shutdown()

}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.mongoClient.Disconnect(ctx); err != nil {
		log.Println("mongo disconnect error:", err)
	}
	if err := a.redisClient.Close(); err != nil {
		log.Println("redis disconnect error:", err)
	}

	return a.server.Shutdown(ctx)
}
