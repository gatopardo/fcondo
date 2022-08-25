package controller

import (
        "log"
	"net/http"
//        "fmt"
        "strings"
        "strconv"
        "time"
        "encoding/json"

	"github.com/gatopardo/acad/ch14/fcondo/app/model"
	"github.com/gatopardo/acad/ch14/fcondo/app/shared/view"
	"github.com/gatopardo/acad/ch14/fcondo/app/shared/email"
	"github.com/gorilla/sessions"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
       "github.com/julienschmidt/httprouter"
  )

  type  STotals struct{
       SCuot    int64
       SIng     int64
       SEgre    int64
       SAtra    int64
       SAmount  int64
  }

   type  ArPay struct {
	Apto    string
	Final   time.Time
	APaym   []model.CuotApt
   }

   type TotPay struct {
        Value   int64
	Fname   string
	Lname   string
	Email   string
	Codigo  string
	Final   time.Time
	APaym   []model.CuotApt
   }

type AptEstadJ struct {
        Fecha       time.Time       `db:"fecha" bson:"fecha"`
        Cuota             int64     `db:"cuota" bson:"cuota"`
        Amount            int64     `db:"amount" bson:"amount"`
        Balance           int64     `db:"balance" bson:"balance"`
}

type AptEstadL struct {
	Apt         string
        Period      time.Time
	LisEstad     []AptEstadJ
}


// ---------------------------------------------------
// MailSendGet despliega formulario para enviar correo
func MailSendGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)

	lisApts, err := model.Apts()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay apartas", view.FlashError})
             http.Redirect(w, r, "/apto/list", http.StatusFound)
         }
	v                  := view.New(r)
	v.Name              = "report/rptmail"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["LisApts"]   =  lisApts
        v.Vars["Level"]     =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
  func  getContent(r *http.Request)( tema, content string){
            tim        := time.Now()
            fec        := tim.Format(layout)
            hour       := tim.Format(timeLayout)
            stm        := "Fecha "+fec +" Hora : " + hour
	    tema       = r.FormValue("tema")
	    content    = stm +"\n" + r.FormValue("content")
       return
  }
// ---------------------------------------------------
   func sendPost(sess *sessions.Session, lisPers []model.Person, tema, content string){
             for _,person := range lisPers{
                 to := person.Email
                err := email.SendEmail(to, tema,content);
                if err != nil {
                   sess.AddFlash(view.Flash{"Error enviando ", view.FlashError})
		   log.Println("Error Enviando", err)
                }
	    }
   }
// ---------------------------------------------------
// MailSendPOST procesa la forma enviada con contenido
func MailSendPOST(w http.ResponseWriter, r *http.Request) {
	sess          := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
	    var lisPers []model.Person
	    var pers      model.Person
            tema, content := getContent(r)
	    lisApts, err := model.Apts()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay apartas", view.FlashError})
                 http.Redirect(w, r, "/apto/list", http.StatusFound)
             }
             for i, _ := range lisApts {
		 sin    := strconv.Itoa(i)
		 sid    := r.FormValue(sin)
		 if len(sid) > 0{
                    aid,_ := atoi32(sid)
		    pers,_, err = model.EmailByAptId(aid)
		    if err != nil {
	                   log.Println(err)
                            sess.AddFlash(view.Flash{"No user", view.FlashError})
		    }
		    lisPers  =  append(lisPers, pers)
		 }
	     }
	     go sendPost(sess,lisPers, tema, content)
             sess.AddFlash(view.Flash{"Envio exitoso!", view.FlashSuccess})
        }
	http.Redirect(w, r, "/user/list", http.StatusFound)
 }
// ---------------------------------------------------
// RptAptGET reporte estado de apto
func RptAptGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriods, err := model.Periods()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay periodos", view.FlashError})
         }
	v                     := view.New(r)
	v.Name                 = "report/rptper"
	v.Vars["token"]        = csrfbanana.Token(w, r, sess)
        v.Vars["LisPeriods"]   = lisPeriods
        v.Vars["Level"]        =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
// RptAptPOST reporte estado de apto
func RptAptPOST(w http.ResponseWriter, r *http.Request) {
	var pers model.Person
	var apt model.Aparta
	var peridi model.Periodo
	var peridf model.Periodo
	var err error
	sess := model.Instance(r)
        uid, ok       := sess.Values["id"].(uint32)
	if ! ok {
             log.Println("No uint32 value in session")
	}
        sPeridf    :=  r.FormValue("idf")
	fperid,_   := atoi32(sPeridf)
        sPeridi    :=  r.FormValue("idi")
	iperid,_   := atoi32(sPeridi)
	action    := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
             pers, apt, err = model.ApartaByUserId(uid)
	     if err != nil {
	        log.Println(err)
                sess.AddFlash(view.Flash{"No apto", view.FlashError})
	     }
             peridf.Id = fperid
             err := (&peridf).PeriodById()
             if err != nil {
	         log.Println(err)
                 sess.AddFlash(view.Flash{"No hay periodo", view.FlashError})
             }
             peridi.Id = iperid
             err = (&peridi).PeriodById()
             if err != nil {
	         log.Println(err)
                 sess.AddFlash(view.Flash{"No hay periodo", view.FlashError})
             }
	     lisPaym, _            := model.Payments(apt.Id, peridf.Inicio, peridi.Inicio)
	     lon                   :=  len(lisPaym)
	     value                 := lisPaym[lon - 1].Balance
             v                     := view.New(r)
             v.Name                 = "report/rptapto"
	     v.Vars["token"]        = csrfbanana.Token(w, r, sess)
             v.Vars["Apt"]          = apt
             v.Vars["Pers"]         = pers
             v.Vars["Perid"]        = peridf
	     v.Vars["Valor"]        = value
             v.Vars["LisPaym"]      = lisPaym
             v.Vars["Level"]        =  sess.Values["level"]
	     v.Render(w)
        }else{
	  http.Redirect(w, r, "/cuota/list", http.StatusFound)
	 }
 }

// ---------------------------------------------------
// ---------------------------------------------------
// RptLisAptGet reporte estado de aptos
func RptLisAptGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriods, err := model.Periods()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay periodos", view.FlashError})
             http.Redirect(w, r, "/apto/list", http.StatusFound)
         }
	var lisApts []model.Aparta
        lisApts, err = model.Apts()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay apartas", view.FlashError})
             http.Redirect(w, r, "/apto/list", http.StatusFound)
         }
	v                     := view.New(r)
	v.Name                 = "report/rptlsapt"
	v.Vars["token"]        = csrfbanana.Token(w, r, sess)
        v.Vars["LisPeriods"]   = lisPeriods
        v.Vars["LisApts"]     = lisApts
        v.Vars["Level"]        =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
// RptLisAptPOST reporte estado de aptos
func RptLisAptPOST(w http.ResponseWriter, r *http.Request) {
	var pers model.Person
	var apt model.Aparta
	var peridi model.Periodo
	var peridf model.Periodo
	var  lisPay []TotPay
	var  aPay   TotPay
	sess := model.Instance(r)
	action    := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            sPeridf    :=   r.FormValue("idf")
            fperid,_   :=   atoi32(sPeridf)
            sPeridi    :=   r.FormValue("idi")
            iperid,_   :=   atoi32(sPeridi)
            lisApts, err := model.Apts()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay apartas", view.FlashError})
                 http.Redirect(w, r, "/apto/list", http.StatusFound)
             }
             peridf.Id = fperid
             err  = (&peridf).PeriodById()
             if err != nil {
	         log.Println(err)
                 sess.AddFlash(view.Flash{"No periodo Final", view.FlashError})
              }
             peridi.Id = iperid
             err = (&peridi).PeriodById()
            if err != nil {
	        log.Println(err)
                sess.AddFlash(view.Flash{"No periodo Inicial", view.FlashError})
             }

             for i, _ := range lisApts {
		 sin    := strconv.Itoa(i)
		 sid    := r.FormValue(sin)
		 if len(sid) > 0{
                    aid,_ := atoi32(sid)
                    pers, apt, err = model.EmailByAptId(aid)
                    if err != nil {
	                 log.Println(err)
                         sess.AddFlash(view.Flash{"No apto", view.FlashError})
	           }
	            aPay.APaym, _ = model.Payments(aid, peridf.Inicio, peridi.Inicio)
		    if err != nil {
	                   log.Println(err)
                            sess.AddFlash(view.Flash{"No Payments", view.FlashError})
		    }else{
	               value      :=  aPay.APaym[len(aPay.APaym) - 1].Balance
		       aPay.Value  =  value
		       aPay.Final  =  peridf.Final
		       aPay.Lname  =  pers.Lname
		       aPay.Fname  =  pers.Fname
		       aPay.Email  =  pers.Email
		       aPay.Codigo =  apt.Codigo
		       lisPay      =  append(lisPay, aPay)
	            }
		 }
	     }

             v                     := view.New(r)
             v.Name                 = "report/rptlsaptper"
	     v.Vars["token"]        = csrfbanana.Token(w, r, sess)
             v.Vars["LisPay"]      = lisPay
             v.Vars["Level"]        =  sess.Values["level"]
	     v.Render(w)
         }else{
	  http.Redirect(w, r, "/cuota/list", http.StatusFound)
	 }
 }
// ---------------------------------------------------
// JRptCondGET reporte estado de condo
func JRptCondGET(w http.ResponseWriter, r *http.Request) {
	var periodo model.Periodo
        var lisAmt  []model.AmtCond
        var js []byte
        var params httprouter.Params
	sess := model.Instance(r)
        params      = context.Get(r, "params").(httprouter.Params)
	sfec       :=  params.ByName("fec")[:10]
	dtfec,err  :=  time.Parse(layout, sfec)
        if err != nil {
	        log.Println(err)
	}else{
//        dtfec       =  time.Date(dtfec.Year(), dtfec.Month(),dtfec.Day(), 0, 0, 0, 0, time.Local)
        err         = (&periodo).PeriodByFec(dtfec)
        if err     != nil {
	        log.Println(err)
        }else{
          lisAmt, err           = model.Amounts( periodo.Id )
          if err != nil {
            sess.AddFlash(view.Flash{"No lista pagos ", view.FlashError})
            log.Println(err)
          }else{
            js, err =  json.Marshal(lisAmt)
            if err == nil{
//               fmt.Println(" json " + string(js))
               w.Header().Set("Content-Type", "application/json")
               w.Write(js)
	       return
            }
           }
          }
          }
          log.Println("JRptCondGet  ", err)
          http.Error(w, err.Error(), http.StatusInternalServerError)
          return

 }
// ---------------------------------------------------
// RptCondGet reporte estado de condo
func RptCondGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriods, err := model.Periods()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay periodos", view.FlashError})
         }
	v                     := view.New(r)
	v.Name                 = "report/condper"
	v.Vars["token"]        = csrfbanana.Token(w, r, sess)
        v.Vars["LisPeriods"]   = lisPeriods
        v.Vars["Level"]        =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
// RptCondPOST reporte estado de condominio
func RptCondPOST(w http.ResponseWriter, r *http.Request) {
	var periodo model.Periodo
        var lisAmt  []model.AmtCond
	sess := model.Instance(r)
        sPerid    :=  r.FormValue("id")
	perid,_   := atoi32(sPerid)
	action    := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            periodo.Id = perid
            err := (&periodo).PeriodById()
            if err != nil {
	        log.Println(err)
                sess.AddFlash(view.Flash{"No periodo", view.FlashError})
             }

	    lisAmt, err           = model.Amounts( perid )
            if err != nil {
	        log.Println(err)
                sess.AddFlash(view.Flash{"No hay montos pagados", view.FlashError})
             }
             v                     := view.New(r)
             v.Name                 = "report/condapto"
	     v.Vars["token"]        =  csrfbanana.Token(w, r, sess)
             v.Vars["LisAmt"]       =  lisAmt
             v.Vars["Per"]          =  periodo
             v.Vars["Level"]        =  sess.Values["level"]
	     v.Render(w)
         }else{
	  http.Redirect(w, r, "/cuota/list", http.StatusFound)
	 }
 }
// ---------------------------------------------------
// RptAllCondGet reporte estado de apto
func RptAllCondGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriods, err := model.Periods()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay periodos", view.FlashError})
         }
	v                     := view.New(r)
	v.Name                 = "report/condpertot"
	v.Vars["token"]        = csrfbanana.Token(w, r, sess)
        v.Vars["LisPeriods"]   = lisPeriods
        v.Vars["Level"]        =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
// RptAllCondPOST reporte estado de apto
func RptAllCondPOST(w http.ResponseWriter, r *http.Request) {
	var periodo model.Periodo
	var lisAmt []model.AmtCond
	var lisCuot []model.CuotaN
	var lisEgre []model.EgresoN
	var lisIngre []model.IngresoN

	sess := model.Instance(r)
        sPerid    :=  r.FormValue("id")
	perid,_   := atoi32(sPerid)
	action    := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            periodo.Id = perid
            err := (&periodo).PeriodById()
            if err != nil {
	        log.Println(err)
                sess.AddFlash(view.Flash{"No periodo", view.FlashError})
             }
	     lisCuot, lisIngre, lisEgre,lisAmt, err  = model.MoneyFlow( perid)
	     if err != nil {
	     }
             var sum,sa int64
	     var stot STotals
	     for _, value := range lisCuot{
		     sum +=  value.Amount
	     }
             stot.SCuot  = sum
	     sum = 0
	     for _, value := range lisIngre{
		     sum +=  value.Amount
	     }
             stot.SIng  = sum
	     sum = 0
	     for _, value := range lisEgre{
		     sum +=  value.Amount
	     }
             stot.SEgre  = sum
	     sum = 0
	     for _, value := range lisAmt{
		     sum +=  value.Amount
		     sa  +=  value.Atraso
	     }
             stot.SAmount  = sum
             stot.SAtra    = sa

             v                     := view.New(r)
             v.Name                 = "report/condtot"
	     v.Vars["token"]        =  csrfbanana.Token(w, r, sess)
             v.Vars["LisCuot"]      =  lisCuot
             v.Vars["LisIngre"]     =  lisIngre
             v.Vars["LisEgre"]      =  lisEgre
             v.Vars["LisAmt"]       =  lisAmt
	     v.Vars["STot"]         =  stot
             v.Vars["Per"]          =  periodo
             v.Vars["Level"]        =  sess.Values["level"]
	     v.Render(w)
         }else{
	  http.Redirect(w, r, "/cuota/list", http.StatusFound)
	 }
 }
// ---------------------------------------------------
// ---------------------------------------------------
  func procesaApt(lsAptTots []model.AmtAptTot)(pos int, scuot, smonto int64){
         var  sdif int64
	 scuot  = int64(0)
	 smonto = int64(0)
	 sdif   =  int64(0)
	 sfec  := time.Date(1,1,1,0,0,0,0,time.UTC)
//         sf    := dFini.Format(formato)
         pos = 0
	 j := 0
	 k := 0
         for i, amt := range lsAptTots {
             k = i
	     smonto = smonto +  amt.Monto
	     if sfec != amt.Inicio {
	          scuot = scuot + amt.Cuota
	     }
	     sdif =  scuot - smonto
	     if sdif <= int64(0){
		     j = pos
		     pos = i
	     }
	     lsAptTots[i].Dif = sdif
	     sfec = amt.Inicio
	 }
         if  (k - pos) <  5{
             pos = j
	     if  (j  - pos) < 5 && (k - 10) > 0{
                pos = k - 10
	     }
	 }
        return
  }

// ---------------------------------------------------
// getAptEstad get apt State up to a period
  func getAptEstad(apt string,  dFini time.Time)(aptEstad []AptEstadJ, err error){
	  var aptDet AptEstadJ
	  dIni,_ := time.Parse("001-01-01" , formato)
        amtTots, err := model.AptDetails(apt, dIni , dFini )
	if err != nil {
		log.Println(err)
		return
	}
        p,_,_ := procesaApt(amtTots)
	lon :=  len(amtTots)
	for i  := p; i < lon; i++ {
		aptDet = AptEstadJ{ Fecha:   amtTots[i].Fecha,
                                          Cuota:   amtTots[i].Cuota,
                                          Amount:  amtTots[i].Monto,
                                          Balance: amtTots[i].Dif,
	                                }
	   aptEstad  = append(aptEstad, aptDet)
         }
      return
  }
// ---------------------------------------------------
// JRptCondDetGET reporte estado de condo en detalles
func JRptCondDetGET(w http.ResponseWriter, r *http.Request) {
	var periodo model.Periodo
	var apts   []model.Apt
        var lisEstad     []AptEstadL
	var aptEstad       AptEstadL
        var lisAptEstad  []AptEstadJ
        var js []byte
	var params httprouter.Params
	sess := model.Instance(r)
        params      = context.Get(r, "params").(httprouter.Params)
	sfec       :=  params.ByName("fec")[:10]
	dtfec,err  :=  time.Parse(layout, sfec)
        err         = (&periodo).PeriodByFec(dtfec)
        if err     != nil {
	        log.Println("JRptCondDetGET",err)
	        log.Fatalln("JRptCondDetGET",err)
        }
        apts, err = model.LisApts() // Apts()
        if err     != nil {
            sess.AddFlash(view.Flash{"No lista Apts ", view.FlashError})
	        log.Println("JRptCondDetGET",err)
	        log.Fatalln("JRptCondDetGET",err)
        }
	for _, apt := range apts {
             lisAptEstad, err       = getAptEstad(string(apt), dtfec  )
              if err != nil {
                log.Println("JRptCondDetGET",err)
              }else{
		      aptEstad = AptEstadL{ Period:  dtfec,
		                            Apt:     string(apt),
		                            LisEstad: lisAptEstad,
		                                }
		  lisEstad = append(lisEstad, aptEstad)
             }
            }
            js, err =  json.Marshal(lisEstad)
            if err == nil{
               w.Header().Set("Content-Type", "application/json")
               w.Write(js)
	       return
            }
          log.Println("JRptCondGET  ", err)
          http.Error(w, err.Error(), http.StatusInternalServerError)
          return
 }
// ---------------------------------------------------
