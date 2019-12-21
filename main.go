package main

import (
  "github.com/gorilla/mux"
  "github.com/PuerkitoBio/goquery"
  "fmt"
  "net/http"
  "encoding/json"
  "strings"
)

type Pagedata struct { //jsonの構造
  URL   []string
}

const ID = "wa_jwa_jpon2"

func GetPage(url string, target string) ([]string) {
  var array []string
  doc, _ := goquery.NewDocument(url)
  doc.Find("span").Each(func(_ int, s *goquery.Selection) {
    url := s.Text()
    if strings.Contains(url, target) {
		array = append(array,url)
    }
  })

  return array
}

func handlerBirthday(w http.ResponseWriter, r *http.Request) {
  url := "https://twitter.com/" + ID //任意のurl取れるように改造したい

  pagedata := GetPage(url, "誕生日")
  pages := Pagedata{pagedata}

  res, err := json.Marshal(pages)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  //*は危険なので個別指定にしておくのが良さそう fixme
  w.Header().Set("Access-Control-Allow-Origin", "*")

  w.Header().Set("Content-Type", "application/json")
  w.Write(res)
}

func handlerFollowing(w http.ResponseWriter, r *http.Request) {
  url := "https://twitter.com/"+ ID +"/following" //任意のurl取れるように改造したい

  pagedata := GetPage(url, "@")
  pages := Pagedata{pagedata}

  res, err := json.Marshal(pages)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  //*は危険なので個別指定にしておくのが良さそう fixme
  w.Header().Set("Access-Control-Allow-Origin", "*")

  w.Header().Set("Content-Type", "application/json")
  w.Write(res)
}

func BirthdayHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  sn := vars["id"]
  url := "https://twitter.com/" + sn //任意のurl取れるように改造したい

  fmt.Println(sn)

  pagedata := GetPage(url, "誕生日")
  pages := Pagedata{pagedata}

  res, err := json.Marshal(pages)
  fmt.Println(pages)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  //*は危険なので個別指定にしておくのが良さそう fixme
  w.Header().Set("Access-Control-Allow-Origin", "*")

  w.Header().Set("Content-Type", "application/json")
  w.Write(res)
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/bd/{id}", BirthdayHandler)
  //http.HandleFunc("/", handlerBirthday)
  //http.HandleFunc("/following", handlerFollowing)
  
  fmt.Printf("server is running\n8080port\n")
  http.ListenAndServe(":8080", r)   // サーバーを起動するよ!
}