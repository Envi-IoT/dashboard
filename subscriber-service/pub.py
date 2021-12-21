import paho.mqtt.client as paho
import json
import time
import random
broker="localhost"
port=1883
def on_publish(client,userdata,result):
    print("data published \n")
    pass
client1= paho.Client("control1")
client1.on_publish = on_publish
client1.connect(broker,port)       
for x in range(10):
    meas = {"messageID": x, "sensorID": "bin1", "temperature": x, "timestamp": time.time()}
    ret= client1.publish("sensors",json.dumps(meas))
    time.sleep(10)