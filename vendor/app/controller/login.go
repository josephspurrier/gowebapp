package controller

import (
	"fmt"
	"log"
	"net/http"

	"app/model"
	"app/shared/passhash"
	"app/shared/session"
	"app/shared/view"

	"github.com/gorilla/sessions"
	"github.com/josephspurrier/csrfbanana"
)

const (
	// Name of the session variable that tracks login attempts
	sessLoginAttempt = "login_attempt"
)

// loginAttempt increments the number of login attempts in sessions variable
func loginAttempt(sess *sessions.Session) {
	// Log the attempt
	if sess.Values[sessLoginAttempt] == nil {
		sess.Values[sessLoginAttempt] = 1
	} else {
		sess.Values[sessLoginAttempt] = sess.Values[sessLoginAttempt].(int) + 1
	}
}

// LoginGET displays the login page
func LoginGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "login/login"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	// Refill any form fields
	view.Repopulate([]string{"email"}, r.Form, v.Vars)
	v.Render(w)
}

// LoginPOST handles the login form submission
func LoginPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Prevent brute force login attempts by not hitting MySQL and pretending like it was invalid :-)
	if sess.Values[sessLoginAttempt] != nil && sess.Values[sessLoginAttempt].(int) >= 5 {
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
	result, err := model.UserByEmail(email)

	// Determine if user exists
	if err == model.ErrNoResult {
		loginAttempt(sess)
		sess.AddFlash(view.Flash{"Password is incorrect - Attempt: " + fmt.Sprintf("%v", sess.Values[sessLoginAttempt]), view.FlashWarning})
		sess.Save(r, w)
	} else if err != nil {
		// Display error message
		log.Println(err)
		sess.AddFlash(view.Flash{"There was an error. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else if passhash.MatchString(result.Password, password) {
		if result.StatusID != 1 {
			// User inactive and display inactive message
			sess.AddFlash(view.Flash{"Account is inactive so login is disabled.", view.FlashNotice})
			sess.Save(r, w)
		} else {
			// Login successfully
			session.Empty(sess)
			sess.AddFlash(view.Flash{"Login successful!", view.FlashSuccess})
			sess.Values["id"] = result.UserID()
			sess.Values["email"] = email
			sess.Values["first_name"] = result.FirstName
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	} else {
		loginAttempt(sess)
		sess.AddFlash(view.Flash{"Password is incorrect - Attempt: " + fmt.Sprintf("%v", sess.Values[sessLoginAttempt]), view.FlashWarning})
		sess.Save(r, w)
	}

	// Show the login page again
	LoginGET(w, r)
}

// LogoutGET clears the session and logs the user out
func LogoutGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// If user is authenticated
	if sess.Values["id"] != nil {
		session.Empty(sess)
		sess.AddFlash(view.Flash{"Goodbye!", view.FlashNotice})
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
