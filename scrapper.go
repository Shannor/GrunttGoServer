package main

import (
	"net/http"
	"log"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// Base String for url
const baseURL = "http://www.readcomics.tv/"

//Move Struct to specific Functions
type Comic struct {
	Title string
	Link string
	Category string
}

func allComicScrape(w http.ResponseWriter, r *http.Request) {

	doc, err := goquery.NewDocument(baseURL + "comic-list")
	if err != nil{
		//Return error here
		log.Fatal(err)
	}

	var comicList []Comic


	doc.Find(".container li").Each(func(index int, item *goquery.Selection){
		comic := Comic{}
		comic.Title = item.Children().Text()
		comic.Link, _ = item.Children().Attr("href")
		comic.Category = item.Parent().SiblingsFiltered("div").Text()

		if _, err := strconv.Atoi(comic.Category); err == nil{
			//Category is a number
			comic.Category = "#"
		}

		comicList = append(comicList, comic)
	})

	res, _ := json.Marshal(comicList)
	fmt.Fprintf(w, "Last Comic: %s", string(res))
}



// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }

func main() {
	http.HandleFunc("/", allComicScrape)
	http.ListenAndServe(":8000",nil)
}

