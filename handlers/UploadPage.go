package handlers

import (
	"../core"

	"gopkg.in/macaron.v1"
)

// UploadPage - Request handler for index page
func UploadPage(ctx *macaron.Context, man *core.DownloadManager) {
	ctx.Data["MangaFolder"] = man.MangaFolder

	man.Download("https://mangalib.me/futari-to-futari")

	ctx.HTML(200, "upload")
}
