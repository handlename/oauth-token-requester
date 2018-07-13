package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/mrjones/oauth"
)

var serviceProviders = map[string]*oauth.ServiceProvider{
	"hatena": &oauth.ServiceProvider{
		RequestTokenUrl:   "https://www.hatena.com/oauth/initiate",
		AuthorizeTokenUrl: "https://www.hatena.ne.jp/oauth/authorize",
		AccessTokenUrl:    "https://www.hatena.com/oauth/token",
	},
	"twitter": &oauth.ServiceProvider{
		RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	},
}

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "enable debug mode")
	flag.Parse()

	service := read("Service Provider")
	serviceProvider := serviceProviders[service]
	if serviceProvider == nil {
		exitWithError(fmt.Errorf("unknown service provider: " + service))
	}

	consumerKey := read("Consumer Key")
	consumerSecret := read("Consumer Secret")

	consumer := oauth.NewConsumer(consumerKey, consumerSecret, *serviceProvider)
	consumer.Debug(debug)

	requestToken, loginURL, err := consumer.GetRequestTokenAndUrl("oob")
	if err != nil {
		exitWithError(err)
	}

	fmt.Printf("Open %s and get code\n", loginURL)
	fmt.Printf("Code: ")

	code := read("Code")

	token, err := consumer.AuthorizeToken(requestToken, code)
	if err != nil {
		exitWithError(err)
	}

	fmt.Println("Token =", token.Token)
	fmt.Println("Secret =", token.Secret)
}

func read(message string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(message, ": ")
	scanner.Scan()
	return scanner.Text()
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
