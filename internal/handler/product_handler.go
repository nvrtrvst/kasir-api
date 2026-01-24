package handler

import (
	"encoding/json"
	"kasir-api/internal/model"
	"kasir-api/internal/repository" // Ambil data slice dari sini
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

// GetAllProducts godoc
// @Summary      Dapatkan Semua Produk
// @Description  Mengambil semua data produk yang tersedia
// @Tags         Products
// @Accept       json
// @Produce      json
// @Success      2200  {array}  object{id=int,nama=string,harga=int,stok=int}
// @Router       /api/produk [get]
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repository.Produk) // Pakai repository.Produk
}

// GetProdukByID godoc
// @Summary      Dapatkan Produk By ID
// @Description  Mengambil data produk yang tersedia by ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param id path int true "Produk ID"
// @Success      2200  {array}   object{id=int,nama=string,harga=int,stok=int}
// @Router       /api/produk/{id} [get]
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk id", http.StatusBadRequest)
		return
	}

	for _, p := range repository.Produk { // Pakai repository.Produk
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, "Produk Tidak Ada", http.StatusNotFound)
}

// CreateProduct godoc
// @Summary      Buat Produk Baru
// @Description  Menambahkan produk baru ke dalam sistem
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        produk  body  object{nama=string,harga=int,stok=int}  true  "Data Produk"
// @Success      201       {object}  object{id=int,nama=string,harga=int,stok=int}
// @Router       /api/produk [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var newProduk model.Produk
	err := json.NewDecoder(r.Body).Decode(&newProduk)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	//masukkan data kedalam variable produk
	newProduk.ID = len(repository.Produk) + 1
	repository.Produk = append(repository.Produk, newProduk)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(newProduk)
}

// UpdateProduct godoc
// @Summary      Update Produk
// @Description  Memperbarui informasi produk yang ada
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id        path    int     true  "Produk ID"
// @Param        produk  body  object{nama=string,harga=int,stok=int}  true  "Data Produk"
// @Success      200       {object}  object{id=int,nama=string,harga=int,stok=int}
// @Router       /api/produk/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid Produk id", http.StatusBadRequest)
		return
	}

	var updateData model.Produk
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Invalid", http.StatusBadRequest)
		return
	}

	for i := range repository.Produk {
		if repository.Produk[i].ID == id {
			updateData.ID = id
			repository.Produk[i] = updateData
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateData)
			return
		}
	}
	http.Error(w, "Produk Tidak Ada", http.StatusNotFound)
}

// DeleteProduk godoc
// @Summary      Hapus Produk
// @Description  Menghapus produk yang ada
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param id path int true "Produk ID"
// @Success      200       {object}  object{message=string}
// @Router       /api/produk/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk id", http.StatusBadRequest)
		return
	}

	for i, p := range repository.Produk {
		if p.ID == id {
			repository.Produk = append(repository.Produk[:i], repository.Produk[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Sukses Delete"})
			return
		}
	}
	http.Error(w, "Produk Tidak Ada", http.StatusNotFound)
}
