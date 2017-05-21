package models

type GeoJson struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}

type Hours struct {
	From uint16 `json:"from"`
	To   uint16 `json:"to"`
}
