package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/internal/models"
	"strings"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		// 1. Tambahkan FOR UPDATE untuk mengunci baris produk (mencegah race condition)
		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1 FOR UPDATE", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("produk dengan ID %d tidak ditemukan", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		// 2. Validasi stok sebelum lanjut
		if stock < item.Quantity {
			return nil, fmt.Errorf("stok produk '%s' tidak cukup (sisa: %d)", productName, stock)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		// 3. Update stok
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// 4. Simpan Header Transaksi
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// 5. BULK INSERT Details (Satu Query untuk semua detail)
	if len(details) > 0 {
		query := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES "
		var values []interface{}
		var placeholders []string

		for i, d := range details {
			n := i * 4
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d)", n+1, n+2, n+3, n+4))
			values = append(values, transactionID, d.ProductID, d.Quantity, d.Subtotal)
		}

		query += strings.Join(placeholders, ",")
		_, err = tx.Exec(query, values...)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) GetReportByDateRange(startDate, endDate time.Time) (*models.SalesReport, error) {
	report := &models.SalesReport{}

	// 1. Hitung Total Revenue dan Total Transaksi
	// created_at di database biasanya TIMESTAMP dengan timezone, jadi perlu hati-hati.
	// Asumsi database menyimpan dalam UTC atau lokal, kita pass range waktu yang sesuai.
	// Query untuk aggregate
	querySummary := `
		SELECT 
			COALESCE(SUM(total_amount), 0), 
			COUNT(id)
		FROM transactions 
		WHERE created_at >= $1 AND created_at <= $2
	`
	err := repo.db.QueryRow(querySummary, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// 2. Cari Produk Terlaris
	// Join transaction_details dengan transactions untuk filter tanggal, kemudian group by product_id
	// Sebenarnya transaction_details tidak punya tanggal, jadi harus JOIN ke transactions.
	// Tapi kita butuh nama produk juga, ada di products table.
	// Atau jika nama produk disimpan di transaction_details (seperti di struct Model TransactionDetail ada ProductName?), 
    // Tapi di struct TransactionDetail tidak ada field ProductName yang mapped ke DB row secara eksplisit.
    // Di func CreateTransaction: `INSERT INTO transaction_details ...` tidak menyimpan nama produk.
    // Jadi harus JOIN ke tabel products untuk dapat nama.

	queryBestSeller := `
		SELECT 
			p.name, 
			COALESCE(SUM(td.quantity), 0) as qty_terjual
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE t.created_at >= $1 AND t.created_at <= $2
		GROUP BY p.id, p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`
	
	err = repo.db.QueryRow(queryBestSeller, startDate, endDate).Scan(&report.ProdukTerlaris.Nama, &report.ProdukTerlaris.QtyTerjual)
	if err == sql.ErrNoRows {
		// Tidak ada transaksi, biarkan kosong atau set default
		report.ProdukTerlaris = models.ProductSales{Nama: "-", QtyTerjual: 0}
	} else if err != nil {
		return nil, err
	}

	return report, nil
}
