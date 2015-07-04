package route

import (
	"net/http"

	"github.com/josephspurrier/gowebapp/controller"
	hr "github.com/josephspurrier/gowebapp/route/middleware/httprouterwrapper"
	"github.com/josephspurrier/gowebapp/route/middleware/logrequest"
	"github.com/josephspurrier/gowebapp/route/middleware/pprofhandler"
	"github.com/josephspurrier/gowebapp/shared/session"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
)

// Load the routes and middleware
func Load() http.Handler {
	return middleware(routes())
}

// *****************************************************************************
// Routes
// *****************************************************************************

func routes() *httprouter.Router {
	r := httprouter.New()

	// Set 404 handler
	r.NotFound = http.HandlerFunc(controller.Error404)

	// Serve static files, no directory browsing
	r.GET("/static/*filepath", hr.HandlerFunc(controller.Static))

	// Home page
	r.GET("/", hr.Handler(http.HandlerFunc(controller.Index)))

	// Login
	r.GET("/login", hr.HandlerFunc(controller.LoginGET))
	r.POST("/login", hr.HandlerFunc(controller.LoginPOST))
	r.GET("/logout", hr.HandlerFunc(controller.Logout))

	// Register
	r.GET("/register", hr.HandlerFunc(controller.RegisterGET))
	r.POST("/register", hr.HandlerFunc(controller.RegisterPOST))

	// About
	r.GET("/about", hr.HandlerFunc(controller.AboutGET))

	// Enable Pprof
	r.GET("/debug/pprof/*pprof", pprofhandler.Handler)

	return r
}

// *****************************************************************************
// Middleware
// *****************************************************************************

func middleware(h http.Handler) http.Handler {
	// Prevents CSRF and Double Submits
	cs := csrfbanana.New(h, session.Store, session.Name)
	cs.FailureHandler(http.HandlerFunc(controller.InvalidToken))
	cs.ClearAfterUsage(true)
	cs.ExcludeRegexPaths([]string{"/static(.*)"})
	csrfbanana.TokenLength = 32
	csrfbanana.TokenName = "token"
	h = cs

	// Log every request
	h = logrequest.Handler(h)

	// Clear handler for Gorilla Context
	h = context.ClearHandler(h)

	return h
}
