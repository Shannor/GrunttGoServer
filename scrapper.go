package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"scrapper/utils"

	"github.com/PuerkitoBio/goquery"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
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

//Request Format -> /popular-comics/{page Number}?url={website}
//TODO: change pagenumber to a query param
func getPopularComics(w http.ResponseWriter, r *http.Request) {

	pageNumber := r.URL.Path[len("/popular-comics/"):]
	//TODO: check if pageNumber is a number and not more than that
	if _, err := strconv.Atoi(pageNumber); err != nil {
		//Throw wrong formated request type error
		return
	}

	var (
		doc         *goquery.Document
		err, docErr error
		resp        *http.Response
		url         string
	)

	choice := r.URL.Query().Get("url")
	switch choice {
	case comicExtraURLParam:
		url = COMICEXTRA + "popular-comic/" + pageNumber
	case readcomicsURLParam:
		url = READCOMIC + "filterList?page=" + pageNumber + "&sortBy=views&asc=false"
	default:
		return
	}

	if appengine.IsDevAppServer() {
		c := appengine.NewContext(r)
		client := urlfetch.Client(c)
		resp, err = client.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	} else {
		resp, err = http.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}

	//Error with inital url request
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != 200 {
		http.Error(w, resp.Status, resp.StatusCode)
		return
	}

	defer resp.Body.Close()

	//Error with goquery
	if docErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type PopComic struct {
		Title      string `json:"title"`
		Link       string `json:"link"`
		Img        string `json:"img"`
		IssueCount int    `json:"issueCount"`
	}

	var popularcomics []PopComic

	switch choice {
	case comicExtraURLParam:
		doc.Find(".cartoon-box").Each(func(index int, item *goquery.Selection) {
			comic := PopComic{}
			//Gets top level information
			comic.Title = item.Find("h3").Children().Text()
			comic.Link, _ = item.Find("h3").Children().Attr("href")
			comic.Img, _ = item.Find("img").Attr("src")
			count := item.Find(".detail").First().Text()
			split := strings.Split(count, " ")
			val, err := strconv.Atoi(split[0])
			if err == nil {
				comic.IssueCount = val
			}

			popularcomics = append(popularcomics, comic)
		})
	case readcomicsURLParam:
		doc.Find(".media").Each(func(index int, item *goquery.Selection) {
			comic := PopComic{}
			//Gets top level information
			obj := item.Find(".chart-title")
			comic.Title = obj.Text()
			comic.Link, _ = obj.Attr("href")
			comic.Img, _ = item.Find("img").Attr("src")
			count := item.Find("i").Parent().Text()
			count = strings.TrimSpace(count)
			val, err := strconv.Atoi(count)
			if err == nil {
				comic.IssueCount = val
			}

			popularcomics = append(popularcomics, comic)
		})
	default:
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(popularcomics)
	w.Write(res)
}

//Chapter : Struct to represent chapters
type Chapter struct {
	ChapterName string `json:"chapterName"`
	Link        string `json:"link"`
	ReleaseDate string `json:"releaseDate"`
}

//Request Format -> /chapter-list/{Comic Name}?url={type}
func getChapters(w http.ResponseWriter, r *http.Request) {

	comicName := r.URL.Path[len("/chapter-list/"):]
	// TODO: Check if name is right or wrong
	var (
		url string
	)

	choice := r.URL.Query().Get("url")
	switch choice {
	case comicExtraURLParam:
		url = COMICEXTRA + "comic/" + comicName
	case readcomicsURLParam:
		url = READCOMIC + "comic/" + comicName
	default:
		return
	}

	doc, err := utils.GetGoQueryDoc(url, r)
	if err != nil {
		utils.ErrorHandler(w, http.StatusBadRequest, err)
	}

	var chapters []Chapter
	switch choice {
	case comicExtraURLParam:

		doc.Find("#list").Children().Each(func(index int, item *goquery.Selection) {
			chapter := Chapter{}
			obj := item.Find("a")
			chapter.Link, _ = obj.Attr("href")
			chapter.ChapterName = obj.Text()
			chapter.ReleaseDate = item.Find("td").Last().Text()
			chapters = append(chapters, chapter)
		})

		pageCount := doc.Find(".general-nav").Children().Length() - 1

		if pageCount > 0 {
			chapterChannels := make(chan []Chapter, pageCount)

			go getExtraChapters(pageCount, r, comicName, chapterChannels)

			for i := range chapterChannels {
				chapters = append(chapters, i...)
			}
		}
	case readcomicsURLParam:
		doc.Find(".chapters").Children().Each(func(index int, item *goquery.Selection) {
			chapter := Chapter{}
			obj := item.Find("a")
			chapter.Link, _ = obj.Attr("href")
			chapter.ChapterName = obj.Text()
			chapter.ReleaseDate = strings.TrimSpace(item.Find(".date-chapter-title-rtl").Text())
			chapters = append(chapters, chapter)
		})

	default:
		log.Println("No or incorrcect URL param.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(chapters)
	w.Write(res)
}

func getExtraChapters(extras int, r *http.Request, comicName string, cc chan []Chapter) {
	for i := 2; i <= extras; i++ {

		var (
			doc         *goquery.Document
			err, docErr error
			resp        *http.Response
		)
		url := COMICEXTRA + "comic/" + comicName + "/" + strconv.Itoa(i)

		if appengine.IsDevAppServer() {
			c := appengine.NewContext(r)
			client := urlfetch.Client(c)
			resp, err = client.Get(url)
			doc, docErr = goquery.NewDocumentFromResponse(resp)
		} else {
			resp, err = http.Get(url)
			doc, docErr = goquery.NewDocumentFromResponse(resp)
		}

		if err != nil {
			log.Printf(err.Error())
			return
		}

		if resp.StatusCode != 200 {
			log.Printf(resp.Status)
			return
		}

		if docErr != nil {
			log.Printf(docErr.Error())
			return
		}

		defer resp.Body.Close()

		var chapters []Chapter
		doc.Find("#list").Children().Each(func(index int, item *goquery.Selection) {
			chapter := Chapter{}
			obj := item.Find("a")
			chapter.Link, _ = obj.Attr("href")
			chapter.ChapterName = obj.Text()
			chapter.ReleaseDate = item.Find("td").Last().Text()
			chapters = append(chapters, chapter)

		})

		cc <- chapters
	}
	close(cc)
}

//Request Format -> /read-comic/{Comic Name}/{Chapter Number}?url={param}
func readComic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/read-comic/"):]
	//0 = ComicName, 1 = Chapter Number
	paths := strings.Split(path, "/")
	//TODO: check how split returns resutls
	if len(paths) < 2 {
		// errorHandler(w, r, http.StatusBadRequest,err)
		return
	}

	var (
		url string
	)

	choice := r.URL.Query().Get("url")
	switch choice {
	case comicExtraURLParam:
		url = COMICEXTRA + paths[0] + "/chapter-" + paths[1]
	case readcomicsURLParam:
		url = READCOMIC + "comic/" + paths[0] + "/" + paths[1]
	default:
		return
	}

	doc, err := utils.GetGoQueryDoc(url, r)
	if err != nil {
		utils.ErrorHandler(w, http.StatusBadRequest, err)
	}
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
		var (
			doc         *goquery.Document
			err, docErr error
			resp        *http.Response
		)

		if appengine.IsDevAppServer() {
			c := appengine.NewContext(r)
			client := urlfetch.Client(c)
			resp, err = client.Get(pageURL)
			doc, docErr = goquery.NewDocumentFromResponse(resp)
		} else {
			resp, err = http.Get(pageURL)
			doc, docErr = goquery.NewDocumentFromResponse(resp)
		}

		if err != nil {
			log.Printf(err.Error())
			return
		}
		if resp.StatusCode != 200 {
			log.Printf(resp.Status)
			return
		}

		if docErr != nil {
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
	var (
		doc         *goquery.Document
		err, docErr error
		resp        *http.Response
	)

	if queryParams.Get("name") != "" {
		url += queryParams.Get("name")
	} else {
		utils.ErrorHandler(w, http.StatusBadRequest, err)
		return
	}

	if appengine.IsDevAppServer() {
		c := appengine.NewContext(r)
		client := urlfetch.Client(c)
		resp, err = client.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	} else {
		resp, err = http.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}

	if err != nil {
		utils.ErrorHandler(w, http.StatusInternalServerError, err)
		return
	}

	if resp.StatusCode != 200 {
		http.Error(w, resp.Status, resp.StatusCode)
		return
	}

	if docErr != nil {
		http.Error(w, docErr.Error(), http.StatusInternalServerError)
	}

	defer resp.Body.Close()

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
