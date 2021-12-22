package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	controllers "subs/controllers"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

var writeAPI api.WriteAPI
var influxClient influxdb2.Client

var knt int
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("sensors: %s\n", msg.Payload())
	dataPoint := controllers.CreateTCANPoint(string(msg.Payload()))
	controllers.PushToDb(dataPoint, writeAPI)
	knt++
}
var knt2 int
var f2 MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("health: %s\n", msg.Payload())
	dataPoint := controllers.CreateHealthPoint(string(msg.Payload()))
	controllers.PushToDb(dataPoint, writeAPI)
	knt2++
}

func main() {

	writeAPI, influxClient = controllers.InitializeInfluxDb()
	defer controllers.CloseInfluxClient(writeAPI, influxClient)

	knt = 0
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	opts := MQTT.NewClientOptions().AddBroker("tcp://mosquitto:1883/")
	opts.SetClientID("influx-pusher")
	opts.SetDefaultPublishHandler(f)
	topic := "sensors"
	topic2 := "health"

	opts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		if token2 := c.Subscribe(topic2, 0, f2); token2.Wait() && token2.Error() != nil {
			panic(token2.Error())
		}
	}
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to server\n")
	}
	<-c
}
