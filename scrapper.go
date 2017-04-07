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


func allComicScrape(w http.ResponseWriter, r *http.Request) {

	doc, err := goquery.NewDocument(baseURL + "comic-list")
	if err != nil{
		//Return error here
		log.Fatal(err)
	}

	type Comic struct {
		Title string
		Link string
		Category string
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
	fmt.Fprintf(w, string(res))
}

func getPopularComics(w http.ResponseWriter, r *http.Request) {

	pageNumber := r.URL.Path[len("/popular-comics/"):]
	doc, err := goquery.NewDocument(baseURL + "popular-comic/" + pageNumber)
	if err != nil{
		log.Fatal(err)
	}

	type Genre struct {
		Name string
		GenreLink string
	}

	type PopComic struct{
		Title string
		Link string
		Img string
		Genres []Genre
	}


	var popularcomics []PopComic
	

	doc.Find(".manga-box").Each(func(index int, item *goquery.Selection){
		comic := PopComic{}
		//Gets top level information
		comic.Title = item.Find("h3").Children().Text()
		comic.Link, _ = item.Find("h3").Children().Attr("href")
		comic.Img, _ = item.Find("img").Attr("src")

		item.Find(".tags").Each(func(index int, child *goquery.Selection){
			genre := Genre{}
			genre.Name = child.Text()
			genre.GenreLink, _ = child.Attr("href")
			comic.Genres = append(comic.Genres, genre)
		})

		popularcomics = append(popularcomics, comic)
	})
	res, _ := json.Marshal(popularcomics)
	fmt.Fprintf(w, string(res))
}

// func getChapters(w http.ResponseWriter, r *http.Request) {
// 	//Get the Comic name out the URL
// 	comicName := r.URL.Path[len("/listchapters/"):]
// 	doc, err := goquery.NewDocument(baseURL + "comic/" + comicName)

// 	var chapters := []struct{
// 		chapter
// 	}
// }


// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }

func main() {
	http.HandleFunc("/comic-list-AZ", allComicScrape)
	http.HandleFunc("/popular-comics/",getPopularComics)
	http.ListenAndServe(":8000",nil)
}

