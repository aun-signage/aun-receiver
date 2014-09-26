# aun-receiver

`aun-receiver` receives tweets and IRC messages and sends them to MQTT.
It is intended to use in combination with `aun-subscreen`.

## Deploy to Heroku

### Create an application

    $ heroku create --buildpack https://github.com/kr/heroku-buildpack-go.git

### Add CloudMQTT

    $ heroku addons:add cloudmqtt:test

### Set environment variables

    $ git push heroku master

### Set environment variables

See [below](#configuration).

### Activate worker

    $ heroku ps:scale worker=1

## Configuration

If you want to receive twitter stream, you need to set these variables.

### MQTT\_URL

URL of destination MQTT. Note that you can not specify `mqtt` here; you need use `tcp` instead.

If you are using CloudMQTT, you can set MQTT\_URL automatically as follows:

    $ heroku config:set MQTT_URL=`heroku config:get CLOUDMQTT_URL | sed -e "s/^mqtt/tcp/"`

### TWITTER\_AUTH

Credentials for twitter. Tokens should be concatenated with ':'.

```
TWITTER_AUTH=[Consumer key]:[Consumer secret]:[Access token]:[Access token secret]
```

### TWITTER\_QUERY

Strings to track.

Example:

```
TWITTER_QUERY=aun,rubykaigi
```

* `TWITTER_QUERY` should be comma separated

## IRC Receiver Configurations

If you want to receive messages from IRC, you need to set these variables.

### IRC\_SERVER

Hostname of the IRC server.

Example:

```
IRC_SERVER=irc.example.com
```


### IRC\_PORT

Port number of the IRC server.

Example:

```
IRC_PORT=6667
```

### IRC\_NICK

Nick to use to connect the IRC server.

Example:

```
IRC_NICK=aun-receiver
```

### IRC\_CHANNELS

Channels to join.

EXAMPLE:

```
IRC_CHANNELS=#test1,#test2
```

* `IRC_CHANNELS` should be comma separated
