package workers

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//
// import (
// 	"strconv"
// 	"strings"
//
// 	"github.com/PuerkitoBio/goquery"
// )
//
// const (
// 	ReadComicsURL      = "http://readcomics.website/"
// 	ReadComicsURLParam = "rcw"
// )
func (rcw *ReadComics) MatchesTag(tag string) bool {
	if tag == rcw.Tag {
		return true
	}
	return false
}

func (rcw *ReadComics) CreateAllComicsURL() string {
	return rcw.BaseURL + "changeMangaList?type=text"
}

func (rcw *ReadComics) CreatePopularComicsURL(pageNumber int) string {
	return rcw.BaseURL + "filterList?page=" + strconv.Itoa(pageNumber) + "&sortBy=views&asc=fals"
}

func (rcw *ReadComics) CreateComicChapterListURL(comicName string) string {
	return rcw.BaseURL + "comic/" + comicName
}
func (rcw *ReadComics) CreateChapterPagesURL(comicName string, chapterNumber int) string {
	return rcw.BaseURL + "comic/" + comicName + "/" + strconv.Itoa(chapterNumber)
}

func (rcw *ReadComics) CreateComicDescriptionURL(comicName string) string {
	return rcw.BaseURL + "comic/" + comicName
}

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

//
// func GetChaptersReadComics(doc *goquery.Document) []Chapter {
// 	var chapters []Chapter
// 	doc.Find(".chapters").Children().Each(func(index int, item *goquery.Selection) {
// 		chapter := Chapter{}
// 		obj := item.Find("a")
// 		chapter.Link, _ = obj.Attr("href")
// 		chapter.ChapterName = obj.Text()
// 		chapter.ReleaseDate = strings.TrimSpace(item.Find(".date-chapter-title-rtl").Text())
// 		chapters = append(chapters, chapter)
// 	})
//
// 	return chapters
// }
