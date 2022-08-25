package route

import (
	"net/http"
        "log"

	"github.com/gatopardo/acad/ch14/fcondo/app/controller"
	"github.com/gatopardo/acad/ch14/fcondo/app/route/middleware/acl"
	 hr  "github.com/gatopardo/acad/ch14/fcondo/app/route/middleware/httprouterwrapper"
	"github.com/gatopardo/acad/ch14/fcondo/app/route/middleware/logrequest"
	"github.com/gatopardo/acad/ch14/fcondo/app/route/middleware/pprofhandler"
	"github.com/gatopardo/acad/ch14/fcondo/app/model"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

  var Flogger  *log.Logger

// Load returns the routes and middleware
func Load() http.Handler {
           Flogger.Println("HTTP routes Load")
	return middleware(routes())
}

// LoadHTTPS returns the HTTP routes and middleware
func LoadHTTPS() http.Handler {
           Flogger.Println("HTTPS routes LoadHTTPS")
	return middleware(routes())
}

// LoadHTTP returns the HTTPS routes and middleware
func LoadHTTP() http.Handler {
           Flogger.Println("HTTPS routes LoadHTTP")
	return middleware(routes())

	// Uncomment this and comment out the line above to always redirect to HTTPS
	//return http.HandlerFunc(redirectToHTTPS)
}

// Optional method to make it easy to redirect from HTTP to HTTPS
func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
           Flogger.Println("HTTP redirect")
	http.Redirect(w, req, "https://"+req.Host, http.StatusMovedPermanently)
}

// *****************************************************************************
// Routes
// *****************************************************************************

func routes() *httprouter.Router {
	r := httprouter.New()
	// Set 404 handler
	r.NotFound = alice.
		New().
		ThenFunc(controller.Error404)

	// Serve static files, no directory browsing
	r.GET("/static/*filepath", hr.Handler(alice.
		New().
		ThenFunc(controller.Static)))

	// Home page
	r.GET("/", hr.Handler(alice.
		New().
		ThenFunc(controller.IndexGET)))

	// Login
	r.GET("/login", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.LoginGET)))
	r.POST("/login", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.LoginPOST)))
        r.GET("/jlogin/:cuenta/:password", hr.Handler(alice.
                New(acl.DisallowAuth).
                ThenFunc(controller.JLoginGET)))
//        r.POST("/jlogin", hr.Handler(alice.
//                New(acl.DisallowAuth).
//                ThenFunc(controller.JLoginPOST)))
	r.GET("/logout", hr.Handler(alice.
		New().
		ThenFunc(controller.LogoutGET)))

// Register
	r.GET("/user/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RegisterGET)))
	r.POST("/user/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RegisterPOST)))
//          Register update
	r.GET("/user/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RegisUpGET)))
	r.POST("/user/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RegisUpPOST)))
//          List
	r.GET("/user/list", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RegisLisGET)))
//           delete 
        r.GET("/user/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
                ThenFunc(controller.RegisDelGET)))
        r.POST("/user/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
                ThenFunc(controller.RegisDelPOST)))

// Apartamento
	r.GET("/apto/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AptGET)))
	r.POST("/apto/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AptPOST)))
	//  update
	r.GET("/apto/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AptUpGET)))
	r.POST("/apto/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AptUpPOST)))
	r.GET("/japt/:fec1/:fec2/:id", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.JAptGET)))
//          List
	r.GET("/apto/list/", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AptLisGET)))
//          Delete
	r.GET("/apto/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AptDeleteGET)))
         r.POST("/apto/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
                 ThenFunc(controller.AptDeletePOST)))

// Categorias
	r.GET("/categoria/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.TipoGET)))
	r.POST("/categoria/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.TipoPOST)))
//	//  update
	r.GET("/categoria/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.TipoUpGET)))
	r.POST("/categoria/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.TipoUpPOST)))
////          List
	r.GET("/categoria/list", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.TipoLisGET)))
////          Delete
	r.GET("/categoria/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.TipoDeleteGET)))
	r.POST("/categoria/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.TipoDeletePOST)))

// Periodos
	r.GET("/period/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.PeriodGET)))
	r.POST("/period/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.PeriodPOST)))
//	//  update
	r.GET("/period/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.PeriodUpGET)))
	r.POST("/period/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.PeriodUpPOST)))
////          List
	r.GET("/period/list", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.PeriodLisGET)))
////          Delete
	r.GET("/period/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.PeriodDeleteGET)))
	r.POST("/period/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.PeriodDeletePOST)))

// Balances
        r.GET("/jbalance/:fec", hr.Handler(alice.
		New(acl.DisallowAnon).
//		New(acl.DisallowAuth).
		ThenFunc(controller.JBalanGET)))
	r.GET("/balance/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BalanGET)))
	r.POST("/balance/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BalanPOST)))
//	//  update
	r.GET("/balance/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BalanUpGET)))
	r.POST("/balance/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BalanUpPOST)))
////          List
	r.GET("/balance/list", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BalanLisGET)))
////          Delete
	r.GET("/balance/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BalanDeleteGET)))
	r.POST("/balance/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BalanDeletePOST)))

// Cuotas
	r.GET("/cuota/periodo/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CuotPerGET)))
	r.POST("/cuota/periodo/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CuotPerPOST)))
	r.POST("/cuota/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CuotRegPOST)))
//	//  update
	r.GET("/cuota/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CuotUpGET)))
	r.POST("/cuota/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CuotUpPOST)))
////          List
	r.GET("/jcuot/:fec", hr.Handler(alice.
//		New(acl.DisallowAuth).
		New(acl.DisallowAnon).
		ThenFunc(controller.JCuotGET)))
//	r.POST("/jcuot", hr.Handler(alice.
//		New(acl.DisallowAnon).
//		ThenFunc(controller.JCuotPOST)))
	r.GET("/cuota/list", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CuotLisGET)))
	r.POST("/cuota/list", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CuotLisPOST)))
////          Delete
	r.GET("/cuota/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CuotDeleteGET)))
	r.POST("/cuota/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CuotDeletePOST)))

// Egresos
	r.GET("/jegre/:fec", hr.Handler(alice.
//		New(acl.DisallowAuth).
		New(acl.DisallowAnon).
		ThenFunc(controller.JEgreGET)))
	r.GET("/egreso/periodo/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EgrePerGET)))
	r.POST("/egreso/periodo/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EgrePerPOST)))
	r.POST("/egreso/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EgreRegPOST)))
//	//  update
	r.GET("/egreso/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EgreUpGET)))
	r.POST("/egreso/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EgreUpPOST)))
////          List
	r.GET("/egreso/list", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EgreLisGET)))
	r.POST("/egreso/list", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EgreLisPOST)))
////          Delete
	r.GET("/egreso/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EgreDeleteGET)))
	r.POST("/egreso/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EgreDeletePOST)))

// Ingresos
	r.GET("/jingre/:fec", hr.Handler(alice.
//		New(acl.DisallowAuth).
		New(acl.DisallowAnon).
		ThenFunc(controller.JIngreGET)))
	r.GET("/ingreso/periodo/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.IngrePerGET)))
	r.POST("/ingreso/periodo/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.IngrePerPOST)))
	r.POST("/ingreso/register", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.IngreRegPOST)))
//	//  update
	r.GET("/ingreso/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.IngreUpGET)))
	r.POST("/ingreso/update/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.IngreUpPOST)))
////          List
	r.GET("/ingreso/list", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.IngreLisGET)))
	r.POST("/ingreso/list", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.IngreLisPOST)))
////          Delete
	r.GET("/ingreso/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.IngreDeleteGET)))
	r.POST("/ingreso/delete/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.IngreDeletePOST)))

// Reporte
	r.GET("/email", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.MailSendGET)))
	r.POST("/email", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.MailSendPOST)))
       r.GET("/report/rptapto", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RptAptGET)))
	r.POST("/report/rptapto", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RptAptPOST)))
       r.GET("/report/rptlisapto", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RptLisAptGET)))
	r.POST("/report/rptlisapto", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RptLisAptPOST)))
       r.GET("/report/rptcondo", hr.Handler(alice.
                New(acl.DisallowAnon).
		ThenFunc(controller.RptCondGET)))
        r.GET("/jcondo/:fec", hr.Handler(alice.
                New(acl.DisallowAuth).
		ThenFunc(controller.JCondoGET)))
	r.POST("/report/rptcondo", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RptCondPOST)))
	r.GET("/report/jrptcondodet/:fec", hr.Handler(alice.
		New(acl.DisallowAnon).
//		New(acl.DisallowAuth).
		ThenFunc(controller.JRptCondDetGET)))
	r.GET("/report/jrptcondo/:fec", hr.Handler(alice.
		New(acl.DisallowAnon).
//		New(acl.DisallowAuth).
		ThenFunc(controller.JRptCondGET)))
       r.GET("/report/rptallcondo", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RptAllCondGET)))
	r.POST("/report/rptallcondo", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.RptAllCondPOST)))

//              New(acl.DisallowAnon).


	// About
	r.GET("/about", hr.Handler(alice.
		New().
		ThenFunc(controller.AboutGET)))

	// Enable Pprof
	r.GET("/debug/pprof/*pprof", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(pprofhandler.Handler)))

	return r
}

// *****************************************************************************
// Middleware
// *****************************************************************************

func middleware(h http.Handler) http.Handler {
	// Prevents CSRF and Double Submits
	cs := csrfbanana.New(h, model.Store, model.Name)
	cs.FailureHandler(http.HandlerFunc(controller.InvalidToken))
	cs.ClearAfterUsage(true)
	cs.ExcludeRegexPaths([]string{"/static(.*)"})
	csrfbanana.TokenLength = 32
	csrfbanana.TokenName = "token"
	csrfbanana.SingleToken = false
	h = cs

      Flogger.Println("middleware", model.Name)

	// Log every request:1
	h = logrequest.Handler(h, Flogger)

	// Clear handler for Gorilla Context
	h = context.ClearHandler(h)

	return h
}
