package services

import (
	"fmt"
	"meal-api/models"
	"sync"
	"time"
)

func CalculateHealthScore(meal models.Meal, wg *sync.WaitGroup) {
	defer wg.Done() //Indicating that this function has been executed
	fmt.Println("Calculating Health Score")
	time.Sleep(3 * time.Second)
	fmt.Println("Health Score Calculated")
}

func CalculateProteinCategory(meal models.Meal, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Calculating Protein Category")
	time.Sleep(4 * time.Second)
	fmt.Println("Protein Category Calculated")
}
