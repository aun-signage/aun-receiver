package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/darashi/aun-receiver/twitter"
)

func clientId() string {
	pid := os.Getpid()
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Sprintf("%d", pid)
	}
	return fmt.Sprintf("%s.%d", hostname, pid)
}

func mqttClient(mqttUrl string) (*MQTT.MqttClient, error) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(mqttUrl)
	opts.SetCleanSession(true)
	opts.SetClientId(clientId())

	opts.SetOnConnectionLost(func(client *MQTT.MqttClient, reason error) {
		log.Fatal("MQTT CONNECTION LOST", reason) // TODO reconnect
	})

	parsed, err := url.Parse(mqttUrl)
	if err != nil {
		return nil, err
	}
	if user := parsed.User; user != nil {
		if username := user.Username(); username != "" {
			opts.SetUsername(username)
		}
		if password, set := user.Password(); set {
			opts.SetPassword(password)
		}
	}

	client := MQTT.NewClient(opts)
	_, err = client.Start()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func main() {
	mqttUrl := os.Getenv("MQTT_URL")
	if mqttUrl == "" {
		log.Fatal("You must specify MQTT_URL environment variable")
	}
	client, err := mqttClient(mqttUrl)
	if err != nil {
		log.Fatal(err)
	}

	twitterAuth := os.Getenv("TWITTER_AUTH")
	if twitterAuth == "" {
		log.Fatal("You must specify TWITTER_AUTH environment variable")
	}
	values := strings.SplitN(twitterAuth, ":", 4)
	if len(values) != 4 {
		log.Fatal("You must specify TWITTER_AUTH in [Consumer key]:[Consumer secret]:[Access token]:[Access token secret] format")
	}

	twitterQuery := os.Getenv("TWITTER_QUERY")
	if twitterQuery == "" {
		log.Fatal("You must specify TWITTER_QUERY environment variable")
	}
	log.Println("Tracking", twitterQuery)

	ch := make(chan string)

	go func() {
		for buf := range ch {
			log.Println(buf)
			client.Publish(MQTT.QOS_ZERO, "tweet", buf)
		}
	}()

	err = twitter.Track(
		values[0],
		values[1],
		values[2],
		values[3],
		twitterQuery,
		ch,
	)
	if err != nil {
		log.Fatal(err)
	}
}
