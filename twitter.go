package main // import "tw"

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
	"time"
)

type User struct {
	SN   string
	Name string
}

type BD struct {
	User  User
	Date  string
	Exist bool
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

func GetBirthday(sn string, target string) BD {
	var bd = BD{Exist: false}
	doc, _ := goquery.NewDocument("https://twitter.com/" + sn)
	doc.Find("span").Each(func(_ int, s *goquery.Selection) {
		elem := s.Text()
		if strings.Contains(elem, target) {
			e := strings.TrimSpace(elem)
			bd = BD{User: User{SN: sn}, Date: e, Exist: true}
		}
	})
	return bd
}

func GetBirthdayAll(users []anaconda.User) chan BD {
	// var bdList []BD
	var wg sync.WaitGroup
	ch := make(chan BD, 200)

	for _, user := range users {
		wg.Add(1)
		go func(name string) {
			u := GetBirthday(name, "誕生日")
			if u.Exist {
				ch <- u
			}
			wg.Done()
		}(user.ScreenName)
	}
	wg.Wait()
	close(ch)

	// return bdList
	return ch
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

func Flatten(double [][]anaconda.User) []anaconda.User {
	var flat []anaconda.User

	for _, single := range double {
		for _, bd := range single {
			flat = append(flat, bd)
		}
	}

	return flat
}

func main() {
	start := time.Now()
	log.Print("start")

	loadEnv()

	api := GetTwitterApi()

	v := url.Values{}
	v.Set("count", "5000")
	v.Set("screen_name", "imas_official")

	//q := "ゆゆ式 -RT"

	// result, _ := api.GetSearch(q, v)
	//	result, _ := api.GetUserTimeline(v)
	//	for _, tweet := range result {
	//		f.Println("----------")
	//		f.Println(tweet.FullText)
	//	}

	// bd := GetBirthday("TwitterJP", "誕生日")
	// b, _ := json.Marshal(bd)
	// f.Printf("%s\n", b)

	idsAll := api.GetFriendsIdsAll(v)
	friends := GetFriendsIdList(idsAll)

	f.Printf("%d friends\n", len(friends))

	size := 100
	chunked := Chunks(friends, size)

	var users [][]anaconda.User
	var wg sync.WaitGroup
	p := url.Values{}
	p.Set("include_entities", "false")

	for _, ids := range chunked {
		wg.Add(1)
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

	// f.Printf("\n%d\n", len(users))
	semi := Flatten(users)

	// for _, user := range semi {
	// 	bd := GetBirthday(user.ScreenName, "誕生日")
	// 	f.Printf("%s:%s\n", user.Name, bd)
	// }

	final := GetBirthdayAll(semi)

	f.Printf("%d friends have Birthday\n", len(final))

	for {
		bd, ok := <-final
		if !ok {
			break
		}
		w, _ := json.Marshal(bd)
		f.Printf("%s\n", w)
	}

	end := time.Now()

	log.Printf("all processes finish in %f", (end.Sub(start).Seconds()))
}
