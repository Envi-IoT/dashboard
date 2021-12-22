package models

//Dustbin data model
type HumidityTemp struct {
	Device      int     `json:"device"`
	Level       float64 `json:"Level"`
	Temperature float64 `json:"Temperature"`
	Humidity    float64 `json:"Humidity"`
}
