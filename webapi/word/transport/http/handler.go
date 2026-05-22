package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"webapi/models"
	"webapi/word"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type errorMessage struct {
	Error string `json:"error"`
}

type message struct {
	Message string `json:"message"`
}

type WordDtoIn struct {
	EngText         string `json:"eng_text"`
	NativeText      string `json:"native_text"`
	Transcription   string `json:"transcription"`
	Difficulty      string `json:"difficulty"`
	CategoryId      string `json:"category_id"`
	ConfirmedUserId string `json:"confirmed_user_id"`
}

type WordDtoUpdate struct {
	EngText       string `json:"eng_text"`
	NativeText    string `json:"native_text"`
	Transcription string `json:"transcription"`
	Difficulty    string `json:"difficulty"`
	CategoryId    string `json:"category_id"`
}

//type WordDtoOut struct {
//	EngText         string `json:"eng_text"`
//	NativeText      string `json:"native_text"`
//	Transcription   string `json:"transcription"`
//	Difficulty      string `json:"difficulty"`
//	CategoryId      string `json:"category_id"`
//	ConfirmedUserId string `json:"confirmed_user_id"`
//	ConfirmedStatus string `json:"confirmed_status"`
//}

func validateDto(dto *WordDtoIn) error {
	if strings.TrimSpace(dto.EngText) == "" || strings.TrimSpace(dto.NativeText) == "" || strings.TrimSpace(dto.Transcription) == "" || strings.TrimSpace(dto.Difficulty) == "" || strings.TrimSpace(dto.ConfirmedUserId) == "" {
		return word.ErrValidation
	}
	return nil
}

func validateUpdate(dto *WordDtoUpdate) error {
	if strings.TrimSpace(dto.EngText) == "" || strings.TrimSpace(dto.NativeText) == "" || strings.TrimSpace(dto.Transcription) == "" || strings.TrimSpace(dto.Difficulty) == "" {
		return word.ErrValidation
	}
	return nil
}

type Handler struct {
	s word.Service
}

func NewHandler(s word.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Healthcheck(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// GetAll godoc
// @Summary Get all words
// @Tags word
// @Produce json
// @Param difficulty query string false "Filter by difficulty"
// @Success 200 {array} models.Word
// @Router /api/words/ [get]
func (h *Handler) GetAll(w http.ResponseWriter, req *http.Request) {
	var words []*models.Word
	decoder := schema.NewDecoder()
	filter := word.WordFilter{Difficulty: ""}
	err := decoder.Decode(&filter, req.URL.Query())
	if err != nil {
		log.Println(err)
	}
	words = h.s.GetAll(req.Context(), filter)
	jsonWords, err := json.Marshal(words)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if len(words) == 0 || err != nil {
		jsonWords = []byte("[]")
	}
	w.Write(jsonWords)
	return
}

// GetById godoc
// @Summary Get word by id
// @Tags word
// @Produce json
// @Param id path string true "word id"
// @Success 200 {object} models.Word
// @Router /api/words/{id} [get]
func (h *Handler) GetById(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	word := h.s.GetById(req.Context(), id)
	jsonWord, err := json.Marshal(word)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if jsonWord == nil || err != nil || word == nil {
		jsonWord = []byte("{}")
	}
	w.Write(jsonWord)
	return
}

// Create godoc
// @Summary Create new word
// @Tags word
// @Security ApiKeyAuth
// @Param word body WordDtoIn required "word"
// @Accepts json WordDtoIn
// @Produce json
// @Success 201
// @Failure 401
// @Failure 422
// @Router /api/words/ [post]
func (h *Handler) Create(w http.ResponseWriter, req *http.Request) {
	createdWord := &WordDtoIn{}
	err := json.NewDecoder(req.Body).Decode(createdWord)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}
	err = validateDto(createdWord)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}
	word, err := h.s.Create(req.Context(), createdWord.EngText, createdWord.NativeText, createdWord.Transcription, createdWord.Difficulty, createdWord.CategoryId, createdWord.ConfirmedUserId)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}
	jsonWord, err := json.Marshal(word)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonWord)
	return
}

// Delete godoc
// @Summary Delete word by id
// @Tags word
// @Security ApiKeyAuth
// @Param id path string true "word id"
// @Produce json
// @Success 200
// @Failure 401
// @Failure 422
// @Router /api/words/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	err := h.s.Delete(req.Context(), id)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}
	w.WriteHeader(http.StatusOK)
	msg, _ := json.Marshal(message{"deleted successfully"})
	w.Write(msg)
	return
}

// Update godoc
// @Summary Update word by id
// @Tags word
// @Security ApiKeyAuth
// @Param id path string true "word id"
// @Param word body WordDtoUpdate required "word"
// @Success 200 {object} models.Word
// @Accepts json WordDtoUpdate
// @Produce json
// @Failure 401
// @Failure 422
// @Router /api/words/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	updatedWord := &WordDtoUpdate{}
	err := json.NewDecoder(req.Body).Decode(updatedWord)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}

	err = validateUpdate(updatedWord)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}
	newWord := &models.Word{
		Id:            id,
		EngText:       updatedWord.EngText,
		NativeText:    updatedWord.NativeText,
		Transcription: updatedWord.Transcription,
		Difficulty:    updatedWord.Difficulty,
		CategoryId:    updatedWord.CategoryId,
	}
	word, err := h.s.UpdateById(req.Context(), id, newWord)
	if err != nil {
		log.Println("service error")
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}
	jsonWord, err := json.Marshal(word)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonWord)
	return
}
