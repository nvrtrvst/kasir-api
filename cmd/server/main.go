package main

import (
	"encoding/json"
	"fmt"
	_ "kasir-api/docs"
	"kasir-api/internal/handler"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Kasir API
// @version         1.0
// @contact.name   Reza Muhammad Akbar
// @contact.url    http://www.kasirapi.com/support
// @contact.email  7bM8A@example.com
// @host      kasir-api-production-2671.up.railway.app
func main() {

	// Jalur untuk membuka UI Swagger
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method == "GET" {
	// 		handler.GetProdukByID(w, r)
	// 	} else if r.Method == "PUT" {
	// 		handler.UpdateProduk(w, r)
	// 	} else if r.Method == "DELETE" {
	// 		handler.DeleteProduk(w, r)
	// 	}
	// })

	// Endpoint untuk mendapatkan semua produk dan menambahkan produk baru
	// http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method == "GET" {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		json.NewEncoder(w).Encode(repository.Produk)
	// 	} else if r.Method == "POST" {
	// 		var produkBaru model.Produk
	// 		json.NewDecoder(r.Body).Decode(&produkBaru)
	// 		produkBaru.ID = len(repository.Produk) + 1
	// 		repository.Produk = append(repository.Produk, produkBaru)

	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusCreated)
	// 		json.NewEncoder(w).Encode(produkBaru)
	// 	}
	// })

	// http.HandleFunc("/api/categories//", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method == "GET" {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		json.NewEncoder(w).Encode(repository.Category)
	// 	}
	// })

	productHandler := handler.NewProductHandler()
	http.HandleFunc("GET /api/produk", productHandler.GetAllProducts)
	http.HandleFunc("POST /api/produk", productHandler.CreateProduct)
	http.HandleFunc("PUT /api/produk/{id}", productHandler.UpdateProduct)
	http.HandleFunc("GET /api/produk/{id}", productHandler.GetProductByID)
	http.HandleFunc("DELETE /api/produk/{id}", productHandler.DeleteProduct)

	categoryHandler := handler.NewCategoryHandler()
	http.HandleFunc("GET /api/categories", categoryHandler.GetAllCategories)
	http.HandleFunc("POST /api/categories", categoryHandler.CreateCategory)
	http.HandleFunc("PUT /api/categories/{id}", categoryHandler.UpdateCategory)
	http.HandleFunc("GET /api/categories/{id}", categoryHandler.GetCategoryByID)
	http.HandleFunc("DELETE /api/categories/{id}", categoryHandler.DeleteCategory)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	fmt.Println("Server Running di port 8080")
	http.ListenAndServe(":8080", nil)
}
