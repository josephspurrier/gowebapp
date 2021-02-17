package main

import (
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/josephspurrier/gowebapp/app/route"
	"github.com/josephspurrier/gowebapp/app/shared/database"
	"github.com/josephspurrier/gowebapp/app/shared/email"
	"github.com/josephspurrier/gowebapp/app/shared/jsonconfig"
	"github.com/josephspurrier/gowebapp/app/shared/recaptcha"
	"github.com/josephspurrier/gowebapp/app/shared/server"
	"github.com/josephspurrier/gowebapp/app/shared/session"
	"github.com/josephspurrier/gowebapp/app/shared/view"
	"github.com/josephspurrier/gowebapp/app/shared/view/plugin"
)

// *****************************************************************************
// Application Logic
// *****************************************************************************

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)
}

func main() {
	// Load the configuration file
	jsonconfig.Load(filepath.Join("config", "config.json"), config)

	// Configure the session cookie store
	session.Configure(config.Session)

	// Connect to database
	database.Connect(config.Database)

	// Configure the Google reCAPTCHA prior to loading view plugins
	recaptcha.Configure(config.Recaptcha)

	// Setup the views
	view.Configure(config.View)
	view.LoadTemplates(config.Template.Root, config.Template.Children)
	view.LoadPlugins(
		plugin.TagHelper(config.View),
		plugin.NoEscape(),
		plugin.PrettyTime(),
		recaptcha.Plugin())

	// Start the listener
	server.Run(route.LoadHTTP(), route.LoadHTTPS(), config.Server)
}

// *****************************************************************************
// Application Settings
// *****************************************************************************

// config the settings variable
var config = &configuration{}

// configuration contains the application settings
type configuration struct {
	Database  database.Info   `json:"Database"`
	Email     email.SMTPInfo  `json:"Email"`
	Recaptcha recaptcha.Info  `json:"Recaptcha"`
	Server    server.Server   `json:"Server"`
	Session   session.Session `json:"Session"`
	Template  view.Template   `json:"Template"`
	View      view.View       `json:"View"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}
