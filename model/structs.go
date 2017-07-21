package model

import (
	"fmt"
	"net/http"
)

type Comic struct {
	Title    string `json:"title"`
	Link     string `json:"link"`
	Category string `json:"category"`
}

type PopularComic struct {
	Title      string `json:"title"`
	Link       string `json:"link"`
	Img        string `json:"img"`
	IssueCount int    `json:"issueCount"`
}

type Chapter struct {
	ChapterName string `json:"chapterName"`
	Link        string `json:"link"`
	ReleaseDate string `json:"releaseDate"`
}

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
