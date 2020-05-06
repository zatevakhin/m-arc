package handlers

import (
	"gopkg.in/macaron.v1"
)

// LoginPage - Request handler for index page
func LoginPage(ctx *macaron.Context) string {
	return "API - LoginPage"
}

// APIUserLogin - sss
func APIUserLogin(ctx *macaron.Context) string {
	return "APIUserLogin"
}
