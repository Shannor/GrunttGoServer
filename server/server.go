package server

import (
	"fmt"
	"net/http"
	"scrapper/model"
	"scrapper/workers"

	"github.com/julienschmidt/httprouter"
)

// //Request Format -> /search-categories
// func getSearchCategories(w http.ResponseWriter, r *http.Request) {
// 	url := COMICEXTRA + "advanced-search"
// 	var (
// 		categories []string
// 	)
//
// 	doc, err := utils.GetGoQueryDoc(url, r)
// 	if err != nil {
// 		model.ErrorHandler(w, http.StatusBadRequest, err)
// 	}
//
// 	doc.Find(".search-checks").ChildrenFiltered("li").Each(func(index int, item *goquery.Selection) {
// 		categories = append(categories, item.Text())
// 	})
//
// 	w.Header().Set("Content-Type", "application/json")
// 	res, _ := json.Marshal(categories)
// 	fmt.Fprintf(w, string(res))
//
// }
//
// //Request Format -> /advanced-search?key={Search Param}&include={Desired Attributes}
// //&exclude={Undesired attributes}&status={status}&page={Page number for results}
// func performAdvancedSearch(w http.ResponseWriter, r *http.Request) {
// 	url := COMICEXTRA + "advanced-search?"
// 	queryParams := r.URL.Query()
//
// 	if queryParams.Get("key") != "" {
// 		url += "key=" + queryParams.Get("key")
// 	}
//
// 	if queryParams.Get("include") != "" {
// 		url += "&wg=" + queryParams.Get("include")
// 	}
//
// 	if queryParams.Get("exclude") != "" {
// 		url += "&wog=" + queryParams.Get("exclude")
// 	}
//
// 	if queryParams.Get("status") != "" {
// 		url += "&status=" + queryParams.Get("status")
// 	}
//
// 	if queryParams.Get("page") != "" {
// 		url += "&page=" + queryParams.Get("page")
// 	}
//
// 	type SearchResult struct {
// 		Title string `json:"title"`
// 		Link  string `json:"link"`
// 		Img   string `json:"img"`
// 		Genre string `json:"genre"`
// 	}
//
// 	var results []SearchResult
//
// 	doc, err := utils.GetGoQueryDoc(url, r)
// 	if err != nil {
// 		model.ErrorHandler(w, http.StatusBadRequest, err)
// 	}
//
// 	doc.Find(".manga-box").Each(func(index int, item *goquery.Selection) {
// 		comic := SearchResult{}
// 		comic.Title = item.Find("h3").Children().Text()
// 		comic.Link, _ = item.Find("h3").Children().Attr("href")
// 		comic.Img, _ = item.Find("img").Attr("src")
// 		comic.Genre = strings.TrimSpace(item.Find(".genre").Text())
// 		results = append(results, comic)
// 	})
//
// 	w.Header().Set("Content-Type", "application/json")
// 	res, _ := json.Marshal(results)
// 	fmt.Fprintf(w, string(res))
// }
//
// //Reqest Format -> /description?name={Comic Name}
// func getDescription(w http.ResponseWriter, r *http.Request) {
//
// 	url := COMICEXTRA + "comic/"
// 	queryParams := r.URL.Query()
//
// 	if queryParams.Get("name") != "" {
// 		url += queryParams.Get("name")
// 	} else {
// 		return
// 	}
//
// 	doc, err := utils.GetGoQueryDoc(url, r)
// 	if err != nil {
// 		return
// 	}
// 	type Response struct {
// 		Description   string `json:"description"`
// 		LargeImg      string `json:"largeImg"`
// 		FormatedName  string `json:"formatedName"`
// 		Name          string `json:"name"`
// 		AlternateName string `json:"alternateName"`
// 		ReleaseYear   string `json:"releaseYear"`
// 		Status        string `json:"status"`
// 		Author        string `json:"author"`
// 		Genre         string `json:"genre"`
// 	}
//
// 	description := Response{FormatedName: queryParams.Get("name")}
//
// 	doc.Find(".manga-details").Find("tr").Each(func(index int, item *goquery.Selection) {
// 		//Chain of ugly if statetments
// 		pair := strings.Split(item.First().Children().Text(), ":")
// 		key := strings.TrimSpace(pair[0])
// 		val := strings.TrimSpace(pair[1])
//
// 		if key == "Name" {
// 			description.Name = val
// 		} else if key == "Alternate Name" {
// 			description.AlternateName = val
// 		} else if key == "Year of Release" {
// 			description.ReleaseYear = val
// 		} else if key == "Status" {
// 			description.Status = val
// 		} else if key == "Author" {
// 			description.Author = val
// 		} else if key == "Genre" {
// 			description.Genre = val
// 		}
//
// 	})
// 	description.LargeImg, _ = doc.Find(".manga-image").ChildrenFiltered("img").Attr("src")
// 	description.Description = strings.TrimSpace(doc.Find(".pdesc").Text())
//
// 	w.Header().Set("Content-Type", "application/json")
// 	res, _ := json.Marshal(description)
// 	fmt.Fprintf(w, string(res))
//
// }

func init() {

	ce := workers.GetComicExtraInstance()
	rcw := workers.GetReadComicsInstance()
	api := model.GetAPIInstance(ce, rcw)
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Welcome to Gruntt-Comics Backend\n")
	})
	router.GET("/all-comics", api.GetAllComics())
	router.GET("/popular-comics/:pageNumber", api.GetPopularComics())
	router.GET("/chapter-list/:comicName", api.GetComicChapterList())
	router.GET("/chapter-pages/:comicName/:chapterNumber", api.GetChapterPages())
	router.GET("/description/:comicName", api.GetComicDescription())
	http.Handle("/", router)
}
