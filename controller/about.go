package controller

import (
	"net/http"

	"github.com/josephspurrier/gowebapp/shared/view"
)

// Displays the default home page
func AboutGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "about"
	v.Render(w)
}
