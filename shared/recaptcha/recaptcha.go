package recaptcha

import (
	"html/template"
	"net/http"

	"github.com/haisum/recaptcha"
)

var (
	recap RecaptchaInfo
)

// RecaptchaInfo has the details for the Google reCAPTCHA
type RecaptchaInfo struct {
	Enabled bool
	Secret  string
	SiteKey string
}

// Configure adds the settings for Google reCAPTCHA
func Configure(c RecaptchaInfo) {
	recap = c
}

// ReadConfig returns the settings for Google reCAPTCHA
func ReadConfig() RecaptchaInfo {
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

// RecaptchaFuncMapp returns a map of functions that are usable in templates
func RecaptchaPlugin() template.FuncMap {
	f := make(template.FuncMap)

	f["RECAPTCHA_SITEKEY"] = func() template.HTML {
		if ReadConfig().Enabled {
			return template.HTML(ReadConfig().SiteKey)
		}

		return template.HTML("")
	}

	return f
}
