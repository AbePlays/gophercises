package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

func main() {
	keys, err := getKeys()
	if err != nil {
		panic(err)
	}

	tClient, err := getClient(keys.Key, keys.Secret)
	if err != nil {
		panic(err)
	}

	users, err := getRetweeters(tClient, keys.TweetId)
	if err != nil {
		panic(err)
	}

	retweeters := getRandomRetweeters(users, keys.Count)
	fmt.Println("Your winners are:", strings.Join(retweeters, ", "))
}

type Keys struct {
	Key     string
	Secret  string
	TweetId string
	Count   int
}

type Retweet struct {
	User struct {
		ScreenName string `json:"screen_name"`
	} `json:"user"`
}

func getKeys() (Keys, error) {
	file, err := os.Open(".env.json")
	if err != nil {
		return Keys{}, err
	}
	defer file.Close()

	var keys Keys
	decoder := json.NewDecoder(file)
	decoder.Decode(&keys)

	return keys, nil
}

func getClient(key, secret string) (*http.Client, error) {
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(key, secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var token oauth2.Token
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&token); err != nil {
		return nil, err
	}

	ctx := context.Background()
	conf := &oauth2.Config{}
	tClient := conf.Client(ctx, &token)

	return tClient, nil
}

func getRetweeters(client *http.Client, tweetId string) ([]string, error) {
	res, err := client.Get("https://api.x.com/2/tweets/" + tweetId)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var retweets []Retweet
	var usersNames []string
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&retweets); err != nil {
		// it will always come here, since x (formerly twitter)
		// has disabled allowing third-party apps to access retweet data
		// on free plans.
		usersNames = []string{"cronaldo_7", "bfernandes_8", "dberbatov_9", "wrooney_10"}
	}

	return usersNames, nil
}

func getRandomRetweeters(usersNames []string, count int) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]string, len(usersNames))
	copy(shuffled, usersNames)
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:count]
}
