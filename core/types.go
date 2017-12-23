package core

import (
	"net/http"
)

type Language struct {
	Name string
	Prefix string
	IsDefault bool
}

type Page struct {
	Name string
	Pattern string
	Description string
	IsDefault bool
	Template string
}

type Route struct {
	LanguagePrefix string
	Pages []Page
}

type Config struct {
	Languages []Language
	Routes []Route
}

type Handler struct {
	RunBeforeHandled bool
	RunAfterHandled bool
	Pattern string
	HandlerFunc func(http.ResponseWriter, *http.Request, *Page) (result bool, err error)
}