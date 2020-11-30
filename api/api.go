package api

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"go-solidary/config"
)

// HandlerBusiness is called to show business list
func HandlerBusiness(w http.ResponseWriter, r *http.Request, c *config.Config) {
	switch r.Method {
	case "GET":
	case "POST":
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

// UpServer start the server
func UpServer(c *config.Config) (err error) {
	http.HandleFunc("/business", func(w http.ResponseWriter, r *http.Request) {
		//HandlerUsers(w, r, c)
	})
	http.HandleFunc("/customer", func(w http.ResponseWriter, r *http.Request) {
		//HandlerFacebookAuth(w, r, c)
	})
	// http.HandleFunc("/users/{key}", HandlerOneUser)
	if c.IsHTTPS() {
		var cert, key []byte
		cert, err = c.GetCertString()
		if err != nil {
			return err
		}
		key, err = c.GetKeyString()
		if err != nil {
			return err
		}
		keyPair, err := tls.X509KeyPair(cert, key)
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{keyPair},
			// Other options
		}
		// Build a server:
		server := http.Server{
			// Other options
			TLSConfig: tlsConfig,
		}
		err = server.ListenAndServeTLS("", "")
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}
