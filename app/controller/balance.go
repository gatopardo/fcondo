package controller

import (
	"log"
	"net/http"
        "fmt"
        "time"
        "strings"
        "encoding/json"

	"github.com/gatopardo/acad/ch14/fcondo/app/model"
	"github.com/gatopardo/acad/ch14/fcondo/app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )
// ---------------------------------------------------
// JBalanGET obtener  balances
func JBalanGET(w http.ResponseWriter, r *http.Request) {
//            fmt.Println("JBalanGET 0"  )
        var params httprouter.Params
	var js []byte
	sess := model.Instance(r)
        params      = context.Get(r, "params").(httprouter.Params)
	sfec       :=  params.ByName("fec")[:10]
	dtfec,err  :=  time.Parse(layout, sfec)
//	sess       := model.Instance(r)
         fAnte    := dtfec.AddDate(0,-1,0)
	 iAmount  := int64(0)
//            fmt.Println("JBalanGET 1",  sfec)
	 lsAmount,err  :=  model.GetAmount(fAnte)
//            fmt.Println("JBalanGET 2", lsAmount, sfec)
        if err != nil {
          sess.AddFlash(view.Flash{"No balance periodo ", view.FlashError})
            log.Println(err)
            log.Fatalln(err)
         }
	 if len(lsAmount) > 0{
            iAmount = lsAmount[0]
            js, err =  json.Marshal(iAmount)
            if err == nil{
               w.Header().Set("Content-Type", "application/json")
               w.Write(js)
	       return
            }else{
            log.Println(err)

            }
        }
         log.Println("JBalanGET  ", err)
          http.Error(w, err.Error(), http.StatusInternalServerError)
          return

 }
// ---------------------------------------------------
// BalanGET despliega formulario crear balances
func BalanGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriod, err := model.Periods()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay periodos ", view.FlashError})
         }
        var lisBalans []model.BalanceN
        lisBalans, err      =  model.Balans()
        if err != nil {
             sess.AddFlash(view.Flash{"No Balances ", view.FlashError})
         }
	v                  := view.New(r)
	v.Name              = "balance/balan"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["LisBalans"] = lisBalans
        v.Vars["LisPeriod"] = lisPeriod
	v.Render(w)
 }
// ---------------------------------------------------
 func getFormBalan(b *  model.BalanceN, r *http.Request)(err error){
	   var nro int64
           b.PeriodId, _   = atoi32(r.FormValue("periodId"))
           nro, err        = money2int64(r.FormValue("amount"))
           if err == nil {
                 b.Amount  =  nro
                 nro,err   = money2int64(r.FormValue("cuota"))
            }
           if err == nil {
                 b.Cuota   =  nro
            }
       return
   }
// ---------------------------------------------------
// BalanPOST procesa la forma enviada con los datos
func BalanPOST(w http.ResponseWriter, r *http.Request) {
        var balan model.BalanceN
	sess          := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
           if validate, missingField := view.Validate(r, []string{"amount",  "cuota"}); !validate {
               sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
               sess.Save(r, w)
               BalanGET(w, r)
               return
	     }
           getFormBalan(&balan, r)
           err := (&balan).BalanByPeriod()
           if err == model.ErrNoResult { // Exito:  no hay balance creado aun 
               ex := (&balan).BalanCreate()
               log.Println("Creating balance")
	       if ex != nil {  // uyy como fue esto ? 
                   log.Println(ex)
                   sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
                   sess.Save(r, w)
                   return
	       }
               sess.AddFlash(view.Flash{"Balance creado: " , view.FlashSuccess})
               sess.Save(r, w)
	   }
      }
	http.Redirect(w, r, "/balance/list", http.StatusFound)
 }

// ---------------------------------------------------
// BalanUpGET despliega la pagina del usuario
func BalanUpGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var balan model.BalanceN
	var params httprouter.Params
	params  = context.Get(r, "params").(httprouter.Params)
	id,_   := atoi32(params.ByName("id"))
        path   :=  "/balance/list"
        balan.Id = id
	err := (&balan).BalanById()
	if err != nil { // Si no existe balance
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No esta balance.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "balance/balanupdate"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Balan"]       = balan
        v.Render(w)
   }

//	     nro, _  =  money2int64(  r.FormValue("cuota") )
// ---------------------------------------------------
 func   getBalanFormUp(b1, b2 model.BalanceN, r * http.Request)(stup string){
        var sf string
        var sup []string

        if   b1.Cuota != b2.Cuota {
             sf  =  fmt.Sprintf( " cuota = %d ", b2.Cuota )
	     sup = append(sup, sf)
           }
        if b1.Amount != b2.Amount {
             sf  =  fmt.Sprintf( " amount = %d ", b2.Amount )
	     sup = append(sup, sf)
           }
          lon := len(sup)
         if lon > 0 {
            sini        :=  "update balances set "
	    now        := time.Now()
	    sf          =  fmt.Sprintf( " , updated_at = '%s' ", now.Format(layout) )
            stup         =  strings.Join(sup, ", ")
            sr          :=  fmt.Sprintf(" where balances.id = %d ", b1.Id)
	    stup         =  sini + stup + sf + sr
          }
// fmt.Println(stup)
         return
  }
// ---------------------------------------------------
// BalanUpPOST procesa la forma enviada con los datos
func BalanUpPOST(w http.ResponseWriter, r *http.Request) {
        var balan, bal model.BalanceN
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        bal.Id     = Id
        path        :=  "/balance/list"
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            getFormBalan(&balan, r)
	    err := (&bal).BalanById()
            st          :=  getBalanFormUp(bal, balan,r)
            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No hay actualizacion solicitada", view.FlashSuccess})
            } else {
             err   =  balan.BalanUpdate(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Balance actualizada exitosamente : " , view.FlashSuccess})
             } else       {
		log.Println(err)
		sess.AddFlash(view.Flash{"Un error ocurrio actualizando.", view.FlashError})
	     }
		sess.Save(r, w)
           }
        }
	http.Redirect(w, r, path, http.StatusFound)
 }
//------------------------------------------------
// BalanLisGET displays the balance page
func BalanLisGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisBalan, err := model.Balans()
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Balances.", view.FlashError})
            sess.Save(r, w)
         }
	v := view.New(r)
	v.Name             = "balance/balanlis"
	v.Vars["token"]    = csrfbanana.Token(w, r, sess)
        v.Vars["LisBalan"] = lisBalan
        v.Vars["Level"]    =  sess.Values["level"]
	v.Render(w)
 }

//------------------------------------------------
// UserDeleteGET handles the note deletion
 func BalanDeleteGET(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        var balan model.BalanceN
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
	id,_        := atoi32(params.ByName("id"))
        path        :=  "/balance/list"
        balan.Id      = id
	err         := (&balan).BalanById()
	if err != nil { // Si no existe el balance
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No esta balance.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "balance/balandelete"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     =  "Eliminar Balance"
        v.Vars["Action"]    =  "/balance/delete"
        v.Vars["Balan"]        = balan
	v.Render(w)
  }

// ---------------------------------------------------
// BalanUpPOST procesa la forma enviada con los datos
func BalanDeletePOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var balan model.Balance
	sess := model.Instance(r)
        var params httprouter.Params
        params       = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        balan.Id      = Id
        path        :=  "/balance/list"
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
             err = balan.BalanDelete()
             if err != nil {
                 log.Println(err)
                 sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
              } else {
                  sess.AddFlash(view.Flash{"Balance borrado!", view.FlashSuccess})
              }
              sess.Save(r, w)
        }
	http.Redirect(w, r, path, http.StatusFound)
 }
// ---------------------------------------------------
