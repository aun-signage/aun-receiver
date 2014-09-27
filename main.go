package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"

	"github.com/darashi/aun-receiver/irc"
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

type MqttMessage struct {
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

	chTweet := make(chan string)
	twitterAuth := os.Getenv("TWITTER_AUTH")
	twitterQuery := os.Getenv("TWITTER_QUERY")
	if twitterAuth != "" && twitterQuery != "" {
		twitterAuthComponents := strings.SplitN(twitterAuth, ":", 4)
		if len(twitterAuthComponents) != 4 {
			log.Fatal("You must specify TWITTER_AUTH in [Consumer key]:[Consumer secret]:[Access token]:[Access token secret] format")
		}
		log.Println("Tracking", twitterQuery)
		go func() {
			err = twitter.Track(
				twitterAuthComponents[0],
				twitterAuthComponents[1],
				twitterAuthComponents[2],
				twitterAuthComponents[3],
				twitterQuery,
				chTweet,
			)
			if err != nil {
				log.Fatal(err)
			}
		}()
	} else {
		log.Println("Twitter receiver not configured. You need to specify TWITTER_AUTH and TWITTER_QUERY.")
	}

	chIrc := make(chan string)
	ircServer := os.Getenv("IRC_SERVER")
	ircPort := os.Getenv("IRC_PORT")
	ircNick := os.Getenv("IRC_NICK")
	ircChannels := os.Getenv("IRC_CHANNELS")
	if ircServer != "" && ircPort != "" && ircNick != "" && ircChannels != "" {
		port, err := strconv.Atoi(ircPort)
		if err != nil {
			log.Fatal(err)
		}
		channels := strings.Split(ircChannels, ",")
		irc.Receive(
			ircServer,
			port,
			ircNick,
			channels,
			chIrc,
		)
	} else {
		log.Println("IRC receiver not configured. You need to specify IRC_SERVER, IRC_PORT, IRC_NICK and IRC_CHANNELS.")
	}

	for {
		select {
		case buf := <-chTweet:
			log.Println(buf)
			client.Publish(MQTT.QOS_ZERO, "social-stream/tweet", buf)
		case buf := <-chIrc:
			log.Println(buf)
			client.Publish(MQTT.QOS_ZERO, "social-stream/irc", buf)
		}
	}
}
