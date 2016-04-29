package controller

import (
	"fmt"
	"log"
	"net/http"

	"app/model"
	"app/shared/session"
	"app/shared/view"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
)

// NotepadReadGET displays the notes in the notepad
func NotepadReadGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	userID := fmt.Sprintf("%s", sess.Values["id"])

	notes, err := model.NotesByUserID(userID)
	if err != nil {
		log.Println(err)
		notes = []model.Note{}
	}

	// Display the view
	v := view.New(r)
	v.Name = "notepad/read"
	v.Vars["first_name"] = sess.Values["first_name"]
	v.Vars["notes"] = notes
	v.Render(w)
}

// NotepadCreateGET displays the note creation page
func NotepadCreateGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "notepad/create"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
}

// NotepadCreatePOST handles the note creation form submission
func NotepadCreatePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"note"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		NotepadCreateGET(w, r)
		return
	}

	// Get form values
	content := r.FormValue("note")

	userID := fmt.Sprintf("%s", sess.Values["id"])

	// Get database result
	err := model.NoteCreate(content, userID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Note added!", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/notepad", http.StatusFound)
		return
	}

	// Display the same page
	NotepadCreateGET(w, r)
}

// NotepadUpdateGET displays the note update page
func NotepadUpdateGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get the note id
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	noteID := params.ByName("id")

	userID := fmt.Sprintf("%s", sess.Values["id"])

	// Get the note
	note, err := model.NoteByID(userID, noteID)
	if err != nil { // If the note doesn't exist
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/notepad", http.StatusFound)
		return
	}

	// Display the view
	v := view.New(r)
	v.Name = "notepad/update"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["note"] = note.Content
	v.Render(w)
}

// NotepadUpdatePOST handles the note update form submission
func NotepadUpdatePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"note"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		NotepadUpdateGET(w, r)
		return
	}

	// Get form values
	content := r.FormValue("note")

	userID := fmt.Sprintf("%s", sess.Values["id"])

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	noteID := params.ByName("id")

	// Get database result
	err := model.NoteUpdate(content, userID, noteID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Note updated!", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/notepad", http.StatusFound)
		return
	}

	// Display the same page
	NotepadUpdateGET(w, r)
}

// NotepadDeleteGET handles the note deletion
func NotepadDeleteGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	userID := fmt.Sprintf("%s", sess.Values["id"])

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	noteID := params.ByName("id")

	// Get database result
	err := model.NoteDelete(userID, noteID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Note deleted!", view.FlashSuccess})
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/notepad", http.StatusFound)
	return
}
