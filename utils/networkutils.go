package utils

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

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
		return nil, fmt.Errorf("Some error")
	}

	if docErr != nil {
		return nil, docErr
	}

	return doc, nil
}
