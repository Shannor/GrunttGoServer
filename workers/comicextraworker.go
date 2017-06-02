package workers

import (
    "strconv"
    "scrapper/model"
    "github.com/PuerkitoBio/goquery"

)

const ComicExtraURL = "http://www.comicextra.com/"
const ComicExtraURLParam = "ce"


func GetAllComics(doc *goquery.Document)([]model.Comic) {
    var comics []model.Comic
    doc.Find(".series-col li").Each(func(index int, item *goquery.Selection) {
        comic := model.Comic{}
        aTag := item.Children()
        comic.Title = aTag.Text()
        comic.Link, _ = aTag.Attr("href")
        comic.Category = item.Parent().SiblingsFiltered("div").Text()

        if _, err := strconv.Atoi(comic.Category); err == nil {
            comic.Category = "#"
        }

        comics = append(comics, comic)
    })
    return comics
}
