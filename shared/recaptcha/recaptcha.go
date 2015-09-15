package recaptcha

import (
	"github.com/haisum/recaptcha"
	"net/http"
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
