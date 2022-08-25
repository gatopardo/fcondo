package logrequest

import (
	"fmt"
	"net/http"
	"time"
        "log"
)

 // Handler logueare la petición HTTP
func Handler(next http.Handler, flog * log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                time.Now()
		flog.Println( r.RemoteAddr, r.Method, r.URL)
		fmt.Println(time.Now().Format("2006-01-02 03:04:05 PM"), r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
