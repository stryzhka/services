package server

import (
	"context"
	"encoding/json"
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
	wordpkg "webapi/word"
	"webapi/word/cache"
	"webapi/word/kafka"
	wordrepo "webapi/word/repository"
	wordsvc "webapi/word/service"
	wordhttp "webapi/word/transport/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
	httpRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Requests currently being processed",
		},
	)
	mongoQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "mongo_query_duration_seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
	kafkaProduced = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kafka_messages_produced_total",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"topic"},
	)
	kafkaConsumed = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kafka_messages_consumed_total",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"topic"},
	)
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
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

func initKafka() (*kafka.Producer, *kafka.Consumer) {
	connString := os.Getenv("kafka_uri")
	log.Println("kafka_uri", connString)
	p := kafka.NewProducer([]string{connString}, "obj-to-users", kafkaProduced)
	c := kafka.NewConsumer([]string{connString}, "users-to-obj", kafkaConsumed)
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
	wordRepository := wordrepo.NewMongoWordRepository(client, redis, mongoQueryDuration)
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

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" || r.URL.Path == "/swagger" {
			next.ServeHTTP(w, r)
			return
		}
		path := r.URL.Path
		if route := mux.CurrentRoute(r); route != nil {
			if tmpl, err := route.GetPathTemplate(); err == nil {
				path = tmpl
			}
		}
		rec := &statusRecorder{ResponseWriter: w, status: 200}
		httpRequestsInFlight.Inc()
		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(r.Method, path))
		next.ServeHTTP(rec, r)
		timer.ObserveDuration()
		httpRequestsInFlight.Dec()
		httpRequestsTotal.WithLabelValues(r.Method, path, strconv.Itoa(rec.status)).Inc()
	})
}

func (a *App) Run(port string) error {
	wordHandler := wordhttp.NewHandler(a.wordService)
	categoryHandler := http2.NewHandler(a.categoryService)
	router := mux.NewRouter()
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		prometheus.NewGoCollector(),
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		httpRequestsTotal,
		httpRequestDuration,
		httpRequestsInFlight,
		mongoQueryDuration,
		kafkaProduced,
		kafkaConsumed)
	router.Use(metricsMiddleware)
	router.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := a.kafkaConsumer.Consume(ctx, func(key, value []byte) error {
			var msg wordpkg.WordConfirmedSuccessMessage
			if err := json.Unmarshal(value, &msg); err != nil {
				return err
			}

			msgCtx, msgCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer msgCancel()

			err := a.wordService.ConfirmWord(msgCtx, msg.ObjectId, msg.ConfirmedAt.String())
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
