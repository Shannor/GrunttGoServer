package utils

import (
	"fmt"
	"net/http"
	"scrapper/model"
	"scrapper/workers"

	"github.com/PuerkitoBio/goquery"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

// RequestError  error for when the provided param doesn't match
type RequestError struct {
	ProvidedParam string
}

//ResponseError error for when there is a response but not a 200
type ResponseError struct {
	ResponseCode int
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("Provided Param( '%s' ) does not match any options.",
		e.ProvidedParam)
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("Error Response code: %d", e.ResponseCode)
}

//ErrorHandler function to write all errors to the page
func ErrorHandler(w http.ResponseWriter, status int, err error) {
	http.Error(w, err.Error(), status)
}

//CreateAllComicsURL function to create the URL for each source
func CreateAllComicsURL(urlParam string) (string, error) {
	switch urlParam {
	case workers.ComicExtraURLParam:
		return workers.ComicExtraURL + "comic-list", nil
	case workers.ReadComicsURLParam:
		return workers.ReadComicsURL + "changeMangaList?type=text", nil
	default:
		return "", &RequestError{urlParam}
	}
}

func CreatePopularComicsURL(urlParam string, pageNumber string) (string, error) {
	switch urlParam {
	case workers.ComicExtraURLParam:
		return workers.ComicExtraURL + "popular-comic/" + pageNumber, nil
	case workers.ReadComicsURLParam:
		return workers.ReadComicsURL + "filterList?page=" +
			pageNumber + "&sortBy=views&asc=false", nil
	default:
		return "", &RequestError{urlParam}
	}
}

func CreateChapterURL(urlParam string, comicName string) (string, error) {
	switch urlParam {
	case workers.ComicExtraURLParam:
		return workers.ComicExtraURL + "comic/" + comicName, nil
	case workers.ReadComicsURLParam:
		return workers.ReadComicsURL + "comic/" + comicName, nil
	default:
		return "", &RequestError{urlParam}

	}
}

func CreateReadComicURL(urlParam string, comicName string, chapterNumber string) (string, error) {
	switch urlParam {
	case workers.ComicExtraURLParam:
		return workers.ComicExtraURL + comicName + "/chapter-" + chapterNumber, nil
	case workers.ReadComicsURLParam:
		return workers.ReadComicsURL + "comic/" + comicName + "/" + chapterNumber, nil
	default:
		return "", &RequestError{urlParam}

	}
}

//GetGoQueryDoc Helper method to return the html body to parse
func GetGoQueryDoc(url string, r *http.Request) (*goquery.Document, error) {

	var (
		doc             *goquery.Document
		httpErr, docErr error
		resp            *http.Response
	)

	if appengine.IsDevAppServer() {
		c := appengine.NewContext(r)
		client := urlfetch.Client(c)
		resp, httpErr = client.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)

	} else {
		resp, httpErr = http.Get(url)
		doc, docErr = goquery.NewDocumentFromResponse(resp)
	}

	defer resp.Body.Close()

	if httpErr != nil {
		return nil, httpErr
	}

	if resp.StatusCode != 200 {
		return nil, &ResponseError{resp.StatusCode}
	}

	if docErr != nil {
		return nil, docErr
	}

	return doc, nil
}

//TODO: Add err returns to all Get methods
func GetAllComics(doc *goquery.Document, param string) []model.Comic {
	switch param {
	case workers.ComicExtraURLParam:
		return workers.GetAllComicsComicExtra(doc)
	case workers.ReadComicsURLParam:
		return workers.GetAllComicsReadComics(doc)
	default:
		return nil
	}
}

func GetPopularComics(doc *goquery.Document, param string) []model.PopularComic {
	switch param {
	case workers.ComicExtraURLParam:
		return workers.GetPopularComicsComicExtra(doc)
	case workers.ReadComicsURLParam:
		return workers.GetPopularComicsReadComics(doc)
	default:
		return nil
	}
}

func GetChapters(doc *goquery.Document, r *http.Request, param string, comicName string) []model.Chapter {
	switch param {
	case workers.ComicExtraURLParam:
		return workers.GetChaptersComicExtra(doc, r, comicName)
	case workers.ReadComicsURLParam:
		return workers.GetChaptersReadComics(doc)
	default:
		return nil
	}
}

func GetChapterImages(doc *goquery.Document, r *http.Request, param string, url string) []string {
	switch param {
	case workers.ComicExtraURLParam:

	case workers.ReadComicsURLParam:

	default:
		return nil
	}
}
