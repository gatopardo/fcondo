package controller

import (
	"log"
	"net/http"
        "strings"
        "fmt"
        "time"
        "encoding/json"

	"github.com/gatopardo/acad/ch14/fcondo/app/model"
	"github.com/gatopardo/acad/ch14/fcondo/app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )
//      Refill any form fields
// view.Repopulate([]string{"name"}, r.Form, v.Vars)
// ---------------------------------------------------
// japt get json service for apt state
 func JAptGET(w http.ResponseWriter, r *http.Request) {
        var params httprouter.Params
        var jpers  model.Jperson
	var peridi, peridf model.Periodo
        var lisPaym []model.CuotApt
	var arPaym ArPay
	var dt11, dt22  time.Time
	var err  error
        params           = context.Get(r, "params").(httprouter.Params)
	sfec1           :=  params.ByName("fec1")[:10]
	sfec2           :=  params.ByName("fec2")[:10]
        sId             :=  params.ByName("id")
        uid,_           :=  atoi32(sId)
	dt11,err         =  time.Parse(layout, sfec1)
        if err == nil {
            dt11       =  time.Date(dt11.Year(), dt11.Month(),dt11.Day(), 0, 0, 0, 0, time.Local)
	    dt22,err   =  time.Parse(layout, sfec2)
            dt22       =  time.Date(dt22.Year(), dt22.Month(),dt22.Day(), 0, 0, 0, 0, time.Local)
            err        =  (&peridi).PeriodByFec(dt11)
	    if err    ==  nil {
               err     =  (&peridf).PeriodByFec(dt22)
            }
        }
        if err      != nil {
	        log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	        return
        }
	_, err         =   (&jpers).JPersByUserId(uid)
	if err == model.ErrNoResult {
           log.Println("JAPTGET ", err)
           http.Error(w, err.Error(), http.StatusBadRequest)
	   return
        }
        lisPaym, err   =  model.Payments(jpers.AptId, peridf.Inicio, peridi.Inicio)
        if err == nil {
           arPaym.Apto   = jpers.Apto
	   arPaym.Final  = peridf.Final
	   arPaym.APaym  = lisPaym
           var js []byte
           js, err =  json.Marshal(arPaym)
           if err == nil{
               w.Header().Set("Content-Type", "application/json")
               w.Write(js)
	       return
           }
	}
           log.Println("JAPTGET 2 ", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
	return
 }
// ---------------------------------------------------
// AptGET despliega la pagina del apto
func AptGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
        lisApts, _ := model.Apts()
	v := view.New(r)
	v.Name = "aparta/apt"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["Title"]  = "Crear Apto"
	v.Vars["Action"]  = "/apto/register"
        v.Vars["LisApts"] =  lisApts
	v.Render(w)
 }
// ---------------------------------------------------
// POST procesa la forma enviada con los datos
func AptPOST(w http.ResponseWriter, r *http.Request) {
        var apt model.Aparta
	sess := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
	    if validate, missingField := view.Validate(r, []string{"codigo"}); !validate {
                 sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
                 sess.Save(r, w)
                 AptGET(w, r)
                 return
	      }
              apt.Codigo           = r.FormValue("codigo")
              apt.Descripcion      = r.FormValue("descripcion")
              err := (&apt).AptByCode()
              if err == model.ErrNoResult { // Exito: no hay apartamento creado aun 
                  ex := (&apt).AptCreate()
	          if ex != nil {  // uyy como fue esto ? 
                     log.Println(ex)
                     sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
                     return
	         } else {  // todo bien
                    sess.AddFlash(view.Flash{"Apto. creado: " +apt.Codigo, view.FlashSuccess})
	         }
              }
          }
          sess.Save(r, w)
	  http.Redirect(w, r, "/apto/list", http.StatusFound)
  }

// ---------------------------------------------------
// AptUpGET despliega la pagina del usuario
func AptUpGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var apt model.Aparta
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	id,_ := atoi32(params.ByName("id"))
	SPag   := params.ByName("pg")
        path   :=  fmt.Sprintf("/apto/list/%s", SPag)
        apt.Id = id
        err := (&apt).AptById()
	if err != nil { // Si no existe el usuario
            log.Println(err)
            sess.AddFlash(view.Flash{"Es raro. No tenemos apartamento.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, path, http.StatusFound)
            return
        }
	v := view.New(r)
	v.Name = "aparta/aptupdate"
	v.Vars["token"]  = csrfbanana.Token(w, r, sess)
	v.Vars["Title"]  = "Actualizar Apto"
	v.Vars["Action"]  = "/apto/update"
        v.Vars["Apto"] = apt
        v.Render(w)
   }
// ---------------------------------------------------
 func   getAptFormUp(a1, a2 model.Aparta, r * http.Request)(stup string){
        var sf string
        var sup []string

        if a1.Descripcion != a2.Descripcion {
	     sf  =  fmt.Sprintf( " descrip = '%s' ", a2.Descripcion )
	     sup = append(sup, sf)
        }
        if a1.Codigo != a2.Codigo {
	     sf  =  fmt.Sprintf( " codigo = '%s' ", a2.Codigo )
	     sup = append(sup, sf)
        }
	lon  := len(sup)
        if  lon > 0 {
            sini        :=  "update apartas set "
	    now         := time.Now()
	    sf           =  fmt.Sprintf( ",  updated_at = '%s' ", now.Format(layout) )
            stup =  strings.Join(sup, ", ")
            sr          :=  fmt.Sprintf(" where apartas.id = %d ", a1.Id)
             stup = sini + stup + sf + sr
        }
        return
  }
// ---------------------------------------------------
// AptUpPOST procesa la forma enviada con los datos
func AptUpPOST(w http.ResponseWriter, r *http.Request) {
        var apt , a2 model.Aparta
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        apt.Id      = Id
	a2.Id       = Id
        path        :=  "/apto/list"
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            err                 :=  (&apt).AptById()
	    if err != nil{
                 log.Println(err)
                 sess.AddFlash(view.Flash{"Error No Apto", view.FlashError})
	    }
            a2.Codigo           = r.FormValue("codigo")
            a2.Descripcion      = r.FormValue("descripcion")
            st                 :=  getAptFormUp(apt, a2, r)
            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No hay actualizacion solicitada", view.FlashSuccess})
            } else {
             err =  apt.AptUpdate(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Apto actualizado exitosamente para: " +apt.Codigo, view.FlashSuccess})
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
// AptLisGET displays the aparta page
func AptLisGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisApts, err := model.Apts()
        if err != nil {
           log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Aptos.", view.FlashError})
            sess.Save(r, w)
         }
	// Display the view
	v := view.New(r)
	v.Name = "aparta/aptlis"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
        v.Vars["LisApt"]   = lisApts
        v.Vars["Level"]    =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
// AptDeleteGET handles the apto deletion
 func AptDeleteGET(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        var apt model.Aparta
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        Id,_ := atoi32(params.ByName("id"))
        apt.Id = Id
        path        :=  "/apto/list"
        err := (&apt).AptById()
        if err != nil {
            log.Println(err)
            sess.AddFlash(view.Flash{"Es raro no tenemos apto.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, path, http.StatusFound)
            return
        }
	v := view.New(r)
	v.Name = "aparta/aptdelete"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     =  "Eliminar Aparta"
        v.Vars["Action"]    =  "/apto/delete"
        v.Vars["Apto"]      =  apt
	v.Render(w)
  }
// ---------------------------------------------------
// AptDeletePOST handles the apto deletion
 func AptDeletePOST(w http.ResponseWriter, r *http.Request) {
        var err error
        sess := model.Instance(r)
        var apt model.Aparta
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        Id,_ := atoi32(params.ByName("id"))
        apt.Id = Id
        path        :=  "/apto/list"
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            err  = apt.AptDelete()
            if err != nil {
                log.Println(err)
                sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
            } else {
                sess.AddFlash(view.Flash{"Apto. borrado!", view.FlashSuccess})
            }
            sess.Save(r, w)
        }
	http.Redirect(w, r, path, http.StatusFound)

  }
// ---------------------------------------------------

