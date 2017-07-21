package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"scrapper/utils"

	"github.com/PuerkitoBio/goquery"
	"google.golang.org/appengine"
)

// COMICEXTRA url
const COMICEXTRA = "http://www.comicextra.com/"

// READCOMIC url
const READCOMIC = "http://readcomics.website/"

const comicExtraURLParam = "ce"
const readcomicsURLParam = "rcw"

//Request format -> /chapter-list/{chpater Name}?url={website}
func allComicsRequest(w http.ResponseWriter, r *http.Request) {

	param := r.URL.Query().Get("url")
	url, formatErr := utils.CreateAllComicsURL(param)

	if formatErr != nil {
		utils.ErrorHandler(w, http.StatusBadRequest, formatErr)
		return
	}

	doc, err := utils.GetGoQueryDoc(url, r)
	if err != nil {
		utils.ErrorHandler(w, http.StatusBadRequest, err)
		return
	}

	comicList := utils.GetAllComics(doc, param)

	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(comicList)
	w.Write(res)
}

//Request Format -> /popular-comics?page={pageNumber}&url={website}
func getPopularComics(w http.ResponseWriter, r *http.Request) {

	pageNumber := r.URL.Query().Get("page")
	if _, err := strconv.Atoi(pageNumber); err != nil {
		//TODO: Throw wrong formated request type error
		return
	}

	choice := r.URL.Query().Get("url")
	url, err := utils.CreatePopularComicsURL(choice, pageNumber)
	if err != nil {
		return
	}

	doc, err := utils.GetGoQueryDoc(url, r)
	if err != nil {
		return
	}

	popularComics := utils.GetPopularComics(doc, choice)

	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(popularComics)
	w.Write(res)
}

//Chapter : Struct to represent chapters

//Request Format -> /chapter-list/{Comic Name}?url={type}
func getChapters(w http.ResponseWriter, r *http.Request) {

	comicName := r.URL.Path[len("/chapter-list/"):]

	choice := r.URL.Query().Get("url")
	url, err := utils.CreateChapterURL(choice, comicName)
	if err != nil {
		return
	}

	doc, err := utils.GetGoQueryDoc(url, r)
	if err != nil {
		utils.ErrorHandler(w, http.StatusBadRequest, err)
	}

	chapters := utils.GetChapters(doc, r, choice, comicName)

	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(chapters)
	w.Write(res)
}

//Request Format -> /read-comic/{Comic Name}/{Chapter Number}?url={param}
func readComic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/read-comic/"):]
	//0 = ComicName, 1 = Chapter Number
	paths := strings.Split(path, "/")
	comicName, chapterNumber := paths[0], paths[1]
	//TODO: check how split returns resutls
	if len(paths) != 2 {
		// errorHandler(w, r, http.StatusBadRequest,err)
		return
	}

	choice := r.URL.Query().Get("url")
	url, err := utils.CreateReadComicURL(choice, comicName, chapterNumber)
	if err != nil {
		utils.ErrorHandler(w, http.StatusBadRequest, err)
		return
	}

	doc, err := utils.GetGoQueryDoc(url, r)
	if err != nil {
		utils.ErrorHandler(w, http.StatusBadRequest, err)
		return
	}
	urls := util.GetChapterImages(doc, r, choice, url)
	var urls []string

	switch choice {
	case comicExtraURLParam:
		numOfPages := doc.Find(".full-select").First().Children().Length()
		pagesChannels := make(chan string, numOfPages)

		go getComicImageURL(url, r, numOfPages, pagesChannels)

		for i := range pagesChannels {
			urls = append(urls, i)
		}
	case readcomicsURLParam:

	default:
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(urls)
	w.Write(res)
}

//Helper function for ReadComic function
func getComicImageURL(url string, r *http.Request, numOfPages int, cc chan string) {
	for i := 1; i <= numOfPages; i++ {
		pageURL := url + "/" + strconv.Itoa(i)

		doc, err := utils.GetGoQueryDoc(pageURL, r)
		if err != nil {
			return
		}
		link, _ := doc.Find("#main_img").Attr("src")
		cc <- link
	}
	close(cc)
}

//Request Format -> /search-categories
func getSearchCategories(w http.ResponseWriter, r *http.Request) {
	url := COMICEXTRA + "advanced-search"
	var (
		categories []string
	)

	doc, err := utils.GetGoQueryDoc(url, r)
	if err != nil {
		utils.ErrorHandler(w, http.StatusBadRequest, err)
	}

	doc.Find(".search-checks").ChildrenFiltered("li").Each(func(index int, item *goquery.Selection) {
		categories = append(categories, item.Text())
	})

	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(categories)
	fmt.Fprintf(w, string(res))

}

//Request Format -> /advanced-search?key={Search Param}&include={Desired Attributes}
//&exclude={Undesired attributes}&status={status}&page={Page number for results}
func performAdvancedSearch(w http.ResponseWriter, r *http.Request) {
	url := COMICEXTRA + "advanced-search?"
	queryParams := r.URL.Query()

	if queryParams.Get("key") != "" {
		url += "key=" + queryParams.Get("key")
	}

	if queryParams.Get("include") != "" {
		url += "&wg=" + queryParams.Get("include")
	}

	if queryParams.Get("exclude") != "" {
		url += "&wog=" + queryParams.Get("exclude")
	}

	if queryParams.Get("status") != "" {
		url += "&status=" + queryParams.Get("status")
	}

	if queryParams.Get("page") != "" {
		url += "&page=" + queryParams.Get("page")
	}

	type SearchResult struct {
		Title string `json:"title"`
		Link  string `json:"link"`
		Img   string `json:"img"`
		Genre string `json:"genre"`
	}

	var results []SearchResult

	doc, err := utils.GetGoQueryDoc(url, r)
	if err != nil {
		utils.ErrorHandler(w, http.StatusBadRequest, err)
	}

	doc.Find(".manga-box").Each(func(index int, item *goquery.Selection) {
		comic := SearchResult{}
		comic.Title = item.Find("h3").Children().Text()
		comic.Link, _ = item.Find("h3").Children().Attr("href")
		comic.Img, _ = item.Find("img").Attr("src")
		comic.Genre = strings.TrimSpace(item.Find(".genre").Text())
		results = append(results, comic)
	})

	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(results)
	fmt.Fprintf(w, string(res))
}

//Reqest Format -> /description?name={Comic Name}
func getDescription(w http.ResponseWriter, r *http.Request) {

	url := COMICEXTRA + "comic/"
	queryParams := r.URL.Query()

	if queryParams.Get("name") != "" {
		url += queryParams.Get("name")
	} else {
		return
	}

	doc, err := utils.GetGoQueryDoc(url, r)
	if err != nil {
		return
	}
	type Response struct {
		Description   string `json:"description"`
		LargeImg      string `json:"largeImg"`
		FormatedName  string `json:"formatedName"`
		Name          string `json:"name"`
		AlternateName string `json:"alternateName"`
		ReleaseYear   string `json:"releaseYear"`
		Status        string `json:"status"`
		Author        string `json:"author"`
		Genre         string `json:"genre"`
	}

	description := Response{FormatedName: queryParams.Get("name")}

	doc.Find(".manga-details").Find("tr").Each(func(index int, item *goquery.Selection) {
		//Chain of ugly if statetments
		pair := strings.Split(item.First().Children().Text(), ":")
		key := strings.TrimSpace(pair[0])
		val := strings.TrimSpace(pair[1])

		if key == "Name" {
			description.Name = val
		} else if key == "Alternate Name" {
			description.AlternateName = val
		} else if key == "Year of Release" {
			description.ReleaseYear = val
		} else if key == "Status" {
			description.Status = val
		} else if key == "Author" {
			description.Author = val
		} else if key == "Genre" {
			description.Genre = val
		}

	})
	description.LargeImg, _ = doc.Find(".manga-image").ChildrenFiltered("img").Attr("src")
	description.Description = strings.TrimSpace(doc.Find(".pdesc").Text())

	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(description)
	fmt.Fprintf(w, string(res))

}

//Add a query param for which url to use
func main() {

	http.HandleFunc("/comic-list-AZ", allComicsRequest)
	http.HandleFunc("/popular-comics/", getPopularComics)
	http.HandleFunc("/chapter-list/", getChapters)
	http.HandleFunc("/read-comic/", readComic)
	http.HandleFunc("/search-categories", getSearchCategories)
	http.HandleFunc("/advanced-search", performAdvancedSearch)
	http.HandleFunc("/description", getDescription)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello fellow, Gruntt User!")
		return
	})
	http.ListenAndServe(":8000", nil)
	appengine.Main()

}
