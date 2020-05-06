package handlers

import (
	"gopkg.in/macaron.v1"
)

// IndexPage - Request handler for index page
func IndexPage(ctx *macaron.Context) string {
	return "INDEX " + ctx.Req.RequestURI
}
