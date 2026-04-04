package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"webapi/category"
	categoryrepo "webapi/category/repository"
	"webapi/category/service"
	http2 "webapi/category/transport/http"
	wordpkg "webapi/word"
	"webapi/word/cache"
	"webapi/word/kafka"
	wordrepo "webapi/word/repository"
	wordsvc "webapi/word/service"
	wordhttp "webapi/word/transport/http"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

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

func initKafka() (*kafka.Producer, *kafka.Consumer) {
	connString := os.Getenv("kafka_uri")
	log.Println("kafka_uri", connString)
	p := kafka.NewProducer([]string{connString}, "obj-to-users")
	c := kafka.NewConsumer([]string{connString}, "users-to-obj")
	return p, c
}

type App struct {
	wordService     wordpkg.Service
	categoryService category.Service
	server          *http.Server
	mongoClient     *mongo.Client
	redisClient     *redis.Client
	kafkaProducer   *kafka.Producer
	kafkaConsumer   *kafka.Consumer
}

func NewApp() *App {
	//db := initWordDb()
	client, err := initMongo()
	redisClient, _ := initRedis()
	kafkaProducer, kafkaConsumer := initKafka()
	if err != nil {
		panic(err)
	}
	//insertTestValues(db, 3)
	//wordRepository := wordrepo.NewInmemRepository(db)
	redis := cache.NewRedisWordCache(redisClient)
	wordRepository := wordrepo.NewMongoWordRepository(client, redis)
	wordService := wordsvc.NewWordService(wordRepository, kafkaProducer)

	categoryRepository := categoryrepo.NewMongoCategoryRepository(client)
	categoryService := service.NewCategoryService(categoryRepository)
	return &App{
		wordService:     wordService,
		mongoClient:     client,
		categoryService: categoryService,
		kafkaConsumer:   kafkaConsumer,
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
	// контекст для консумера — живёт всё время работы приложения
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := a.kafkaConsumer.Consume(ctx, func(key, value []byte) error {
			var msg wordpkg.WordConfirmedSuccessMessage
			if err := json.Unmarshal(value, &msg); err != nil {
				return err
			}
			log.Println("msg: ", msg.ObjectId)

			// отдельный контекст для каждого сообщения с таймаутом
			msgCtx, msgCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer msgCancel()

			err := a.wordService.ConfirmWord(msgCtx, msg.ObjectId)
			if err != nil {
				log.Printf("ConfirmWord error for objectId %s: %v", msg.ObjectId, err)
			}
			return err
		}); err != nil {
			log.Printf("consumer error: %s", err.Error())
			cancel()
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
	if err := a.kafkaProducer.Close(); err != nil {
		log.Println("kafka disconnect error:", err)
	}

	return a.server.Shutdown(ctx)
}
