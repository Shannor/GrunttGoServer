package utils

import (
    "net/http"
    "fmt"
    "strconv"
    "strings"
    "scrapper/workers"
    "github.com/PuerkitoBio/goquery"
    "google.golang.org/appengine"
    "google.golang.org/appengine/urlfetch"
)

const READCOMIC = "http://readcomics.website/"
const readcomicsURLParam = "rcw"

type RequestError struct {
    ProvidedParam string
}

type ResponseError struct{
    ResponseCode int
}

func (e *RequestError) Error() string {
    return fmt.Sprintf("Provided Param( %s ) does not match any options.",
     e.ProvidedParam)
}

func (e *ResponseError) Error() string{
    return fmt.Sprintf("Error Response code: %d", e.ResponseCode)
}

func ErrorHandler(w http.ResponseWriter, status int, err error) {
	http.Error(w, err.Error(), status)
}


func CreateAllComicsURL(urlParam string)(string, error){
    switch urlParam {
    case workers.ComicExtraURLParam:
        return workers.ComicExtraURL + "comic-list", nil
    case readcomicsURLParam:
        return READCOMIC + "changeMangaList?type=text", nil
    default:
        return "", &RequestError{urlParam}
    }
}

func GetGoQueryDoc(url string, r *http.Request)(*goquery.Document, error) {

    var (
        doc         *goquery.Document
        httpErr, docErr error
        resp        *http.Response
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
        return nil,docErr
    }

    return doc, nil
}

func GetAllComics(doc *goquery.Document, param string) []Comic{
    switch param {
    case workers.ComicExtraURLParam:
        return workers.GetAllComics(doc)
    case readcomicsURLParam:
    	doc.Find(".type-content li").Each(func(index int, item *goquery.Selection) {
    		comic := Comic{}
            aTag := item.ChildrenFiltered("a")
    		comic.Title = strings.TrimSpace(aTag.Text())
    		comic.Link, _ = aTag.Attr("href")

    		if _, err := strconv.Atoi(string(comic.Title[0])); err == nil {
    			comic.Category = "#"
    		} else {
    			comic.Category = string(comic.Title[0])
    		}
    		comics = append(comics, comic)
    	})
        return comics
    default:
    	return nil
    }
}
