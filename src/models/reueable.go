package models

type Email struct {
	To      string
	Body    string
	Subject string
}

type GeoJson struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}

type Address struct {
	Route         string  `json:"route"`
	Country       string  `json:"country"`
	Location      GeoJson `json:"location"`
	PostalCode    string  `json:"postal_code"`
	StreetNumber  string  `json:"street_number"`
	AdminAreaLvl1 string  `json:"administrative_area_level_1"`
}

type Hours struct {
	From uint16 `json:"from" validate:"ltefield=To"`
	To   uint16 `json:"to"`
}
