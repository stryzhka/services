package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"testing"
	"webapi/models"
	http2 "webapi/word/transport/http"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAdd100(t *testing.T) {
	client := &http.Client{}
	for i := 0; i < 101; i++ {
		word := &http2.WordDto{
			EngText:       "test|" + strconv.Itoa(i),
			NativeText:    "тест",
			Transcription: "@@@22!!",
			Difficulty:    "easy",
			CategoryId:    uuid.Nil.String(),
		}
		jsonWord, _ := json.Marshal(word)
		req, err := http.NewRequest("POST", "http://localhost:8080/api/words/", bytes.NewReader(jsonWord))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		log.Printf("Отправка %d", i)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	}
}

func TestAdd1000(t *testing.T) {
	client := &http.Client{}
	for i := 0; i < 1001; i++ {
		word := &http2.WordDto{
			EngText:       "test|" + strconv.Itoa(i),
			NativeText:    "тест",
			Transcription: "@@@22!!",
			Difficulty:    "easy",
			CategoryId:    uuid.Nil.String(),
		}
		jsonWord, _ := json.Marshal(word)
		req, err := http.NewRequest("POST", "http://localhost:8080/api/words/", bytes.NewReader(jsonWord))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		log.Printf("Отправка %d", i)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		if resp != nil {
			resp.Body.Close()
			assert.Equal(t, http.StatusCreated, resp.StatusCode)
		}
		//assert.Equal(t, http.StatusCreated, resp.StatusCode)
	}
}

func TestDeleteAll(t *testing.T) {
	client := &http.Client{}

	// GET all words
	req, err := http.NewRequest("GET", "http://localhost:8080/api/words/", nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	if resp == nil {
		t.Fatal("response is nil")
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read body correctly
	jsonWords, err := io.ReadAll(resp.Body) // ← правильное чтение тела
	assert.NoError(t, err)

	var words []models.Word
	err = json.Unmarshal(jsonWords, &words)
	assert.NoError(t, err)

	// Delete each word
	for _, word := range words {
		req, err := http.NewRequest("DELETE", "http://localhost:8080/api/words/"+word.Id, nil)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		if resp != nil {
			resp.Body.Close()
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	}
}
