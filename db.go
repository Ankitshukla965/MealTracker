package main

import (
	"database/sql"
	"fmt"
	"log"
	"meal-api/config"

	_ "github.com/lib/pq" //run go get github.com/lib/pq to install this
)
//config.Config - From config package, Config struct
func ConnectDB(cfg config.Config) *sql.DB {
	connStr := fmt.Sprintf( //Sprintf can be saved in a variable
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
	)
	db, err := sql.Open("postgres", connStr) //sql.Open to connect with DB

	if err != nil {
		log.Fatal("Error opening DB:", err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	log.Println("Connected to PostgreSQL")

	return db
}
