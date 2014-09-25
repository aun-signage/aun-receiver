# aun-social-receiver

## Deployment

    $ heroku create --buildpack https://github.com/kr/heroku-buildpack-go.git
    $ git push heroku master

    $ heroku config:set TWITTER_AUTH=[Consumer key]:[Consumer secret]:[Access token]:[Access token secret]
    $ heroku config:set TWITTER_QUERY=[Query]

    $ heroku config:set IRC_SERVER=[Server] IRC_PORT=[Port] IRC_NICK=[Nick] IRC_CHANNELS=[Channels]

    $ heroku addons:add cloudmqtt:test
    $ heroku config:set MQTT_URL=`heroku config:get CLOUDMQTT_URL | sed -e "s/^mqtt/tcp/"`
    $ heroku ps:scale worker=1

