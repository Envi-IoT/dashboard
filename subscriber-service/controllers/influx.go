package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"

	models "subs/models"
)

//InitializeInfluxDb and set up write api
func InitializeInfluxDb() (api.WriteAPI, influxdb2.Client) {
	clientInflux := influxdb2.NewClientWithOptions("http://influxdb:8086", "telegraf:sensors",
		influxdb2.DefaultOptions().SetBatchSize(10))
	writeAPI := clientInflux.WriteAPI("", "sensors")
	return writeAPI, clientInflux
}

//PushToDb influxdb
func PushToDb(msg string, w api.WriteAPI) {
	sensorData := models.HumidityTemp{}
	json.Unmarshal([]byte(msg), &sensorData)
	p := influxdb2.NewPoint(
		"humidityTemps",
		map[string]string{"Unit": "celcius"},
		map[string]interface{}{
			"messageID":   sensorData.MessageID,
			"sensorID":    sensorData.SensorID,
			"temperature": sensorData.Temperature,
			"humidity":    sensorData.Humidity,
			"timestamp":   sensorData.Timestamp,
		},
		time.Now())
	w.WritePoint(p)
	if w.Errors() != nil {
		fmt.Printf("Write error: %s\n", <-w.Errors())
	}
}

//CloseInfluxClient close client
func CloseInfluxClient(w api.WriteAPI, clientInflux influxdb2.Client) {
	w.Flush()
	clientInflux.Close()
}
