package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"meal-api/config"
	"net/http"
	"sync"
	"time"
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

		//Go Routine - Calling the fuction concurrently by creating a new thread

		//Starting 2 go routines
		var wg sync.WaitGroup                 //Initialize wait group
		wg.Add(2)                             //Informing that there are 2 goroutines to be executed
		go CalculateHealthScore(newMeal, &wg) //&wg because we need to update the same wait group
		go CalculateProteinCategory(newMeal, &wg)
		wg.Wait()
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(newMeal)

		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
func CalculateHealthScore(meal Meal, wg *sync.WaitGroup) {
	defer wg.Done() //Indicating that this function has been executed
	fmt.Println("Calculating Health Score")
	time.Sleep(3 * time.Second)
	fmt.Println("Health Score Calculated")
}

func CalculateProteinCategory(meal Meal, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Calculating Protein Category")
	time.Sleep(4 * time.Second)
	fmt.Println("Protein Category Calculated")
}

func WriteAuditLog(meal Meal) {

	fmt.Println("Writing audit log for:", meal.Name)
	//go routine will sleep for 3 seconds
	time.Sleep(3 * time.Second)

	fmt.Println("Audit completed for:", meal.Name)
}

func main() {
	conf := config.LoadConfig()
	fmt.Printf("DB Host %s\n", conf.DBHost)
	db = ConnectDB(conf)
	defer db.Close() //Defer ---- Run this function when the main function is executed.
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/meals", mealsHandler)
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
