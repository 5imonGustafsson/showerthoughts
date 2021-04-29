package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type config struct {
	mqttPort string
	mqttHost string
	QOS      int
}

func main() {
	config := loadConfig()
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%s", config.mqttHost, config.mqttPort))

	opts.ConnectTimeout = time.Second // Minimal delays on connect
	opts.WriteTimeout = time.Second   // Minimal delays on writes
	opts.KeepAlive = 10               // Keepalive every 10 seconds so we quickly detect network outages
	opts.PingTimeout = time.Second    // local broker so response should be quick

	// Automate connection management (will keep trying to connect and will reconnect if network drops)
	opts.ConnectRetry = true
	opts.AutoReconnect = true

	// Log events
	opts.OnConnectionLost = func(cl mqtt.Client, err error) {
		fmt.Println("connection lost")
	}

	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("connection established")

		// Establish the subscription - doing this here means that it will happen every time a connection is established
		// (useful if opts.CleanSession is TRUE or the broker does not reliably store session data)
		t := c.Subscribe("*", byte(config.QOS), func(_ mqtt.Client, msg mqtt.Message) {
			log.Printf("%+v", msg)
		})
		// the connection handler is called in a goroutine so blocking here would hot cause an issue. However as blocking
		// in other handlers does cause problems its best to just assume we should not block
		go func() {
			_ = t.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
			if t.Error() != nil {
				fmt.Printf("ERROR SUBSCRIBING: %s\n", t.Error())
			} else {
				fmt.Println("subscribed to: *")
			}
		}()
	}
}

func loadConfig() config {
	return config{
		mqttPort: getStrEnv(os.Getenv("MQTT_PORT")),
		mqttHost: getStrEnv(os.Getenv("MQTT_HOST")),
		QOS:      getIntEnv(os.Getenv("MQTT_QOS")),
	}
}

func getStrEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("some error msg"))
	}
	return val
}

func getIntEnv(key string) int {
	val := getStrEnv(key)
	ret, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Sprintf("some error"))
	}
	return ret
}

func getBoolEnv(key string) bool {
	val := getStrEnv(key)
	ret, err := strconv.ParseBool(val)
	if err != nil {
		panic(fmt.Sprintf("some error"))
	}
	return ret
}
