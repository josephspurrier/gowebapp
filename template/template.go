package template

import "embed"

// Assets represents the embedded files.
//go:embed *.tmpl */*.tmpl
var Assets embed.FS
