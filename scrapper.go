package main

import (
	"net/http"
	"log"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

type Chapter struct{
	ChapterName string
	Link string
	ReleaseDate string
}

func getChapters(w http.ResponseWriter, r *http.Request) {
	//Get the Comic name out the URL
	comicName := r.URL.Path[len("/chapter-list/"):]
	doc, err := goquery.NewDocument(baseURL + "comic/" + comicName)

	if err != nil {
		log.Fatal(err)
	}

	var chapters []Chapter

	doc.Find(".basic-list").Children().Each(func(index int, item *goquery.Selection){
		chapter := Chapter{}
		chapter.Link, _ = item.ChildrenFiltered("a").Attr("href")
		chapter.ChapterName = item.ChildrenFiltered("a").Text()
		chapter.ReleaseDate = item.ChildrenFiltered("Span").Text()
		chapters = append(chapters, chapter)

	})

	pageCount :=  doc.Find(".general-nav").Children().Length() - 1

	if pageCount > 0{
		chapterChannels := make(chan []Chapter, pageCount)

		go getExtraChapters(pageCount, comicName, chapterChannels)

		for i := range chapterChannels {
			chapters = append(chapters, i...)
		}
	}

	res, _ := json.Marshal(chapters)
	fmt.Fprintf(w, string(res))
}


func getExtraChapters(extras int, comicName string, cc chan []Chapter){
	for i := 2; i <= extras; i++ {

		doc, err := goquery.NewDocument(baseURL + "comic/" + comicName + "/" + strconv.Itoa(i))
		if err != nil{
			log.Fatal(err)
		}

		var chapters []Chapter
		doc.Find(".basic-list").Children().Each(func(index int, item *goquery.Selection){
			chapter := Chapter{}
			chapter.Link, _ = item.ChildrenFiltered("a").Attr("href")
			chapter.ChapterName = item.ChildrenFiltered("a").Text()
			chapter.ReleaseDate = item.ChildrenFiltered("Span").Text()
			chapters = append(chapters, chapter)
		})

		cc <- chapters
	}
	close(cc)
}

func readComic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/read-comic/"):]
	//0 = ComicName, 1 = Chapter Number
	paths := strings.Split(path, "/")
	if len(paths) < 2{
		//TODO:Change to network response
		log.Fatal("Missing Comic Name or Chapter Number.")
	}
	url := baseURL + paths[0] + "/chapter-" + paths[1]
	doc, err := goquery.NewDocument(url)

	if err != nil{
		log.Fatal(err)
	}

	numOfPages := doc.Find(".full-select").Last().Children().Length()
	pagesChannels := make(chan string, numOfPages)


	go getComicImageURL(url, numOfPages,pagesChannels)

	var urls []string
	for i := range pagesChannels{
		urls = append(urls, i)
	}

	res, _ := json.Marshal(urls)
	fmt.Fprintf(w, string(res))
}

func getComicImageURL(url string, numOfPages int, c chan string ){
	for i := 1; i <= numOfPages; i++{
		pageUrl := url + "/" + strconv.Itoa(i)
		doc, err := goquery.NewDocument(pageUrl)
		if err !=nil{
			log.Fatal(err)
		}
		link, _ := doc.Find("#main_img").Attr("src")
		c <- link
	}
	close(c)
}

func main() {
	http.HandleFunc("/comic-list-AZ", allComicScrape)
	http.HandleFunc("/popular-comics/",getPopularComics)
	http.HandleFunc("/chapter-list/", getChapters)
	http.HandleFunc("/read-comic/", readComic)
	http.ListenAndServe(":8000",nil)
}

