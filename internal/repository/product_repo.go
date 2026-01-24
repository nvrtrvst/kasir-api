package repository

import "kasir-api/internal/model" // Sesuaikan dengan nama modul di go.mod kamu

// Pindahkan variabel data asli kodemu ke sini
var Produk = []model.Produk{
	{ID: 1, Nama: "indomie Garung", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Mizone", Harga: 5500, Stok: 5},
	{ID: 3, Nama: "Kecap", Harga: 10000, Stok: 8},
	{ID: 4, Nama: "Masako", Harga: 2500, Stok: 20},
}
