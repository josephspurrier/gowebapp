package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/josephspurrier/gowebapp/database"
	"github.com/josephspurrier/gowebapp/shared/mysql"
	"github.com/josephspurrier/gowebapp/shared/passhash"
	"github.com/josephspurrier/gowebapp/shared/session"
	"github.com/josephspurrier/gowebapp/shared/view"

	"github.com/josephspurrier/csrfbanana"
)

func LoginGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// If user is authenticated
	if sess.Values["id"] != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Display the view
	v := view.New(r)
	v.Name = "login"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	// Refill any form fields
	view.Repopulate([]string{"email"}, r.Form, v.Vars)
	v.Render(w)
}

func LoginPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// If user is authenticated
	if sess.Values["id"] != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Prevent brute force login attempts by not hitting MySQL and pretending like it was invalid :-)
	if sess.Values["login_attempt"] != nil && sess.Values["login_attempt"].(int) >= 5 {
		log.Println("Brute force login prevented")
		sess.AddFlash(view.Flash{"Sorry, no brute force :-)", view.FlashNotice})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"email", "password"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}

	// Form values
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Get database result
	db, _ := mysql.Instance()
	defer db.Link.Close()
	result := database.User{}
	err := db.Link.Get(&result, "SELECT id, password, status_id, first_name FROM user WHERE email = ? LIMIT 1", email)

	// Determine if password is correct
	if err == sql.ErrNoRows {
		// Log the attempt
		if sess.Values["login_attempt"] == nil {
			sess.Values["login_attempt"] = 1
		} else {
			sess.Values["login_attempt"] = sess.Values["login_attempt"].(int) + 1
		}
		sess.AddFlash(view.Flash{"Password is incorrect - Attempt: " + fmt.Sprintf("%v", sess.Values["login_attempt"]), view.FlashWarning})
		sess.Save(r, w)
	} else if err != nil {
		// Display error message
		log.Println(err)
		sess.AddFlash(view.Flash{"There was an error. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else if passhash.MatchString(result.Password, password) {
		if result.Status_id != 1 {
			// User inactive and display inactive message
			sess.AddFlash(view.Flash{"Account is inactive so login is disabled.", view.FlashNotice})
			sess.Save(r, w)
		} else {
			// Login successfully
			// Clear out all stored values in the cookie
			for k := range sess.Values {
				delete(sess.Values, k)
			}
			sess.AddFlash(view.Flash{"Login successful!", view.FlashSuccess})
			sess.Values["id"] = result.Id
			sess.Values["email"] = email
			sess.Values["first_name"] = result.First_name
			err := sess.Save(r, w)
			if err != nil {
				log.Println(err)
			}
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	} else {
		// Log the attempt
		if sess.Values["login_attempt"] == nil {
			sess.Values["login_attempt"] = 1
		} else {
			sess.Values["login_attempt"] = sess.Values["login_attempt"].(int) + 1
		}

		sess.AddFlash(view.Flash{"Password is incorrect - Attempt: " + fmt.Sprintf("%v", sess.Values["login_attempt"]), view.FlashWarning})
		sess.Save(r, w)
	}

	// Show the login page again
	LoginGET(w, r)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// If user is authenticated
	if sess.Values["id"] != nil {
		// Clear out all stored values in the cookie
		for k := range sess.Values {
			//log.Println("Deleting: ", k)
			delete(sess.Values, k)
		}
		sess.AddFlash(view.Flash{"Goodbye!", view.FlashNotice})

		// Save the cookie
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
