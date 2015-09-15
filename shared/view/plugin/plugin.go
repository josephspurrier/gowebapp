package plugin

import (
	"html/template"
	"log"

	"github.com/josephspurrier/gowebapp/shared/view"
)

// TemplateFuncMap returns a map of functions that are usable in templates
func TemplateFuncMap(v view.View) template.FuncMap {
	f := make(template.FuncMap)

	f["JS"] = func(s string) template.HTML {
		path, err := v.AssetTimePath(s)

		if err != nil {
			log.Println("JS Error:", err)
			return template.HTML("<!-- JS Error: " + s + " -->")
		}

		return template.HTML(`<script type="text/javascript" src="` + path + `"></script>`)
	}

	f["CSS"] = func(s string) template.HTML {
		path, err := v.AssetTimePath(s)

		if err != nil {
			log.Println("CSS Error:", err)
			return template.HTML("<!-- CSS Error: " + s + " -->")
		}

		return template.HTML(`<link rel="stylesheet" type="text/css" href="` + path + `" />`)
	}

	f["NOESCAPE"] = func(name string) template.HTML {
		return template.HTML(name)
	}

	f["LINK"] = func(path, name string) template.HTML {
		return template.HTML(`<a href="` + v.PrependBaseURI(path) + `">` + name + `</a>`)
	}

	return f
}
