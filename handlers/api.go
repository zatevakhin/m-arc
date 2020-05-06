package handlers

import (
	"gopkg.in/macaron.v1"
)

// IndexList - Request handler for index page
func IndexList(ctx *macaron.Context) string {
	return "API - IndexList"
}
