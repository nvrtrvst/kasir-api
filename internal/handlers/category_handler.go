package handlers

import (
	"encoding/json"
	"kasir-api/internal/models"
	"kasir-api/internal/services"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetAllCategories godoc
// @Summary      Dapatkan Semua Kategori
// @Description  Mengambil semua data kategori yang tersedia
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Success      200  {array}  object{id=int,name=string,description=string}
// @Router       /api/categories [get]
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// CategoryHandler godoc
// @Summary      Dapatkan  Kategori By ID
// @Description  Mengambil data kategori yang tersedia by ID
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param id path int true "Category ID"
// @Success      2200  {array}   object{id=int,name=string,description=string}
// @Router       /api/categories/{id} [get]
func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r)
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category id", http.StatusBadRequest)
		return
	}

	category, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// CreateCategory godoc
// @Summary      Buat Kategori Baru
// @Description  Menambahkan kategori baru ke dalam sistem
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        category  body  object{name=string,description=string}  true  "Data Kategori"
// @Success      201       {object}  object{id=int,name=string,description=string}
// @Router       /api/categories [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err = h.service.Create(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(category)
}

// UpdateCategory godoc
// @Summary      Update Kategori
// @Description  Memperbarui informasi kategori yang ada
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id        path    int     true  "Category ID"
// @Param        category  body  object{name=string,description=string}  true  "Data Kategori"
// @Success      200       {object}  object{id=int,name=string,description=string}
// @Router       /api/categories/{id} [put]
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category id", http.StatusBadRequest)
		return
	}

	var updateData models.Category
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Invalid", http.StatusBadRequest)
		return
	}

	updateData.ID = id
	err = h.service.Update(&updateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateData)
	http.Error(w, "Category Tidak Ada", http.StatusNotFound)
}

// DeleteCategory godoc
// @Summary      Hapus Kategori
// @Description  Menghapus kategori yang ada
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param id path int true "Category ID"
// @Success      200       {object}  object{message=string}
// @Router       /api/categories/{id} [delete]
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category id", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Category deleted successfully",
	})
}
