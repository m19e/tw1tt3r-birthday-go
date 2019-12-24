package main

import (
	"github.com/joho/godotenv"
	//_ "github.com/joho/godotenv/autoload"
	"github.com/ChimeraCoder/anaconda"
	
	"log"
	. "os"
	f "fmt"
	"net/url"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed loading .env file")
	}
}

func GetTwitterApi() *anaconda.TwitterApi {
	return anaconda.NewTwitterApiWithCredentials(Getenv("ACCESS_TOKEN_KEY"), Getenv("ACCESS_TOKEN_SECRET"), Getenv("CONSUMER_KEY"), Getenv("CONSUMER_SECRET"))
}

func main() {
	loadEnv()

	api := GetTwitterApi()

	v := url.Values{}
	v.Set("count", "50")
	
	result, _ := api.GetSearch("#ゆゆ式ac -RT", v)
	for _, tweet := range result.Statuses {
		f.Println(tweet.Text)
	}
}