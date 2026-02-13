package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"webapi/category"
	"webapi/models"

	"github.com/gorilla/mux"
)

type errorMessage struct {
	Error string `json:"error"`
}

type message struct {
	Message string `json:"message"`
}

type CategoryDto struct {
	Name string `json:"name"`
}

func validateDto(dto *CategoryDto) error {
	if strings.TrimSpace(dto.Name) == "" {
		return category.ErrValidation
	}
	return nil
}

type Handler struct {
	s category.Service
}

func NewHandler(s category.Service) *Handler {
	return &Handler{s: s}
}

// GetAll godoc
// @Summary Get all categories
// @Tags category
// @Produce json
// @Success 200 {array} models.Category
// @Router /api/categories/ [get]
func (h *Handler) GetAll(w http.ResponseWriter, req *http.Request) {
	var categories []*models.Category
	categories = h.s.GetAll(req.Context())
	jsonCategories, err := json.Marshal(categories)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if len(categories) == 0 || err != nil {
		jsonCategories = []byte("[]")
	}
	w.Write(jsonCategories)
}

// GetById godoc
// @Summary Get category by id
// @Tags category
// @Produce json
// @Param id path string true "category id"
// @Success 200 {object} models.Category
// @Router /api/categories/{id} [get]
func (h *Handler) GetById(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	categoryModel := h.s.GetById(req.Context(), id)
	jsonCategory, err := json.Marshal(categoryModel)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if jsonCategory == nil || err != nil || categoryModel == nil {
		jsonCategory = []byte("{}")
	}
	w.Write(jsonCategory)
}

// GetAllWords godoc
// @Summary Get all words by category id
// @Tags category
// @Produce json
// @Param id path string true "category id"
// @Success 200 {array} models.Word
// @Router /api/categories/{id}/words [get]
func (h *Handler) GetAllWords(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	var words []*models.Word
	words = h.s.GetAllWords(req.Context(), id)
	jsonWords, err := json.Marshal(words)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if len(words) == 0 || err != nil {
		jsonWords = []byte("[]")
	}
	w.Write(jsonWords)
}

// Create godoc
// @Summary Create new category
// @Tags category
// @Param category body CategoryDto required "category"
// @Accepts json CategoryDto
// @Produce json
// @Success 201
// @Failure 422
// @Router /api/categories/ [post]
func (h *Handler) Create(w http.ResponseWriter, req *http.Request) {
	createdCategory := &CategoryDto{}
	err := json.NewDecoder(req.Body).Decode(createdCategory)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}
	err = validateDto(createdCategory)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}
	categoryModel, err := h.s.Create(req.Context(), createdCategory.Name)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}
	jsonCategory, err := json.Marshal(categoryModel)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonCategory)
}

// Delete godoc
// @Summary Delete category by id
// @Tags category
// @Param id path string true "category id"
// @Produce json
// @Success 200
// @Failure 422
// @Router /api/categories/{id} [delete]
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
}

// Update godoc
// @Summary Update category by id
// @Tags category
// @Param id path string true "category id"
// @Param category body CategoryDto required "category"
// @Success 200 {object} models.Category
// @Accepts json CategoryDto
// @Produce json
// @Failure 422
// @Router /api/categories/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	updatedCategory := &CategoryDto{}
	err := json.NewDecoder(req.Body).Decode(updatedCategory)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}
	err = validateDto(updatedCategory)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}
	newCategory := &models.Category{
		Id:   id,
		Name: updatedCategory.Name,
	}
	categoryModel, err := h.s.UpdateById(req.Context(), id, newCategory)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}
	jsonCategory, err := json.Marshal(categoryModel)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		msg, _ := json.Marshal(errorMessage{err.Error()})
		w.Write(msg)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonCategory)
}
