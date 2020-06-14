package core

import (
	"database/sql"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"./types"
	"./utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	CHAPTERS_FOLDER = "ca"
	COVERS_FOLDER   = "co"
	THUMB_FOLDER    = "th"
	META_FILE       = "me.json"
)

type MangaInfo types.MangaInfo

type DownloadManager struct {
	MangaFolder   string
	DataBase      *sql.DB
	PluginManager *PluginManager
}

// Download - Download
func (m *DownloadManager) Download(mangaURL string) {

	// TODO: check if manga in db
	// if not exists starting goroutine

	manga := make(chan string)
	go m.saveManga(manga)
	manga <- mangaURL
}

func (m *DownloadManager) saveManga(ch chan string) {
	mangaURL := <-ch

	plugin := m.PluginManager.GetPluginForSite(mangaURL)
	(*plugin).GetMetaData()
	(*plugin).GetChapters()

	mangaInfo := MangaInfo((*plugin).GetData())

	m.prepareDirectories(&mangaInfo)

	m.saveMangaChapters(&mangaInfo)
	m.saveMangaMetadata(&mangaInfo)

}

func (m *DownloadManager) getMangaPath(mangaInfo *MangaInfo) string {
	mangaID := utils.HashSha1(mangaInfo.URL)
	return filepath.Join(m.MangaFolder, mangaInfo.Hostname, mangaID)
}

func (m *DownloadManager) prepareDirectories(mangaInfo *MangaInfo) {
	path := m.getMangaPath(mangaInfo)
	os.MkdirAll(path, os.ModePerm)
	os.MkdirAll(filepath.Join(path, CHAPTERS_FOLDER), os.ModePerm)
	os.MkdirAll(filepath.Join(path, COVERS_FOLDER), os.ModePerm)
	os.MkdirAll(filepath.Join(path, THUMB_FOLDER), os.ModePerm)
}

func (m *DownloadManager) saveMangaChapters(mangaInfo *MangaInfo) {
	path := m.getMangaPath(mangaInfo)
	for _, chapter := range mangaInfo.Chapters {
		chapterID := utils.HashSha1(chapter.URL)
		os.MkdirAll(filepath.Join(path, CHAPTERS_FOLDER, chapterID), os.ModePerm)
		os.MkdirAll(filepath.Join(path, THUMB_FOLDER, chapterID), os.ModePerm)

		m.saveMangaPages(mangaInfo, &chapter, path, chapterID)
		break
	}
}

func (m *DownloadManager) saveMangaPages(mangaInfo *MangaInfo, chapter *types.MangaChapter, path string, chapterID string) {

	for _, page := range chapter.Pages {
		result := getPageByURL(page.ImageURL)
		pageLocation := filepath.Join(path, CHAPTERS_FOLDER, chapterID, fmt.Sprintf("%d", page.ImageIndex))

		img, format, err := image.Decode(result.Body)
		if err != nil {
			log.Fatal(err)
		}

		fd := utils.OpenOrCreate(pageLocation)
		defer fd.Close()

		if status, err := convertImageToPNG(format, img, fd); !status {
			log.Fatal(err)
		}

	}
}

// func saveImage()

func (m *DownloadManager) saveMangaMetadata(mangaInfo *MangaInfo) {

}

func getPageByURL(url string) *http.Response {

	result, err := http.Get(string(url))

	if nil != err {
		logrus.Fatal(err)
	}

	if 200 != result.StatusCode {
		logrus.WithFields(logrus.Fields{
			"url":  url,
			"code": result.StatusCode,
		}).Fatal("Can't get page.")
	}

	return result
}

func convertImageToPNG(format string, img image.Image, fd *os.File) (bool, error) {
	switch format {
	case "jpeg":
		if err := png.Encode(fd, img); err != nil {
			return false, errors.Wrap(err, "unable to encode png")
		}

		return true, nil
	}

	return false, fmt.Errorf("unable to convert %#v to png", format)
}
