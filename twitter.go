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
	"sync"
)

type User struct {
	SN   string
	Name string
}

type BD struct {
	USER User
	Date string
}

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

// func GetAllUsersLookup(f func(), idsList [][]int64) [][]anaconda.User {
// 	var users [][]anaconda.User
// 	var wg sync.WaitGroup
// 	v := url.Values{}
// 	v.Set("include_entities", "false")

// 	for _, ids := range idsList {
// 		wg.Add(1)
// 		go func(ids []int64) {
// 			defer wg.Done()
// 			res, _ := api.GetUsersLookupByIds(ids, v)
// 			users = append(users, res)
// 		}(ids)
// 	}

// 	wg.Wait()
// 	return users
// }

func GetBirthday(sn string, target string) []BD {
	var bdlist []BD
	doc, _ := goquery.NewDocument("https://twitter.com/" + sn)
	doc.Find("span").Each(func(_ int, s *goquery.Selection) {
		elem := s.Text()
		if strings.Contains(elem, target) {
			e := strings.TrimSpace(elem)
			bdlist = append(bdlist, BD{USER: User{SN: sn}, Date: e})
		}
	})
	return bdlist
}

func Chunks(l []int64, n int) [][]int64 {
	var result [][]int64

	for i := 0; i < len(l); i += n {
		from := i
		to := i + n
		if to > len(l) {
			to = len(l)
		}
		result = append(result, l[from:to])
	}

	return result
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

	idsAll := api.GetFriendsIdsAll(v)

	friends := GetFriendsIdList(idsAll)

	f.Printf("length: %d\n", len(friends))
	for _, n := range friends {
		f.Println(n)
	}

	size := 100
	chunked := Chunks(friends, size)

	var users [][]anaconda.User
	var wg sync.WaitGroup
	p := url.Values{}
	p.Set("include_entities", "false")

	for _, ids := range chunked {
		wg.Add(1)
		log.Print(ids)
		go func(ids []int64) {
			defer wg.Done()
			res, err := api.GetUsersLookupByIds(ids, p)
			if err != nil {
				log.Fatal(err)
			}
			users = append(users, res)
		}(ids)
	}

	wg.Wait()

	// u, _ := json.Marshal(users)
	// f.Printf("%s\n", u)

	// f.Printf("%s\n", users[0][0].Name)
	for _, mas := range users {
		for _, user := range mas {
			bd := GetBirthday(user.ScreenName, "誕生日")
			if len(bd) != 0 {
				f.Printf("%s:%s\n", user.Name, bd[0])
			}
		}
	}
}
