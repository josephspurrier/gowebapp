package session

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	Store *sessions.CookieStore
	Name  string
)

// Session stores session level information
type Session struct {
	Options   sessions.Options `json:"Options"`   // Pulled from: http://www.gorillatoolkit.org/pkg/sessions#Options
	Name      string           `json:"Name"`      // Name for: http://www.gorillatoolkit.org/pkg/sessions#CookieStore.Get
	SecretKey string           `json:"SecretKey"` // Key for: http://www.gorillatoolkit.org/pkg/sessions#CookieStore.New
}

// Start a session
func Start(secretKey string, options sessions.Options, name string) {
	Store = sessions.NewCookieStore([]byte(secretKey))
	Store.Options = &options
	Name = name
}

// Session returns a new session, never returns an error
func Instance(r *http.Request) *sessions.Session {
	session, _ := Store.Get(r, Name)
	return session
}
