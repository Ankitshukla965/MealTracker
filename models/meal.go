package models

type Meal struct {
	ID       int    `json:"id"` // in case of json, change it to id
	Name     string `json:"name"`
	Calories int    `json:"calories"`
	Protein  int    `json:"protein"`
}
