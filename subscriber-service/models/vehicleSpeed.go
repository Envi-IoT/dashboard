package models

import "time"

//VehicleMetrics model
type HumidityTemp struct {
	MessageID   int       `json:"messageID"`
	SensorID    int       `json:"sensorID"`
	Temperature int       `json:"temperature"`
	Humidity    string    `json:"humidity"`
	Timestamp   time.Time `json:"timestamp"`
}
