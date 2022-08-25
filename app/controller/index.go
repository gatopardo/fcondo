package controller

import (
	"net/http"

	"github.com/gatopardo/acad/ch14/fcondo/app/model"
	"github.com/gatopardo/acad/ch14/fcondo/app/shared/view"
)

// IndexGET displays the home page
func IndexGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	session := model.Instance(r)

	if session.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "index/auth"
		v.Render(w)
	} else {
		// Display the view
		v := view.New(r)
		v.Name = "index/anon"
		v.Render(w)
		return
	}
}
