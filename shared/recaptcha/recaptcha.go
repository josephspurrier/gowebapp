package recaptcha

import (
	"html/template"
	"net/http"

	"github.com/haisum/recaptcha"
)

var (
	recap Info
)

// Info has the details for the Google reCAPTCHA
type Info struct {
	Enabled bool
	Secret  string
	SiteKey string
}

// Configure adds the settings for Google reCAPTCHA
func Configure(c Info) {
	recap = c
}

// ReadConfig returns the settings for Google reCAPTCHA
func ReadConfig() Info {
	return recap
}

// Verified returns whether the Google reCAPTCHA was verified or not
func Verified(r *http.Request) bool {
	if !recap.Enabled {
		return true
	}

	// Check the reCaptcha
	re := recaptcha.R{
		Secret: recap.Secret,
	}
	return re.Verify(*r)
}

// Plugin returns a map of functions that are usable in templates
func Plugin() template.FuncMap {
	f := make(template.FuncMap)

	f["RECAPTCHA_SITEKEY"] = func() template.HTML {
		if ReadConfig().Enabled {
			return template.HTML(ReadConfig().SiteKey)
		}

		return template.HTML("")
	}

	return f
}
