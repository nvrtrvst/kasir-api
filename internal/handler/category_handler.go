package handler

import (
	"encoding/json"
	"kasir-api/internal/model"
	"kasir-api/internal/repository"
	"net/http"
	"strconv"
)

type CategoryHandler struct {
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

// GetAllCategories godoc
// @Summary      Dapatkan Semua Kategori
// @Description  Mengambil semua data kategori yang tersedia
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Success      2200  {array}  object{id=int,name=string,description=string}
// @Router       /api/categories [get]
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repository.Category)
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
func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
		return
	}
	for _, c := range repository.Category {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}
	http.Error(w, "Kategori tidak ada", http.StatusNotFound)
}

// CreateCategory godoc
// @Summary      Buat Kategori Baru
// @Description  Menambahkan kategori baru ke dalam sistem
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id path int true "Category ID"
// @Success      201       {object}  object{id=int,name=string,description=string}
// @Router       /api/categories [post]
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory model.Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	//masukkan data kedalam variable kategori
	newCategory.ID = len(repository.Category) + 1
	repository.Category = append(repository.Category, newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(newCategory)
}

// UpdateCategory godoc
// @Summary      Update Kategori
// @Description  Memperbarui informasi kategori yang ada
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param id path int true "Category ID"
// @Success      200       {object}  object{id=int,name=string,description=string}
// @Router       /api/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid Category id", http.StatusBadRequest)
		return
	}

	var updateData model.Category
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Invalid", http.StatusBadRequest)
		return
	}

	for i := range repository.Category {
		if repository.Category[i].ID == id {
			updateData.ID = id
			repository.Category[i] = updateData
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateData)
			return
		}
	}
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
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid Category id", http.StatusBadRequest)
		return
	}

	for i, c := range repository.Category {
		if c.ID == id {
			repository.Category = append(repository.Category[:i], repository.Category[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Sukses Delete Category"})
			return
		}
	}
	http.Error(w, "Data tidak ada", http.StatusNotFound)
}
