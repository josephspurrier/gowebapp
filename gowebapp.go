package main

import (
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/josephspurrier/gowebapp/config"
	"github.com/josephspurrier/gowebapp/controller"
	hr "github.com/josephspurrier/gowebapp/middleware/httprouterwrapper"
	"github.com/josephspurrier/gowebapp/middleware/logrequest"
	"github.com/josephspurrier/gowebapp/middleware/pprofhandler"
	"github.com/josephspurrier/gowebapp/plugin"
	"github.com/josephspurrier/gowebapp/shared/jsonconfig"
	"github.com/josephspurrier/gowebapp/shared/mysql"
	"github.com/josephspurrier/gowebapp/shared/session"
	"github.com/josephspurrier/gowebapp/shared/view"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
)

// *****************************************************************************
// Main
// *****************************************************************************

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)
	
	// Use all cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// Load the configuration file
	jsonconfig.Load("config"+string(os.PathSeparator)+"config.json", config.Raw)

	// Start the session
	session.Start(config.Raw.Session.SecretKey, config.Raw.Session.Options,
		config.Raw.Session.Name)

	// Connect to MySQL
	mysql.Config(config.Raw.MySQL)

	// Setup the views
	view.Config(config.Raw.View)
	view.LoadTemplates(config.Raw.Template.Root, config.Raw.Template.Children)
	view.LoadPlugins(plugin.TemplateFuncMap())

	// Start the HTTP listener
	log.Fatal(http.ListenAndServe(config.ListenAddress(), handlers()))
}

// *****************************************************************************
// Routing
// *****************************************************************************

func router() *httprouter.Router {
	r := httprouter.New()

	// Set 404 handler
	r.NotFound = http.HandlerFunc(controller.Error404)

	// Serve static files, no directory browsing
	r.GET("/static/*filepath", hr.HandlerFunc(controller.Static))

	// Home page
	r.GET("/", hr.Handler(http.HandlerFunc(controller.Index)))
	//r.GET("/", hr.HandlerFunc(controller.Index))

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

func handlers() http.Handler {
	var h http.Handler

	// Route to pages
	h = router()

	// Prevents CSRF and Double Submits
	cs := csrfbanana.New(h, session.Store, config.Raw.Session.Name)
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
