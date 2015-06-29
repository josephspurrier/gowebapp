package config

import (
	"encoding/json"
	"fmt"

	"github.com/josephspurrier/gowebapp/shared/mysql"
	"github.com/josephspurrier/gowebapp/shared/view"

	"github.com/gorilla/sessions"
)

// Settings is a container for the config data
var Raw = &Layout{}

// ListenAddress returns the address for ListenAndServe
func ListenAddress() string {
	return Raw.Server.Hostname + ":" + fmt.Sprintf("%d", Raw.Server.Port)
}

// Layout is the top container
type Layout struct {
	Server   `json:"Server"`
	Session  `json:"Session"`
	View     view.View            `json:"View"`
	MySQL    mysql.ConnectionInfo `json:"MySQL"`
	Template `json:"Template"`
}

// Server stores the hostname and port number
type Server struct {
	Hostname string `json:"Hostname"` // Server name
	Port     int    `json:"Port"`     // Port to listen on
}

// Session stores session level information
type Session struct {
	Options   sessions.Options `json:"Options"`   // Pulled from: http://www.gorillatoolkit.org/pkg/sessions#Options
	Name      string           `json:"Name"`      // Name for: http://www.gorillatoolkit.org/pkg/sessions#CookieStore.Get
	SecretKey string           `json:"SecretKey"` // Key for: http://www.gorillatoolkit.org/pkg/sessions#CookieStore.New
}

type Template struct {
	Root     string   `json:"Root"`
	Children []string `json:"Children"`
}

// ParseJSON unmarshals bytes to structs
func (c *Layout) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}
