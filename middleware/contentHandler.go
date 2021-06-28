package middleware

import (
	"github.com/mlzyplntsyntk/go4web/core"
	"html/template"
	"net/http"
)

var (
	Controllers []ContentController = make([]ContentController, 0)
	templates                       = template.Must(template.ParseFiles("resource/blog.html", "resource/home.html"))
)

type ContentController interface {
	SetPattern() string
	SetView() string
	SetModel() map[string]string
}

type ContentHandler struct {
}

func (t ContentHandler) RunBeforeHandled() bool {
	return false
}

func (t ContentHandler) RunAfterHandled() bool {
	return false
}

func (t ContentHandler) Pattern() string {
	return "*"
}

func (t ContentHandler) Run(w http.ResponseWriter, r *http.Request, route *core.Page) (result bool, err error) {
	for item := range Controllers {
		if Controllers[item].SetPattern() == route.Pattern {
			model := Controllers[item].SetModel()
			view := Controllers[item].SetView()
			err := templates.ExecuteTemplate(w, view+".html", model)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
	return
}

func AddController(ch ContentController) {
	Controllers = append(Controllers, ch)
}
