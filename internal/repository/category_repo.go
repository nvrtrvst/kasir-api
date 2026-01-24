package repository

import "kasir-api/internal/model" // Sesuaikan dengan nama modul di go.mod kamu

// Pindahkan variabel data asli kodemu ke sini
var Category = []model.Category{
	{ID: 1, Name: "Makanan", Description: "Macam-macam Makanan"},
	{ID: 2, Name: "Minuman", Description: "Macam Minuman"},
	{ID: 3, Name: "Alat Mandi", Description: "Macam Alat Mandi"},
	{ID: 4, Name: "Rokok", Description: "Macam Rokok"},
}
