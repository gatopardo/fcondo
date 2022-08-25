package controller

import (
	"log"
	"net/http"
        "strings"
        "fmt"
        "time"

	"github.com/gatopardo/acad/ch14/fcondo/app/model"
	"github.com/gatopardo/acad/ch14/fcondo/app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )

//      Refill any form fields
// view.Repopulate([]string{"name"}, r.Form, v.Vars)
// ---------------------------------------------------

// TipoGET despliega la pagina del tipo
func TipoGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
        lisTipos, _ := model.Tipos()
	lisApli     :=   make([]model.Aplic,2,2)
	lisApli     =   model.ApliLis()
	v := view.New(r)
	v.Name = "tipo/tipo"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["Title"]  = "Crear Tipo"
	v.Vars["Action"]  = "/tipo/register"
        v.Vars["LisApli"] =  lisApli
        v.Vars["LisTipos"] =  lisTipos
	v.Render(w)
 }
// ---------------------------------------------------
// ---------------------------------------------------
// POST procesa la forma enviada con los datos
func TipoPOST(w http.ResponseWriter, r *http.Request) {
        var tipo model.Tipo
	sess      := model.Instance(r)
        action    := r.FormValue("action")
	path      :=  "/tipo/list"
        if ! (strings.Compare(action,"Cancelar") == 0) {
	    if validate, missingField := view.Validate(r, []string{"codigo"}); !validate {
                 sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
                 sess.Save(r, w)
                 TipoGET(w, r)
                 return
	      }
              tipo.Codigo           = r.FormValue("codigo")
              tipo.Aplica           = r.FormValue("aplica")
              tipo.Descripcion      = r.FormValue("descripcion")
              err := (&tipo).TipoByCode()
              if err == model.ErrNoResult { // Exito: no hay tipo creado aun 
                  ex := (&tipo).TipoCreate()
	          if ex != nil {  // uyy como fue esto ? 
                     log.Println(ex)
                     sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
                     return
	         } else {  // todo bien
                    sess.AddFlash(view.Flash{"Tipoo. creado: " +tipo.Codigo, view.FlashSuccess})
	         }
              }
          }
          sess.Save(r, w)
	  http.Redirect(w, r, path, http.StatusFound)
  }

// ---------------------------------------------------
// TipoUpGET despliega la pagina del usuario
func TipoUpGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var tipo model.Tipo
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	id,_ := atoi32(params.ByName("id"))
        path   :=  "/tipo/list"
        tipo.Id = id
	 err := (&tipo).TipoById()
	if err != nil { // Si no existe el usuario
            log.Println(err)
            sess.AddFlash(view.Flash{"Es raro. No tenemos tipo.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, path, http.StatusFound)
            return
        }
	v := view.New(r)
	v.Name = "tipo/tipoupdate"
	v.Vars["token"]  = csrfbanana.Token(w, r, sess)
	v.Vars["Title"]  = "Actualizar Categoria"
	v.Vars["Action"]  = "/tipo/update"
        v.Vars["Tipo"] = tipo
        v.Render(w)
   }
   // ---------------------------------------------------
 func   getTipoFormUp(t1, t2 model.Tipo , r * http.Request)(stup string){
        var sf string
        var sup []string

        if t1.Descripcion != t2.Descripcion {
             sf  =  fmt.Sprintf( " descrip = '%s' ", t2.Descripcion )
             sup = append(sup, sf)
        }
        if t1.Codigo != t2.Codigo {
             sf  =  fmt.Sprintf( " descrip = '%s' ", t2.Codigo )
             sup = append(sup, sf)
        }
	lon :=  len(sup)
        if lon > 0 {
            sini       :=  "update tipos set "
	    now        := time.Now()
	    sf          =  fmt.Sprintf( " , updated_at = '%s' ", now.Format(layout) )
            sr         :=  fmt.Sprintf(" where tipo.id = %d ", t1.Id)
            stup        =  strings.Join(sup, ", ")
            stup        = sini + stup + sf + sr

        }
        return
  }
// ---------------------------------------------------

// ---------------------------------------------------
// TipoUpPOST procesa la forma enviada con los datos
func TipoUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var tipo, t2 model.Tipo
	// Get session
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        tipo.Id      = Id
        t2.Id        = Id
        path        := "/tipo/list"
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            st          :=  getTipoFormUp(tipo, t2, r)
            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No hay actualizacion solicitada", view.FlashSuccess})
            } else {
             err =  tipo.Update(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Tipo actualizado exitosamente para: " +tipo.Codigo, view.FlashSuccess})
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
// TipoLisGET displays the tipo page
func TipoLisGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisTipos, err := model.Tipos()
        if err != nil {
           log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Tipos.", view.FlashError})
            sess.Save(r, w)
         }
	v                 := view.New(r)
	v.Name             = "tipo/tipolis"
	v.Vars["token"]    = csrfbanana.Token(w, r, sess)
        v.Vars["LisTipo"]   = lisTipos
        v.Vars["Level"]    =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
// TipoDeleteGET handles the tipo deletion
 func TipoDeleteGET(w http.ResponseWriter, r *http.Request) {
	// Get session
        sess := model.Instance(r)
        var tipo model.Tipo
        var params httprouter.Params
        params   = context.Get(r, "params").(httprouter.Params)
        Id,_    := atoi32(params.ByName("id"))
        tipo.Id  = Id
        path    :=  "/tipo/list"
        err := (&tipo).TipoById()
        if err != nil {
            log.Println(err)
            sess.AddFlash(view.Flash{"Es raro no tenemos tipo.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, path, http.StatusFound)
            return
        }
	v := view.New(r)
	v.Name = "tipo/tipodelete"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     =  "Eliminar Tipo"
        v.Vars["Action"]    =  "/tipo/delete"
        v.Vars["Tipo"]      =  tipo
	v.Render(w)
  }
// ---------------------------------------------------
// TipoDeletePOST handles the tipo deletion
 func TipoDeletePOST(w http.ResponseWriter, r *http.Request) {
        var err error
        sess := model.Instance(r)
        var tipo model.Tipo
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        Id,_ := atoi32(params.ByName("id"))
        tipo.Id = Id
        path        :=  "/tipo/list"
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            err  = tipo.Delete()
            if err != nil {
                log.Println(err)
                sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
            } else {
                sess.AddFlash(view.Flash{"Tipo. borrado!", view.FlashSuccess})
            }
            sess.Save(r, w)
        }
	http.Redirect(w, r, path, http.StatusFound)

  }
// ---------------------------------------------------

