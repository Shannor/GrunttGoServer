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

type Chapter struct {
	ChapterName string `json:"chapterName"`
	Link        string `json:"link"`
	ReleaseDate string `json:"releaseDate"`
}
