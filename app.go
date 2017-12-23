package main

import (
	"./controller"
	"./core"
	"./middleware"
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
