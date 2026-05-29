package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"webapi/models"

	"github.com/stretchr/testify/assert"

	"net/http"
	"strconv"
	"testing"
	dtos "webapi/word/transport/http"
)

//package integration
//
//import (
//	"context"
//	"encoding/json"
//	"io"
//	"net/http"
//	"testing"
//	"time"
//	"webapi/models"
//	"webapi/word"
//	"webapi/word/cache"
//	"webapi/word/repository"
//
//	"github.com/prometheus/client_golang/prometheus"
//	redis2 "github.com/redis/go-redis/v9"
//	"github.com/stretchr/testify/assert"
//	"github.com/testcontainers/testcontainers-go"
//	"github.com/testcontainers/testcontainers-go/modules/mongodb"
//	"github.com/testcontainers/testcontainers-go/modules/redis"
//	"go.mongodb.org/mongo-driver/v2/mongo"
//	"go.mongodb.org/mongo-driver/v2/mongo/options"
//)
//

func TestQueueCreate(t *testing.T) {
	client := &http.Client{}
	for i := 0; i < 100; i++ {
		word := &dtos.WordDtoIn{
			EngText:         "test|" + strconv.Itoa(i),
			NativeText:      "тест",
			Transcription:   "@@@22!!",
			Difficulty:      "easy",
			CategoryId:      "d89eb250-29ff-4000-bc4d-72b39ebb4a5d",
			ConfirmedUserId: "d89eb250-29ff-4000-bc4d-72b39ebb4a5d",
		}
		jsonWord, _ := json.Marshal(word)
		req, err := http.NewRequest("POST", "http://localhost:8080/api/words/", bytes.NewReader(jsonWord))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		log.Printf("Отправка %d", i)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		var msg []byte
		msg, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(msg))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	}
}

func TestQueueConfirmed(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/api/words/", nil)
	if err != nil {
		log.Fatal(err)
	}
	var objects []*models.Word
	var bodyBytes []byte
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyBytes, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bodyBytes, &objects)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		expected := fmt.Sprintf("test|%d", i)
		found := false
		for _, obj := range objects {
			if obj.EngText == expected && obj.ConfirmedStatus == "confirmed" {
				found = true
				break
			}
		}

		assert.True(t, found, "not found")
	}
}

//func TestAdd100(t *testing.T) {
//	//client := &http.Client{}
//	//for i := 0; i < 101; i++ {
//	//	word := &http2.WordDto{
//	//		EngText:       "test|" + strconv.Itoa(i),
//	//		NativeText:    "тест",
//	//		Transcription: "@@@22!!",
//	//		Difficulty:    "easy",
//	//		CategoryId:    uuid.Nil.String(),
//	//	}
//	//	jsonWord, _ := json.Marshal(word)
//	//	req, err := http.NewRequest("POST", "http://localhost:8080/api/words/", bytes.NewReader(jsonWord))
//	//	assert.NoError(t, err)
//	//	req.Header.Set("Content-Type", "application/json")
//	//	log.Printf("Отправка %d", i)
//	//	resp, err := client.Do(req)
//	//	assert.NoError(t, err)
//	//	assert.Equal(t, http.StatusCreated, resp.StatusCode)
//	//}
//}
//
//func TestAdd100000(t *testing.T) {
//	//client := &http.Client{}
//	//for i := 1; i < 100001; i++ {
//	//	word := &http2.WordDto{
//	//		EngText:       "test|" + strconv.Itoa(i),
//	//		NativeText:    "тест",
//	//		Transcription: "@@@22!!",
//	//		Difficulty:    "easy",
//	//		CategoryId:    uuid.Nil.String(),
//	//	}
//	//	jsonWord, _ := json.Marshal(word)
//	//	req, err := http.NewRequest("POST", "http://localhost:8080/api/words/", bytes.NewReader(jsonWord))
//	//	assert.NoError(t, err)
//	//	req.Header.Set("Content-Type", "application/json")
//	//	log.Printf("Отправка %d", i)
//	//	resp, err := client.Do(req)
//	//	assert.NoError(t, err)
//	//	if resp != nil {
//	//		resp.Body.Close()
//	//		assert.Equal(t, http.StatusCreated, resp.StatusCode)
//	//	}
//	//	//assert.Equal(t, http.StatusCreated, resp.StatusCode)
//	//}
//}
//
//func TestDeleteAll(t *testing.T) {
//	client := &http.Client{}
//
//	// GET all words
//	req, err := http.NewRequest("GET", "http://localhost:8080/api/words/", nil)
//	assert.NoError(t, err)
//
//	resp, err := client.Do(req)
//	assert.NoError(t, err)
//	if resp == nil {
//		t.Fatal("response is nil")
//	}
//	defer resp.Body.Close()
//	assert.Equal(t, http.StatusOK, resp.StatusCode)
//
//	// Read body correctly
//	jsonWords, err := io.ReadAll(resp.Body) // ← правильное чтение тела
//	assert.NoError(t, err)
//
//	var words []models.Word
//	err = json.Unmarshal(jsonWords, &words)
//	assert.NoError(t, err)
//
//	// Delete each word
//	for _, word := range words {
//		req, err := http.NewRequest("DELETE", "http://localhost:8080/api/words/"+word.Id, nil)
//		assert.NoError(t, err)
//
//		resp, err := client.Do(req)
//		assert.NoError(t, err)
//		if resp != nil {
//			resp.Body.Close()
//			assert.Equal(t, http.StatusOK, resp.StatusCode)
//		}
//	}
//}
//
//func setupContainers(t *testing.T) (*mongo.Client, word.Cache, func()) {
//	ctx := context.Background()
//
//	// MongoDB контейнер
//	mongoC, err := mongodb.RunContainer(ctx, testcontainers.WithImage("mongo:7"))
//	if err != nil {
//		t.Fatal(err)
//	}
//	mongoURI, _ := mongoC.ConnectionString(ctx)
//	mongoClient, _ := mongo.Connect(options.Client().ApplyURI(mongoURI))
//
//	// Redis контейнер
//	redisC, err := redis.RunContainer(ctx, testcontainers.WithImage("redis:7"))
//	if err != nil {
//		t.Fatal(err)
//	}
//	redisAddr, _ := redisC.ConnectionString(ctx)
//	redisClient := redis2.NewClient(&redis2.Options{
//		Addr: redisAddr,
//	})
//	redisCache := cache.NewRedisWordCache(redisClient) // твоя реализация word.Cache
//
//	// Засеять данные
//	seedWords(t, mongoClient, 100_000)
//
//	cleanup := func() {
//		mongoClient.Disconnect(ctx)
//		mongoC.Terminate(ctx)
//		redisC.Terminate(ctx)
//	}
//
//	return mongoClient, redisCache, cleanup
//}
//
//func seedWords(t *testing.T, client *mongo.Client, count int) {
//	t.Helper()
//	coll := client.Database("words").Collection("words")
//
//	docs := make([]interface{}, count)
//	for i := range docs {
//		docs[i] = models.Word{
//			EngText:         string(rune(i)),
//			NativeText:      string(rune(i)),
//			Transcription:   string(rune(i)),
//			Difficulty:      "1",
//			CategoryId:      "0000",
//			ConfirmedUserId: "0000",
//		}
//	}
//
//	_, err := coll.InsertMany(context.Background(), docs,
//		options.InsertMany().SetOrdered(false),
//	)
//	if err != nil {
//		t.Fatal("seed failed:", err)
//	}
//}
//
//func TestGetAll_CacheVsNoDB(t *testing.T) {
//	mongoClient, redisCache, cleanup := setupContainers(t)
//	defer cleanup()
//
//	ctx := context.Background()
//	hist := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "test"}, []string{"op"})
//	repo := repository.NewMongoWordRepository(mongoClient, redisCache, hist)
//
//	filter := word.WordFilter{} // пустой — берём все
//
//	// --- Первый запрос: кеша нет, идём в Mongo ---
//	start := time.Now()
//	words1 := repo.GetAll(ctx, filter)
//	coldDuration := time.Since(start)
//
//	if len(words1) == 0 {
//		t.Fatal("expected words, got 0")
//	}
//
//	// --- Второй запрос: данные уже в Redis ---
//	start = time.Now()
//	words2 := repo.GetAll(ctx, filter)
//	hotDuration := time.Since(start)
//
//	t.Logf("Mongo (cold): %v | Redis (hot): %v | speedup: %.1fx",
//		coldDuration, hotDuration,
//		float64(coldDuration)/float64(hotDuration),
//	)
//
//	// Кеш должен быть быстрее
//	if hotDuration >= coldDuration {
//		t.Errorf("cache не дал прироста: cold=%v hot=%v", coldDuration, hotDuration)
//	}
//
//	// Данные должны совпадать
//	if len(words1) != len(words2) {
//		t.Errorf("разное кол-во: mongo=%d redis=%d", len(words1), len(words2))
//	}
//}
//
//// Бенчмарк для go test -bench=.
//func BenchmarkGetAll_Cold(b *testing.B) {
//	mongoClient, redisCache, cleanup := setupContainers(&testing.T{})
//	defer cleanup()
//
//	ctx := context.Background()
//	hist := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "bench"}, []string{"op"})
//	repo := repository.NewMongoWordRepository(mongoClient, redisCache, hist)
//	filter := word.WordFilter{}
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		redisCache.InvalidateAll(ctx) // сбрасываем кеш перед каждой итерацией
//		repo.GetAll(ctx, filter)
//	}
//}
//
//func BenchmarkGetAll_Hot(b *testing.B) {
//	mongoClient, redisCache, cleanup := setupContainers(&testing.T{})
//	defer cleanup()
//
//	ctx := context.Background()
//	hist := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "bench"}, []string{"op"})
//	repo := repository.NewMongoWordRepository(mongoClient, redisCache, hist)
//	filter := word.WordFilter{}
//
//	// Прогрев — заполняем кеш
//	repo.GetAll(ctx, filter)
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		repo.GetAll(ctx, filter) // кеш всегда горячий
//	}
//}
