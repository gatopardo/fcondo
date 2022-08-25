
package pprofhandler

import (
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

// Handler routes the pprof pages using httprouter
func Handler(w http.ResponseWriter, r *http.Request) {

	p := context.Get(r, "params").(httprouter.Params)

	switch p.ByName("pprof") {
	case "/cmdline":
		pprof.Cmdline(w, r)
	case "/profile":
		pprof.Profile(w, r)
	case "/symbol":
		pprof.Symbol(w, r)
	default:
		pprof.Index(w, r)
	}
}
