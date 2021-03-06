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

	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	resp, httpErr = client.Get(url)
	doc, docErr = goquery.NewDocumentFromResponse(resp)

	defer resp.Body.Close()

	if httpErr != nil {
		return nil, fmt.Errorf("GetGoQueryDoc: Http error: %s", httpErr.Error())
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GetGoQueryDoc: Response code is %d not 200", resp.StatusCode)
	}

	if docErr != nil {
		return nil, fmt.Errorf("GetGoQueryDoc: Document error: %s", docErr.Error())
	}

	return doc, nil
}
