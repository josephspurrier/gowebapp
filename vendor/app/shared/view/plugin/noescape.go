package plugin

import (
	"html/template"
)

// NoEscape returns a template.FuncMap
// * NOESCAPE prevents escaping variable
func NoEscape() template.FuncMap {
	f := make(template.FuncMap)

	f["NOESCAPE"] = func(name string) template.HTML {
		return template.HTML(name)
	}

	return f
}
