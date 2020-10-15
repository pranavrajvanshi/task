package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Article struct {
	title              string    `json:"title"`
	subtitle           string    `json:"subtitle"`
	Id                 string    `json:"Id"`
	Creation_Timestamp time.Time `json:Timestamp`
}

var Articles []Article
var mu sync.Mutex
var wg sync.WaitGroup

func homePage(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
func handleRequests() {

	http.HandleFunc("/", homePage)

	http.HandleFunc("/articles", getorpost)

	// http.HandleFunc("/articles/search", searchQuery)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func getorpost(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		wg.Add(1)
		createNewArticle(w, r, &wg)
		wg.Wait()
	} else {

		returnAllArticles(w, r)

	}
}
func createNewArticle(w http.ResponseWriter, r *http.Request, wg *sync.WaitGroup) {
	mu.Lock()
	if r.Method == "POST" {
		reqBody, _ := ioutil.ReadAll(r.Body)
		var article Article
		json.Unmarshal(reqBody, &article)

		article.Id = Articles[len(Articles)-1].Id + "1"
		article.Creation_Timestamp = time.Now()
		Articles = append(Articles, article)

		json.NewEncoder(w).Encode(article)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	mu.Unlock()
	wg.Done()
}
func returnSingleArticle(w http.ResponseWriter, r *http.Request, p string) {

	key, err := strconv.Atoi(p)
	fmt.Println(err)
	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}
	}
	fmt.Println("Endpoint Hit: returnSingleArticle")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {

	value_js, err := json.Marshal(Articles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(value_js)
	w.Header().Set("Content-Type", "application/json")
	w.Write(value_js)

	fmt.Println("Endpoint Hit: returnAllArticles")

}

func main() {
	Articles = []Article{
		Article{Id: 1, Title: "AA", SubTitle: "Article Description", Content: "Article Content", Creation_Timestamp: time.Now()}}

	handleRequests()
}
