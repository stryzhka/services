package http

import (
	"encoding/json"
	"net/http"
	"services/webapi/models"
	"services/webapi/word"
	"strings"
)

type errorMessage struct {
	Error string `json:"error"`
}

type message struct {
	Message string `json:"message"`
}

type WordDto struct {
	EngText       string `json:"eng_text"`
	NativeText    string `json:"native_text"`
	Transcription string `json:"transcription"`
	Difficulty    string `json:"difficulty"`
	CategoryId    string `json:"category_id"` //TODO тут по другому тоже но мне лень
}

func validateDto(dto *WordDto) error {
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
// @Tags profile
// @Produce json
// @Success 200 {array} models.Word
// // @Failure 401
// @Router /api/words/ [get]
func (h *Handler) GetAll(w http.ResponseWriter, req *http.Request) {
	var words []*models.Word
	words = h.s.GetAll(req.Context())
	jsonWords, err := json.Marshal(words)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if len(words) == 0 || err != nil {
		jsonWords = []byte("[]")
	}
	w.Write(jsonWords)
	return
}
