package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"

	models "subs/models"
)

//InitializeInfluxDb and set up write api
func InitializeInfluxDb() (api.WriteAPI, influxdb2.Client) {
	clientInflux := influxdb2.NewClientWithOptions("http://influxdb:8086", "telegraf:sensors",
		influxdb2.DefaultOptions().SetBatchSize(10))
	writeAPI := clientInflux.WriteAPI("", "sensors")
	return writeAPI, clientInflux
}

//create TCAN data point
func CreateTCANPoint(msg string) *write.Point {
	sensorData := models.HumidityTemp{}
	json.Unmarshal([]byte(msg), &sensorData)
	p := influxdb2.NewPoint(
		"Tcan",
		map[string]string{"Unit": "celcius"},
		map[string]interface{}{
			"device":      sensorData.Device,
			"Level":       sensorData.Level,
			"temperature": sensorData.Temperature,
			"humidity":    sensorData.Humidity,
		},
		time.Now())
	return p
}

//create health data point
func CreateHealthPoint(msg string) *write.Point {
	data := models.Health{}
	json.Unmarshal([]byte(msg), &data)
	p := influxdb2.NewPoint(
		"Health",
		map[string]string{"Unit": "volts"},
		map[string]interface{}{
			"device":    data.Device,
			"battery":   data.Battery,
			"bootCount": data.BootCount,
		},
		time.Now())
	return p
}

//PushToDb influxdb
func PushToDb(p *write.Point, w api.WriteAPI) {
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
