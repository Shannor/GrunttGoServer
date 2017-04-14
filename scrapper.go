package main

import (
	"net/http"
	"log"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"github.com/PuerkitoBio/goquery"
)

// Base String for url
const baseURL = "http://www.readcomics.tv/"

func errorHandler(w http.ResponseWriter, r *http.Request, status int, err error) {
	http.Error(w, err.Error() , status)
}
//Request format -> /chapter-list/{chpater Name} 
func allComicsRequest(w http.ResponseWriter, r *http.Request) {

	var (
		doc *goquery.Document
	 	err, docErr error
		resp *http.Response
	)

	url := baseURL + "comic-list"

	if appengine.IsDevAppServer(){

		c := appengine.NewContext(r)
		client := urlfetch.Client(c)
		resp, err = client.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}else{

		resp, err = http.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}

	if err != nil{
		http.Error(w, err.Error() , http.StatusInternalServerError)
		return 
	}

	if resp.StatusCode != 200{
		http.Error(w, resp.Status, resp.StatusCode)
		return
	}

	if docErr != nil{

		http.Error(w, docErr.Error() , http.StatusInternalServerError)
		return 
	}

	defer resp.Body.Close()

	type Comic struct {
		Title string `json:"title"`
		Link string `json:"link"`
		Category string `json:"category"`
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

	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(comicList)
	fmt.Fprintf(w, string(res))
}

//Request Format -> /popular-comics/{page Number}
func getPopularComics(w http.ResponseWriter, r *http.Request) {

	pageNumber := r.URL.Path[len("/popular-comics/"):]
	//TODO: check if pageNumber is a number and not more than that
	var ( 
		doc *goquery.Document
		err, docErr error
		resp *http.Response
	)
	url := baseURL + "popular-comic/" + pageNumber

	if appengine.IsDevAppServer(){
		c := appengine.NewContext(r)
		client := urlfetch.Client(c)
		resp, err = client.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}else{
		resp, err = http.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}

	//Error with inital url request
	if err != nil{
		http.Error(w, err.Error() , http.StatusInternalServerError)
		return	
	}

	if resp.StatusCode != 200 {
		http.Error(w, resp.Status, resp.StatusCode)
		return
	}

	defer resp.Body.Close()

	//Error with goquery
	if docErr != nil{
		http.Error(w, err.Error() , http.StatusInternalServerError)
		return 
	}

	type Genre struct {
		Name string `json:"name"`
		GenreLink string `json:"genreLink"`
	}

	type PopComic struct{
		Title string `json:"title"`
		Link string `json:"link"`
		Img string 	`json:"img"`
		Genres []Genre `json:"genres"`
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
	
	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(popularcomics)
	fmt.Fprintf(w, string(res))
}

type Chapter struct{
	ChapterName string `json:"chapterName"`
	Link string `json:"link"`
	ReleaseDate string `json:"releaseDate"`
}
//Request Format -> /chapter-list/{Comic Name}
func getChapters(w http.ResponseWriter, r *http.Request) {

	comicName := r.URL.Path[len("/chapter-list/"):]
	// TODO: Check if name is right or wrong
	url := baseURL + "comic/" + comicName
	var ( 
		doc *goquery.Document
		err, docErr error
		resp *http.Response
	)

	if appengine.IsDevAppServer(){
		c := appengine.NewContext(r)
		client := urlfetch.Client(c)
		resp, err = client.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}else{
		resp, err = http.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}


	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err)
		return 
	}

	if resp.StatusCode != 200{
		http.Error(w, resp.Status, resp.StatusCode)
		return
	}

	if docErr != nil{
		errorHandler(w, r, http.StatusInternalServerError, err)
		return
	}

	defer resp.Body.Close()

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

		go getExtraChapters(pageCount,r ,comicName, chapterChannels)

		for i := range chapterChannels {
			chapters = append(chapters, i...)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(chapters)
	fmt.Fprintf(w, string(res))
}


func getExtraChapters(extras int, r *http.Request ,comicName string, cc chan []Chapter){
	for i := 2; i <= extras; i++ {

		var ( 
			doc *goquery.Document
			err, docErr error
			resp *http.Response
		)
		url := baseURL + "comic/" + comicName + "/" + strconv.Itoa(i)

		if appengine.IsDevAppServer(){
			c := appengine.NewContext(r)
			client := urlfetch.Client(c)
			resp, err = client.Get(url)
			doc, docErr = goquery.NewDocumentFromResponse(resp)
		}else{
			resp, err = http.Get(url)
			doc, docErr = goquery.NewDocumentFromResponse(resp)
		}

		if err != nil{
			log.Printf(err.Error())
			return 
		}

		if resp.StatusCode != 200{
			log.Printf(resp.Status)
			return
		}

		if docErr != nil{
			log.Printf(docErr.Error())
			return 
		}

		defer resp.Body.Close()

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

//Request Format -> /read-comic/{Comic Name}/{Chapter Number}
func readComic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/read-comic/"):]
	//0 = ComicName, 1 = Chapter Number
	paths := strings.Split(path, "/")
	//TODO: check how split returns resutls
	if len(paths) < 2{
		// errorHandler(w, r, http.StatusBadRequest,err)
		return 
	}
	url := baseURL + paths[0] + "/chapter-" + paths[1]
	var ( 
		doc *goquery.Document
		err, docErr error
		resp *http.Response
	)

	if appengine.IsDevAppServer(){
		c := appengine.NewContext(r)
		client := urlfetch.Client(c)
		resp, err = client.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}else{
		resp, err = http.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}

	if err != nil{
		errorHandler(w, r, http.StatusInternalServerError, err)
		return 
	}

	if resp.StatusCode != 200{
		http.Error(w, resp.Status, resp.StatusCode)
		return 
	}

	if docErr != nil{
		http.Error(w, docErr.Error(), http.StatusInternalServerError)
	}

	numOfPages := doc.Find(".full-select").Last().Children().Length()
	pagesChannels := make(chan string, numOfPages)

	defer resp.Body.Close()

	go getComicImageURL(url,r ,numOfPages,pagesChannels)

	var urls []string
	for i := range pagesChannels{
		urls = append(urls, i)
	}

	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(urls)
	fmt.Fprintf(w, string(res))
}
//Helper function for ReadComic function
func getComicImageURL(url string, r *http.Request ,numOfPages int, cc chan string ){
	for i := 1; i <= numOfPages; i++{
		pageUrl := url + "/" + strconv.Itoa(i)
		var ( 
			doc *goquery.Document
			err, docErr error
			resp *http.Response
		)

		if appengine.IsDevAppServer(){
			c := appengine.NewContext(r)
			client := urlfetch.Client(c)
			resp, err = client.Get(pageUrl)
			doc, docErr = goquery.NewDocumentFromResponse(resp)
		}else{
			resp, err = http.Get(pageUrl)
			doc, docErr = goquery.NewDocumentFromResponse(resp)
		}

		if err !=nil{
			log.Printf(err.Error())
			return 
		}
		if resp.StatusCode != 200{
			log.Printf(resp.Status)
			return 
		}

		if docErr != nil{
			log.Printf(docErr.Error())
			return
		}

		defer resp.Body.Close()

		link, _ := doc.Find("#main_img").Attr("src")
		cc <- link
	}
	close(cc)
}

//Request Format -> /search-categories
func getSearchCategories(w http.ResponseWriter, r *http.Request){
    url := baseURL + "advanced-search"
	var ( 
		categories []string
		doc *goquery.Document
		err, docErr error
		resp *http.Response
	)

	if appengine.IsDevAppServer(){
		c := appengine.NewContext(r)
		client := urlfetch.Client(c)
		resp, err = client.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}else{
		resp, err = http.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}

    if err != nil{
    	errorHandler(w, r, http.StatusInternalServerError, err)
		return
    }

    if resp.StatusCode != 200{
    	http.Error(w, resp.Status, resp.StatusCode)
    	return
    }

    if docErr != nil{
    	http.Error(w, docErr.Error(), http.StatusInternalServerError)
    }

    defer resp.Body.Close()

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
	url := baseURL + "advanced-search?"
	queryParams := r.URL.Query()

	var (
		doc *goquery.Document
		err, docErr error
		resp *http.Response
	)

	if queryParams.Get("key") != ""{
		url += "key=" + queryParams.Get("key")
	}

	if queryParams.Get("include") != ""{
		url += "&wg=" + queryParams.Get("include")
	}

	if queryParams.Get("exclude") != ""{
		url += "&wog=" + queryParams.Get("exclude")
	}

	if queryParams.Get("status") != ""{
		url += "&status=" + queryParams.Get("status")
	}

	if queryParams.Get("page") != ""{
		url += "&page=" + queryParams.Get("page")
	}

	type SearchResult struct{
		Title string `json:"title"`
		Link string `json:"link"`
		Img string 	`json:"img"`
		Genre string  `json:"genre"`
	}

	var results []SearchResult

	if appengine.IsDevAppServer(){
		c := appengine.NewContext(r)
		client := urlfetch.Client(c)
		resp, err = client.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}else{
		resp, err = http.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}

    if err != nil{
    	errorHandler(w, r, http.StatusInternalServerError,err)
		return
    }

    if resp.StatusCode != 200{	
    	http.Error(w, resp.Status, resp.StatusCode)
    	return
    }

    if docErr != nil{
    	http.Error(w, docErr.Error(), http.StatusInternalServerError)
    	return
    }

    defer resp.Body.Close()

    doc.Find(".manga-box").Each(func(index int, item *goquery.Selection){
    	comic := SearchResult{}
    	comic.Title = item.Find("h3").Children().Text()
    	comic.Link, _ = item.Find("h3").Children().Attr("href")
    	comic.Img, _ = item.Find("img").Attr("src")
    	comic.Genre = strings.TrimSpace(item.Find(".genre").Text())
    	results = append(results, comic)
    })

    w.Header().Set("Content-Type", "application/json")
    res, _ := json.Marshal(results)
	fmt.Fprintf(w,string(res))
}

//Reqest Format -> /description?name={Comic Name}
func getDescription(w http.ResponseWriter, r *http.Request) {

	url := baseURL + "comic/"
	queryParams := r.URL.Query()
	var (
		doc *goquery.Document
		err, docErr error
		resp *http.Response
	)

	if queryParams.Get("name") != ""{
		url += queryParams.Get("name")
	}else{
		errorHandler(w, r, http.StatusBadRequest,err)
		return
	}

	if appengine.IsDevAppServer(){
		c := appengine.NewContext(r)
		client := urlfetch.Client(c)
		resp, err = client.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}else{
		resp, err = http.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}

    if err != nil{
    	errorHandler(w, r, http.StatusInternalServerError, err)
		return
    }

    if resp.StatusCode != 200{
    	http.Error(w, resp.Status, resp.StatusCode)
    	return
    }

    if docErr != nil{
    	http.Error(w, docErr.Error(), http.StatusInternalServerError )
    }

    defer resp.Body.Close()

    type Response struct{
    	Description string `json:"description"`
    	LargeImg string `json:"largeImg"`
    	FormatedName string  `json:"formatedName"`
    	Name string `json:"name"`
    	AlternateName string `json:"alternateName"`
    	ReleaseYear string `json:"releaseYear"`
    	Status string `json:"status"`
    	Author string `json:"author"`
    	Genre string `json:"genre"`
    }
    
 	description := Response{ FormatedName: queryParams.Get("name")}

    doc.Find(".manga-details").Find("tr").Each(func(index int, item *goquery.Selection){
    	//Chain of ugly if statetments
    	pair := strings.Split(item.First().Children().Text(), ":")
    	key := strings.TrimSpace(pair[0])
    	val := strings.TrimSpace(pair[1])

    	if key == "Name" {
    		description.Name = val
    	}else if key == "Alternate Name"{
    		description.AlternateName = val	
    	}else if key == "Year of Release"{
    		description.ReleaseYear = val
    	}else if key == "Status"{
    		description.Status = val
    	}else if key == "Author"{
    		description.Author = val
    	}else if key == "Genre"{
    		description.Genre = val
    	}


    })
    description.LargeImg, _ = doc.Find(".manga-image").ChildrenFiltered("img").Attr("src")
    description.Description = strings.TrimSpace(doc.Find(".pdesc").Text())

    w.Header().Set("Content-Type", "application/json")
    res, _ := json.Marshal(description)
    fmt.Fprintf(w, string(res))

}
func main() {

	http.HandleFunc("/comic-list-AZ", allComicsRequest)
	http.HandleFunc("/popular-comics/",getPopularComics)
	http.HandleFunc("/chapter-list/", getChapters)
	http.HandleFunc("/read-comic/", readComic)
	http.HandleFunc("/search-categories", getSearchCategories)
	http.HandleFunc("/advanced-search", performAdvancedSearch)
	http.HandleFunc("/description", getDescription)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello fellow, Gruntt User!")
		return 
	})
	http.ListenAndServe(":8000",nil)
	appengine.Main()

}

