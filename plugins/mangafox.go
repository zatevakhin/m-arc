package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"../core/types"
	"../core/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

// GetPluginName - Plugin name
func GetPluginName() string {
	return "fanfox.net"
}

// GetPluginVersion - Version number
func GetPluginVersion() string {
	return "0.0.1"
}

// MangaInfo - Object contains data about manga
type MangaInfo types.MangaInfo

// MangaInfoExport - exported
var MangaInfoExport MangaInfo

func (manga *MangaInfo) SetUrl(url string) {
	manga.URL = url
}

func (manga *MangaInfo) GetData() types.MangaInfo {
	return *manga
}

func (manga *MangaInfo) GetMetaData() {
	result := getPageByURL(manga.URL)
	parseMangaMetadata(manga, result)
}

func getTranslationStatus(status *string) types.TranslationStatus {
	var statusCode types.TranslationStatus

	return statusCode
}

func getMangaStatus(status *string) types.MangaCompletionStatus {
	var statusCode types.MangaCompletionStatus
	switch *status {
	case "Completed":
		statusCode = types.MS_FINISHED
	case "Ongoing":
		statusCode = types.MS_IN_PROGRESS
	default:
		statusCode = types.MS_UNKNOWN
	}
	return statusCode
}

func getPageByURL(url string) *http.Response {

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(``)))

	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	cookie := http.Cookie{Name: "isAdult", Value: "1"}
	req.AddCookie(&cookie)

	client := &http.Client{Timeout: time.Second * 10}

	result, err := client.Do(req)

	if nil != err {
		logrus.Fatal(err)
	}

	if 200 != result.StatusCode {
		logrus.WithFields(logrus.Fields{
			"url":  url,
			"code": result.StatusCode,
		}).Fatal("Can't get page.")
	}

	logrus.Debug(result)

	return result
}

func parseMangaMetadata(manga *MangaInfo, response *http.Response) {
	document, err := goquery.NewDocumentFromResponse(response)

	if nil != err {
		log.Fatal(err)
	}

	parsedURL, err := url.Parse(manga.URL)

	if err != nil {
		logrus.Fatalf("The URL wasn't parsed!")
	}

	detailInfoRight := document.Find("div.detail-info-right")

	manga.Title = detailInfoRight.Find("span.detail-info-right-title-font").Text()

	status := detailInfoRight.Find("span.detail-info-right-title-tip").Text()
	manga.MangaStatus = getMangaStatus(&status)

	manga.Description = detailInfoRight.Find("p.fullcontent").Text()

	genres := detailInfoRight.Find("p.detail-info-right-tag-list")

	genres.Find("a").Each(func(i int, s *goquery.Selection) {
		manga.Genres = append(manga.Genres, strings.Trim(s.Text(), " "))
	})

	chapters := document.Find("ul.detail-main-list")

	defaultRegex := regexp.MustCompile(`Vol\.(?P<volume>[\d.]+)\s+Ch\.(?P<chapter>[\d.]+)(?:(?:\s-)+?(?P<name>.+))?$`)

	hostname := parsedURL.Hostname()
	scheme := parsedURL.Scheme

	chapterItems := chapters.Find("li")
	chapterItemsCount := chapterItems.Length()
	chapterItems.Each(func(i int, s *goquery.Selection) {
		link := s.Find("a")

		hrefAttr, _ := link.Attr("href")
		titleAttr, _ := link.Attr("title")

		result := utils.FindNamedMatches(defaultRegex, titleAttr)

		var chapter types.MangaChapter

		chapter.Index = chapterItemsCount
		chapterItemsCount--

		chapter.VolumeNumber = result["volume"]
		chapter.ChapterNumber = result["chapter"]
		chapter.IsSpecial = strings.Contains(result["chapter"], ".")
		chapter.Name = strings.Trim(result["name"], " ")
		chapter.URL = fmt.Sprintf("%s://%s%s", scheme, hostname, hrefAttr)

		manga.Chapters = append(manga.Chapters, chapter)
	})
}

func (manga *MangaInfo) GetChapters() {

	for i := range manga.Chapters {
		result := getPageByURL(manga.Chapters[i].URL)
		parseMangaChapter(&manga.Chapters[i], result)
		break
	}

}

func parseMangaChapter(chapter *types.MangaChapter, response *http.Response) {
	document, err := goquery.NewDocumentFromResponse(response)

	if nil != err {
		log.Fatal(err, document)
	}

	var foundScript string

	document.Find("script").EachWithBreak(func(i int, s *goquery.Selection) bool {
		script := s.Text()

		isContains := strings.Contains(script, "comicid")
		isContains = isContains && strings.Contains(script, "chapterid")
		isContains = isContains && strings.Contains(script, "imagepage")
		isContains = isContains && strings.Contains(script, "imagecount")

		if isContains {
			foundScript = script
		}

		return !isContains
	})

	defaultRegex := regexp.MustCompile(`comicid[^\d]+(?P<comicid>\d+).*chapterid[^\d]+(?P<chapterid>\d+).*imagepage[^\d]+(?P<imagepage>\d+).*imagecount[^\d]+(?P<imagecount>\d+)`)
	result := utils.FindNamedMatches(defaultRegex, foundScript)

	document.Find("div.pager-list-left")

	log.Println(result)
}

func appendChapter(manga *MangaInfo, chapter string, link string, index int) {

}
