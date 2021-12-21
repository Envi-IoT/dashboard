package controllers

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

//InitializeInfluxDb and set up write api
func InitializeInfluxDb() (api.QueryAPI, influxdb2.Client) {
	clientInflux := influxdb2.NewClientWithOptions("http://influxdb:8086", "telegraf:sensors",
		influxdb2.DefaultOptions().SetBatchSize(10))
	queryAPI := clientInflux.QueryAPI("")
	return queryAPI, clientInflux
}

//get past 24 hours of records
func GetRecords(queryAPI api.QueryAPI) {
	result, err := queryAPI.Query(context.Background(), `from(bucket:"sensors")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "humidityTemps")`)
	if err == nil {
		for result.Next() {
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			fmt.Printf("row: %s\n", result.Record().String())
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	} else {
		fmt.Printf("Query error: %s\n", err.Error())
	}
}

//CloseInfluxClient close client
func CloseInfluxClient(w api.WriteAPI, clientInflux influxdb2.Client) {
	w.Flush()
	clientInflux.Close()
}
