package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"meal-api/config"
	"meal-api/models"
	"meal-api/repository"
	"meal-api/services"
	"net/http"
	"sync"
)

var db *sql.DB // Provided by go as Reference/handle to database connection manager
// Pass it to every DB function to not create a new db connection - that's why *
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Meal API is running"))
}

func mealsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {

		meals, err := repository.FetchMeals(db)

		if err != nil {
			http.Error(w, "Failed to fetch meals", http.StatusInternalServerError)
			return
		}
		//Serialize to convert into Json object
		json.NewEncoder(w).Encode(meals)
	}

	if r.Method == http.MethodPost {

		newMeal := models.Meal{}

		err := json.NewDecoder(r.Body).Decode(&newMeal)

		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = repository.SaveMeal(db, newMeal)

		if err != nil {
			http.Error(w, "Failed to create meal", http.StatusInternalServerError)
			return
		}

		//Go Routine - Calling the fuction concurrently by creating a new thread

		//Starting 2 go routines
		var wg sync.WaitGroup //Initialize wait group
		wg.Add(2)             //Informing that there are 2 goroutines to be executed
		go services.CalculateHealthScore(newMeal, &wg)
		go services.CalculateProteinCategory(newMeal, &wg) //&wg because we need to update the same wait group
		wg.Wait()
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(newMeal)

		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {
	conf, err := config.MustLoad()
	if err == nil {
		fmt.Printf("DB Host %s\n", conf.DBHost)
		db = ConnectDB(*conf) //db is defined as a global var of *sql.db type
		defer db.Close()      //Defer ---- Run this function when the main function is executed.
		http.HandleFunc("/health", healthHandler)
		http.HandleFunc("/meals", mealsHandler)
		fmt.Println("Server started on port 8080")
		http.ListenAndServe(":8080", nil)
	}
}
