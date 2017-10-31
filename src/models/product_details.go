package models

type NutritionDetails struct {
	Body   string `json:"body" bson:"body"`
	Header string `json:"header" bson:"header"`
}

type NutritionNutrients struct {
	Total         string               `json:"total" bson:"total"`
	Label         string               `json:"label" bson:"label"`
	PctDailyValue string               `json:"pct_daily_value" bson:"pct_daily_value"`
	SubCats       []NutritionNutrients `bson:"subcategories" json:"subcategories"`
}

type Nutrition struct {
	Calories             uint32               `json:"calories" bson:"calories"`
	Disclaimer           string               `json:"disclaimer" bson:"disclaimer"`
	ServingSize          string               `json:"serving_size" bson:"serving_size"`
	ServingsPerContainer string               `json:"servings_per_container" bson:"servings_per_container"`
	Nutrients            []NutritionNutrients `json:"nutrients" bson:"nutrients"`
}
