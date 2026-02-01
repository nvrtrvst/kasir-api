package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Open database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	//Test Conn
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set Connection pool settings (opstional tp reccomencd)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Sukses DB")
	return db, nil
}
