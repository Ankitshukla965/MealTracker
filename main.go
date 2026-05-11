package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"meal-api/config"
	"net/http"
)

var db *sql.DB

type Meal struct {
	ID       int    `json:"id"` // in case of json, change it to id
	Name     string `json:"name"`
	Calories int    `json:"calories"`
	Protein  int    `json:"protein"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Meal API is running"))
}

func mealsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {

		//db.Query will return all the words
		rows, err := db.Query(
			"SELECT id, name, calories, protein FROM meals",
		)

		if err != nil {
			http.Error(w, "Failed to fetch meals", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var meals []Meal

		//Move from row to row (Iteration)
		for rows.Next() {

			var meal Meal

			//copies DB Columns into row variable that is of Meal type defined above

			err := rows.Scan(
				&meal.ID,
				&meal.Name,
				&meal.Calories,
				&meal.Protein,
			)

			if err != nil {
				http.Error(w, "Failed to scan meal", http.StatusInternalServerError)
				return
			}

			meals = append(meals, meal)
		}

		json.NewEncoder(w).Encode(meals)

		return
	}

	if r.Method == http.MethodPost {

		newMeal := Meal{}

		err := json.NewDecoder(r.Body).Decode(&newMeal)

		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		query := `
	INSERT INTO meals (name, calories, protein)
	VALUES ($1, $2, $3)
	RETURNING id
	`

		err = db.QueryRow(
			query,
			newMeal.Name,
			newMeal.Calories,
			newMeal.Protein,
		).Scan(&newMeal.ID)

		if err != nil {
			http.Error(w, "Failed to create meal", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(newMeal)

		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {
	conf := config.LoadConfig()
	fmt.Printf("DB Host %s\n", conf.DBHost)
	db = ConnectDB(conf)
	defer db.Close()
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/meals", mealsHandler)
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
