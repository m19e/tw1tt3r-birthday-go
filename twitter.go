package main

import (
	"github.com/joho/godotenv"
	//_ "github.com/joho/godotenv/autoload"
	"github.com/ChimeraCoder/anaconda"
	"github.com/PuerkitoBio/goquery"

	"encoding/json"
	f "fmt"
	"log"
	"net/url"
	. "os"
	"strings"
)

type BD struct {
	Name string
	Date string
}

type BDList []BD

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed loading .env file")
	}
}

func GetTwitterApi() *anaconda.TwitterApi {
	return anaconda.NewTwitterApiWithCredentials(
		Getenv("ACCESS_TOKEN_KEY"),
		Getenv("ACCESS_TOKEN_SECRET"),
		Getenv("CONSUMER_KEY"),
		Getenv("CONSUMER_SECRET"),
	)
}

func GetFriendsIdList(ids chan anaconda.FriendsIdsPage) []int64 {
	friends := make([]int64, 0)
	friendsChan := ids

friendsLoop:
	for {
		select {
		case p, ok := <-friendsChan:
			if ok {
				for _, id := range p.Ids {
					friends = append(friends, id)
				}
			} else {
				break friendsLoop
			}
		}
	}
	return friends
}

func GetBirthday(sn string, target string) BDList {
	var bdlist BDList
	doc, _ := goquery.NewDocument("https://twitter.com/" + sn)
	doc.Find("span").Each(func(_ int, s *goquery.Selection) {
		elem := s.Text()
		if strings.Contains(elem, target) {
			e := strings.TrimSpace(elem)
			bdlist = append(bdlist, BD{Name: sn, Date: e})
		}
	})
	return bdlist
}

func Chunks(l []int64, n int) chan []int64 {
	ch := make(chan []int64)

	go func() {
		for i := 0; i < len(l); i += n {
			from := i
			to := i + n
			if to > len(l) {
				to = len(l)
			}
			ch <- l[from:to]
		}
		close(ch)
	}()
	return ch
}

func main() {
	loadEnv()

	api := GetTwitterApi()

	v := url.Values{}
	v.Set("count", "5000")
	v.Set("screen_name", "imassc_official")

	//q := "ゆゆ式 -RT"

	// result, _ := api.GetSearch(q, v)
	//	result, _ := api.GetUserTimeline(v)
	//	for _, tweet := range result {
	//		f.Println("----------")
	//		f.Println(tweet.FullText)
	//	}

	bd := GetBirthday("TwitterJP", "誕生日")
	b, _ := json.Marshal(bd)
	f.Printf("%s\n", b)

	ids := api.GetFriendsIdsAll(v)

	friends := GetFriendsIdList(ids)

	f.Printf("length: %d\n", len(friends))
	for _, n := range friends {
		f.Println(n)
	}

	size := 100
	for l := range Chunks(friends, size) {
		f.Printf("length: %d\n", len(l))
		f.Println(l)
	}
}
