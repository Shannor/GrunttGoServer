package workers

import (
	"fmt"
	"net/http"
	"scrapper/utils"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//GetTag Returns the tag of readComics
func (rcw *ReadComics) GetTag() string {
	return rcw.Tag
}

//CreateAllComicsURL returns the url for all comics
func (rcw *ReadComics) CreateAllComicsURL() string {
	return rcw.BaseURL + "changeMangaList?type=text"
}

//CreatePopularComicsURL returns the popular comics url
func (rcw *ReadComics) CreatePopularComicsURL(pageNumber int) string {
	return rcw.BaseURL + "filterList?page=" + strconv.Itoa(pageNumber) + "&sortBy=views&asc=fals"
}

//CreateComicChapterListURL returns url for chapter list
func (rcw *ReadComics) CreateComicChapterListURL(comicName string) string {
	return rcw.BaseURL + "comic/" + comicName
}

//CreateChapterPagesURL return the url for the pages for readcomics
func (rcw *ReadComics) CreateChapterPagesURL(comicName string, chapterNumber int) string {
	return rcw.BaseURL + "comic/" + comicName + "/" + strconv.Itoa(chapterNumber)
}

//CreateComicDescriptionURL creates the description url for a comic
func (rcw *ReadComics) CreateComicDescriptionURL(comicName string) string {
	return rcw.BaseURL + "comic/" + comicName
}

//GetAllComics method to scrape comics from readcomics
func (rcw *ReadComics) GetAllComics(doc *goquery.Document) (Comics, error) {
	var comics Comics
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
	return comics, nil
}

//GetPopularComics method to return the popluar comics on the website
func (rcw *ReadComics) GetPopularComics(doc *goquery.Document) (PopularComics, error) {
	var comics PopularComics

	doc.Find(".media").Each(func(index int, item *goquery.Selection) {
		comic := PopularComic{}
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

		comics = append(comics, comic)
	})
	return comics, nil
}

//GetComicChapterListPageAmount returns the max number of chapters pages
func (rcw *ReadComics) GetComicChapterListPageAmount(doc *goquery.Document) int {
	return 1
}

//GetComicChapterList Returns all the chapters listed on the page
func (rcw *ReadComics) GetComicChapterList(comicName string, numOfPages int, r *http.Request) (Chapters, error) {
	var chapters Chapters
	url := rcw.BaseURL + "comic/" + comicName

	doc, err := utils.GetGoQueryDoc(url, r)
	if err != nil {
		return Chapters{}, fmt.Errorf("GetComicChapterList error. Error: %s", err.Error())
	}

	doc.Find(".chapters").Children().Each(func(index int, item *goquery.Selection) {
		chapter := Chapter{}
		obj := item.Find("a")
		chapter.Link, _ = obj.Attr("href")
		chapter.ChapterName = obj.Text()
		chapter.ReleaseDate = strings.TrimSpace(item.Find(".date-chapter-title-rtl").Text())
		chapters = append(chapters, chapter)
	})

	return chapters, nil
}

//GetNumberOfPages return the number of pages in a chapter
func (rcw *ReadComics) GetNumberOfPages(doc *goquery.Document) int {
	return doc.Find(".selectpicker").First().Children().Length()
}

//GetChapterPages return the urls for the pages in a chapter
func (rcw *ReadComics) GetChapterPages(comicName string, chapterNumber int, numOfPages int, r *http.Request) ([]string, error) {
	baseURL := rcw.BaseURL + "comic/" + comicName + "/" + strconv.Itoa(chapterNumber)
	pagesChannels := make(chan string, numOfPages)
	var urls []string
	go rcw.getComicImageURL(baseURL, numOfPages, r, pagesChannels)

	for url := range pagesChannels {
		urls = append(urls, url)
	}
	return urls, nil
}

//GetComicImageURL go routine to get the acutal urls
func (rcw *ReadComics) getComicImageURL(url string, numOfPages int, r *http.Request, cc chan string) {
	for i := 1; i <= numOfPages; i++ {
		pageURL := url + "/" + strconv.Itoa(i)

		doc, err := utils.GetGoQueryDoc(pageURL, r)
		if err != nil {
			return
		}
		link, _ := doc.Find("#ppp").Find("img").Attr("src")
		cc <- strings.TrimSpace(link)
	}
	close(cc)
}

func (rcw *ReadComics) GetComicDescription(doc *goquery.Document) (Description, error) {
	var description Description
	description.Name = strings.TrimSpace(doc.Find(".listmanga-header").First().Text())
	description.Description = strings.TrimSpace(doc.Find(".manga.well").ChildrenFiltered("p").Text())
	description.LargeImg, _ = doc.Find(".img-responsive").Attr("src")

	doc.Find(".dl-horizontal").Find("dd").Each(func(index int, item *goquery.Selection) {
		switch index {
		case 0:
			description.Genre = strings.TrimSpace(item.Text())
		case 1:
			description.Status = strings.TrimSpace(item.Text())
		case 2:
			description.AlternateName = strings.TrimSpace(item.Text())
		case 3:
			description.Author = strings.TrimSpace(item.Text())
		case 4:
			description.ReleaseYear = strings.TrimSpace(item.Text())
		}
	})
	return description, nil
}
