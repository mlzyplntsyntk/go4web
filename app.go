package main 

import (
	"./core"
	"./middleware"
)

func addHandlers() {
	core.AddHandler(middleware.ContentHandler{})
}

func main() {
	addHandlers()
	core.Start(8080, "./config.json")
}
