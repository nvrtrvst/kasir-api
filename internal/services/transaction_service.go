package services

import (
	"kasir-api/internal/models"
	"kasir-api/internal/repositories"
	"time"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetDailyReport(date time.Time) (*models.SalesReport, error) {
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

	return s.repo.GetReportByDateRange(startDate, endDate)
}

func (s *TransactionService) GetReportByDateRange(startDate, endDate time.Time) (*models.SalesReport, error) {
	// Adjust endDate to include the full day if needed, or assume caller handles it.
	// For "YYYY-MM-DD" parsing, usually we get 00:00:00.
	// We should probably set endDate to 23:59:59 of that day.
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
	return s.repo.GetReportByDateRange(startDate, endDate)
}
