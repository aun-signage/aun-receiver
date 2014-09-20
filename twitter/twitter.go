package twitter

import (
	"bufio"
	"log"

	"github.com/mrjones/oauth"
)

func Track(
	consumerKey, consumerSecret, accessToken, accessTokenSecret,
	track string,
	ch chan<- string,
) error {
	c := oauth.NewConsumer(
		consumerKey,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		})

	params := map[string]string{"track": track}
	response, err := c.Post(
		"https://stream.twitter.com/1.1/statuses/filter.json",
		params,
		&oauth.AccessToken{
			Token:  accessToken,
			Secret: accessTokenSecret})
	if err != nil {
		return err
	}
	defer response.Body.Close()
	log.Println("Stream connected")

	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		text := scanner.Text()
		if text != "" {
			ch <- text
		}
	}

	return scanner.Err()
}
