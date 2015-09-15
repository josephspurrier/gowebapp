package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"

	"github.com/josephspurrier/gowebapp/route"
	"github.com/josephspurrier/gowebapp/shared/database"
	"github.com/josephspurrier/gowebapp/shared/jsonconfig"
	"github.com/josephspurrier/gowebapp/shared/server"
	"github.com/josephspurrier/gowebapp/shared/session"
	"github.com/josephspurrier/gowebapp/shared/view"
	"github.com/josephspurrier/gowebapp/shared/view/plugin"
)

// *****************************************************************************
// Application Logic
// *****************************************************************************

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// Load the configuration file
	jsonconfig.Load("config"+string(os.PathSeparator)+"config.json", config)

	// Configure the session cookie store
	session.Configure(config.Session)

	// Connect to database
	database.Connect(config.Database)

	// Setup the views
	view.Configure(config.View)
	view.LoadTemplates(config.Template.Root, config.Template.Children)
	view.LoadPlugins(plugin.TemplateFuncMap(config.View))

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
	Database database.Databases `json:"Database"`
	Server   server.Server      `json:"Server"`
	Session  session.Session    `json:"Session"`
	Template view.Template      `json:"Template"`
	View     view.View          `json:"View"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}
