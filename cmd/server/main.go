package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/docs"
	_ "kasir-api/docs"
	"kasir-api/internal/handlers"
	"kasir-api/internal/repositories"
	"kasir-api/internal/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

// @title           Kasir API
// @version         1.0
// @contact.name   Reza Muhammad Akbar
// @contact.url    http://www.kasirapi.com/support
// @contact.email  7bM8A@example.com
// @host
// @BasePath  /
// @schemes         https http
func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// if _, err := os.Stat(".env"); err == nil {
	// 	viper.SetConfigFile(".env")
	// }

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("Error reading .env file:", err)
		}
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	//setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed database ", err)
	}
	defer db.Close()

	// ... pring-print version ...

	// Cek PORT untuk menentukan environment
	port := os.Getenv("PORT")

	if port == "" {
		// SETINGAN LOKAL
		port = "8080"
		docs.SwaggerInfo.Host = "localhost:8080"
		docs.SwaggerInfo.Schemes = []string{"http"}
	} else {
		// SETINGAN RAILWAY
		docs.SwaggerInfo.Host = "kasir-api-production-2671.up.railway.app"
		docs.SwaggerInfo.Schemes = []string{"https"}
	}

	docs.SwaggerInfo.BasePath = "/"

	// Jalur untuk membuka UI Swagger
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Setup routes
	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/products/", productHandler.HandleProductByID)
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	// Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout) // POST

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	fmt.Println("Server Running di port:" + config.Port)

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("Gagal Running Server")
	}
}
