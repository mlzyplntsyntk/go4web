package main

import (
	"github.com/mlzyplntsyntk/go4web/controller"
	"github.com/mlzyplntsyntk/go4web/core"
	"github.com/mlzyplntsyntk/go4web/middleware"
)

func addHandlers() {
	handler := middleware.ContentHandler{}
	middleware.AddController(controller.BlogController{})
	middleware.AddController(controller.IndexController{})
	core.AddHandler(handler)
}

func main() {
	addHandlers()
	core.Start(8080, "./config.json")
}
