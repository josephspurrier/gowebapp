package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/josephspurrier/gowebapp/config"
)

// Run starts the HTTP and/or HTTPS listener
func Run(handlers http.Handler) {
	if config.Raw.Server.UseHTTP && config.Raw.Server.UseHTTPS {
		go func() {
			startHTTPS(handlers)
		}()

		startHTTP(handlers)
	} else if config.Raw.Server.UseHTTP {
		startHTTP(handlers)
	} else if config.Raw.Server.UseHTTPS {
		startHTTPS(handlers)
	} else {
		log.Println("Config file does not specify a listener to start")
	}
}

// startHTTP starts the HTTP listener
func startHTTP(handlers http.Handler) {
	fmt.Println(time.Now().Format("2006-01-02 03:04:05 PM"), "Running HTTP "+httpAddress())

	// Start the HTTP listener
	log.Fatal(http.ListenAndServe(httpAddress(), handlers))
}

// startHTTPs starts the HTTPS listener
func startHTTPS(handlers http.Handler) {
	fmt.Println(time.Now().Format("2006-01-02 03:04:05 PM"), "Running HTTPS "+httpsAddress())

	// Start the HTTPS listener
	log.Fatal(http.ListenAndServeTLS(httpsAddress(), config.Raw.Server.CertFile, config.Raw.Server.KeyFile, handlers))
}

// httpAddress returns the HTTP address
func httpAddress() string {
	return config.Raw.Server.Hostname + ":" + fmt.Sprintf("%d", config.Raw.Server.HTTPPort)
}

// httpsAddress returns the HTTPS address
func httpsAddress() string {
	return config.Raw.Server.Hostname + ":" + fmt.Sprintf("%d", config.Raw.Server.HTTPSPort)
}
