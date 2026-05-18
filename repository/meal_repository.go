package repository

import (
	"database/sql"
	"meal-api/models"
)

func SaveMeal(db *sql.DB, meal models.Meal) error {

	query := ` INSERT INTO meals(name, calories, protein) VALUES ($1, $2, $3) RETURNING id `

	err := db.QueryRow( //Use Query Row when we need to return a row.
		query,
		meal.Name,
		meal.Calories,
		meal.Protein,
	).Scan(&meal.ID) //ID is returned and saved into the meal object.

	return err
}

func FetchMeals(db *sql.DB) ([]models.Meal, error) {
	//return type is the Slice of the Meal struct
	rows, err := db.Query(
		"SELECT id, name, calories, protein FROM meals",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close() //Later close the rows query

	var meals []models.Meal

	for rows.Next() {

		var meal models.Meal //For each row save it in a single meal object

		err := rows.Scan(
			&meal.ID, //It should be in the same order as Query
			&meal.Name,
			&meal.Calories,
			&meal.Protein,
		)

		if err != nil {
			return nil, err
		}

		meals = append(meals, meal)
	}

	return meals, nil
}
