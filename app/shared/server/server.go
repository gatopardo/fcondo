package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
        "os"
        "strconv"

        "github.com/gatopardo/acad/ch14/fcondo/app/route"
)

// Server stores the hostname and port number
type Server struct {
        Remote    bool  `json:"Remote"`   // Server origin
        Origin    string `json:"Origin"`   // Server origin
	Hostname  string `json:"Hostname"`  // Server name
	UseHTTP   bool   `json:"UseHTTP"`   // Listen on HTTP
	UseHTTPS  bool   `json:"UseHTTPS"`  // Listen on HTTPS
	HTTPPort  int    `json:"HTTPPort"`  // HTTP port
	HTTPSPort int    `json:"HTTPSPort"` // HTTPS port
	CertFile  string `json:"CertFile"`  // HTTPS certificate
	KeyFile   string `json:"KeyFile"`   // HTTPS private key
}

// Run starts the HTTP and/or HTTPS listener
func Run(httpHandlers http.Handler, httpsHandlers http.Handler, s Server) {
	 fmt.Println("Server al inicio ", s.Origin )
        if  s.Remote   {
              sport := os.Getenv("PORT")
              iport, _ :=  strconv.Atoi(sport)
              s.HTTPPort = iport
              s.HTTPSPort = iport
              s.Hostname =  ""
              route.Flogger.Println(httpsAddress(s))
         }else{
              route.Flogger.Println(httpAddress(s))
	 }
	if s.UseHTTP && s.UseHTTPS {
		go func() {
			startHTTPS(httpsHandlers, s)
		}()

		startHTTP(httpHandlers, s)
	} else if s.UseHTTP {
		startHTTP(httpHandlers, s)
	} else if s.UseHTTPS {
		startHTTPS(httpsHandlers, s)
	} else {
		log.Println("Fichero Config no specifica servidor para iniciar")
	}
}

// startHTTP starts the HTTP listener
func startHTTP(handlers http.Handler, s Server) {
	fmt.Println(time.Now().Format("2006-01-02 03:04:05 PM"), "Running HTTP "+httpAddress(s))

	// Start the HTTP listener
	log.Fatal(http.ListenAndServe(httpAddress(s), handlers))
}

// startHTTPs starts the HTTPS listener
func startHTTPS(handlers http.Handler, s Server) {
	fmt.Println(time.Now().Format("2006-01-02 03:04:05 PM"), "Running HTTPS "+httpsAddress(s))

	// Start the HTTPS listener
	if s.Remote {
             log.Fatal(http.ListenAndServe(httpsAddress(s), handlers))
        } else {
             log.Fatal(http.ListenAndServeTLS(httpsAddress(s), s.CertFile, s.KeyFile, handlers))
        }
}

// httpAddress returns the HTTP address
func httpAddress(s Server) string {
	return s.Hostname + ":" + fmt.Sprintf("%d", s.HTTPPort)
}

// httpsAddress returns the HTTPS address
func httpsAddress(s Server) string {
	return s.Hostname + ":" + fmt.Sprintf("%d", s.HTTPSPort)
}
