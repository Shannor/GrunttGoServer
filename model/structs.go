package model

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
