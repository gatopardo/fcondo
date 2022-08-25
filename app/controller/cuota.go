package controller

import (
	"log"
	"net/http"
        "fmt"
        "strings"
        "time"
        "encoding/json"

	"github.com/gatopardo/acad/ch14/fcondo/app/model"
	"github.com/gatopardo/acad/ch14/fcondo/app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )
  // -------------------------------------------------:--
  // JCondoGET reporte ingresos cuotas
 func JCondoGET(w http.ResponseWriter, r *http.Request) {
	var perid model.Periodo
        var lisAmt     []model.AmtCond
        var amtPerCond   model.AmtPerCond
        var params httprouter.Params
	sess := model.Instance(r)
        params      = context.Get(r, "params").(httprouter.Params)
	sfec       :=  params.ByName("fec")[:10]
	dtfec,err  :=  time.Parse(layout, sfec)
        if err != nil {
              sess.AddFlash(view.Flash{"Formato Fecha Errado ", view.FlashError})
	      log.Println(err)
	}else{
        dtfec       =  time.Date(dtfec.Year(), dtfec.Month(),dtfec.Day(), 0, 0, 0, 0, time.Local)
        err         = (&perid).PeriodByFec(dtfec)
        if err     != nil {
	        log.Println(err)
        }else{
	  lisAmt, err           = model.Amounts( perid.Id )
          if err != nil {
            log.Println(err)
          }else{
	    amtPerCond.Fecha = perid.Inicio
	    amtPerCond.LisAmt = lisAmt
            var js []byte
            js, err =  json.Marshal(amtPerCond)
            if err == nil{
               w.Header().Set("Content-Type", "application/json")
               w.Write(js)
	       return
            }
           }
          }
          }
          log.Println("JCondo  ", err)
          http.Error(w, err.Error(), http.StatusInternalServerError)
          return
 }

// ---------------------------------------------------
// JCuotGET reporte ingresos cuotas
 func JCuotGET(w http.ResponseWriter, r *http.Request) {
	var periodo model.Periodo
        var lisCuot  []model.CuotaN
        var params httprouter.Params
	sess := model.Instance(r)
        params      = context.Get(r, "params").(httprouter.Params)
	sfec       :=  params.ByName("fec")[:10]
	dtfec,err  :=  time.Parse(layout, sfec)
        if err != nil {
             sess.AddFlash(view.Flash{"Formato Fecha Errado ", view.FlashError})
	        log.Println("JCuotGET",err)
	}else{
        dtfec       =  time.Date(dtfec.Year(), dtfec.Month(),dtfec.Day(), 0, 0, 0, 0, time.Local)
        err         = (&periodo).PeriodByFec(dtfec)
        if err     != nil {
	        log.Println(err)
        }else{
          lisCuot, err           = model.CuotLim( periodo.Id )
          if err != nil {
             sess.AddFlash(view.Flash{"No hay cuotas Periodo ", view.FlashError})
            log.Println(err)
          }else{
            var js []byte
            js, err =  json.Marshal(lisCuot)
            if err == nil{
               w.Header().Set("Content-Type", "application/json")
               w.Write(js)
	       return
            }
           }
          }
          }
          log.Println("JCuot  ", err)
          http.Error(w, err.Error(), http.StatusInternalServerError)
          return
 }
  // -------------------------------------------------:--
// CuotPerGET despliega formulario escoger periodo
func CuotPerGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriod, err := model.Periods()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay periodos ", view.FlashError})
	     log.Println("CuotPerGet",err)
         }
	v                  := view.New(r)
	v.Name              = "cuota/cuotper"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     =  "Escoger Periodo"
        v.Vars["Action"]    =  "/cuota/periodo/register"
        v.Vars["LisPeriod"] = lisPeriod
	v.Render(w)
 }
 //----------------------------------------------------
// CuotPerPOST procesa la forma enviada con periodo
func CuotPerPOST(w http.ResponseWriter, r *http.Request) {
        var cuot model.CuotaN
        var period model.Periodo
	sess          := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            var lisTipo []model.Tipo
            var lisCuot []model.CuotaN
            cuot.PeriodId,  _   =  atoi32(r.FormValue("id"))
            period.Id           =  cuot.PeriodId
            _                   =  (&period).PeriodById()
            cuot.Period         =  period.Inicio
            lisApts, err       :=  model.Apts()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay aptos ", view.FlashError})
            }
            lisTipo,  err        = model.Tipos()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay tipos ", view.FlashError})
             }
             lisCuot, _          = (&cuot).CuotsPer()
	     v                  := view.New(r)
	     v.Name              = "cuota/cuotreg"
             v.Vars["token"]     = csrfbanana.Token(w, r, sess)
	     v.Vars["Title"]     =  "Crear Cuota"
             v.Vars["Action"]    =  "/cuota/register"
             v.Vars["Cuot"]      = cuot
             v.Vars["LisApt"]    = lisApts
             v.Vars["LisTip"]    = lisTipo
             v.Vars["LisCuots"]  = lisCuot
             v.Render(w)
        }
//	http.Redirect(w, r, "/cuota/list", http.StatusFound)
 }
// ---------------------------------------------------
 func getFormCuot(c *  model.CuotaN, r *http.Request, b bool)(err error){
           formato         :=  "2006/01/02"
           formato2        :=  "2006-01-02"
           c.ApartaId, _    =  atoi32(r.FormValue("aptId"))
           c.TipoId, _      =  atoi32(r.FormValue("tipId"))
           c.PeriodId, _  =  atoi32(r.FormValue("periodId"))
           stPeriod        := r.FormValue("period")
           stFecha         := r.FormValue("fecha")
           c.Period,_       =  time.Parse(formato,stPeriod)
	   if b {
              c.Fecha, _    =  time.Parse(formato2,stFecha)
            }else{
              c.Fecha, _    =  time.Parse(formato,stFecha)
	    }
           ramount         :=  r.FormValue("amount")
           samount         :=   strings.ReplaceAll(ramount, ",","")
           unr, err        :=  money2int64(samount)
           if err == nil {
                 c.Amount   =  unr
            }
       return
   }
// ---------------------------------------------------
// CuotRegPOST despliega formulario crear cuota
func CuotRegPOST(w http.ResponseWriter, r *http.Request) {
        var cuot   model.CuotaN
        var period  model.Periodo
        var err  error
	sess   := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
           getFormCuot(&cuot, r, true)
           period.Id           =  cuot.PeriodId
           err                 =  (&period).PeriodById()
           err                 =  (&cuot).CuotCreate()
           if err != nil {  // uyy como fue esto ? 
               log.Println(err)
               sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
               return
           } else {  // todo bien
                sess.AddFlash(view.Flash{"Cuota. creada: " , view.FlashSuccess})
           }
            var lisApto []model.Aparta
            var lisTipo []model.Tipo
            var lisCuot []model.CuotaN
            lisApto, err  = model.Apts()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay apartas ", view.FlashError})
            }
            lisTipo, err  = model.Tipos()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay tipos ", view.FlashError})
            }
            lisCuot, err         = (&cuot).CuotsPer()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay cuotas ", view.FlashError})
            }
            v                   := view.New(r)
            v.Name               = "cuota/cuotreg"
            v.Vars["token"]      = csrfbanana.Token(w, r, sess)
	    v.Vars["Title"]     =  "Guardar Cuota"
            v.Vars["Action"]    =  "/cuota/register"
            v.Vars["Cuot"]      = cuot
            v.Vars["LisApt"]     = lisApto
            v.Vars["LisTip"]     = lisTipo
            v.Vars["LisCuots"]   = lisCuot
	    v.Render(w)
        }
	http.Redirect(w, r, "/cuota/list", http.StatusFound)
 }
// ---------------------------------------------------
// CuotUpGET despliega la pagina del usuario
func CuotUpGET(w http.ResponseWriter, r *http.Request) {
        var lisTipo []model.Tipo
	sess := model.Instance(r)
        var cuot model.CuotaN
	var params httprouter.Params
	params   = context.Get(r, "params").(httprouter.Params)
	id,_    := atoi32(params.ByName("id"))
        cuot.Id  = id
        path    :=  "/cuota/list"
        lisApts, err       :=  model.Apts()
        if err  != nil {
             sess.AddFlash(view.Flash{"No hay aptos ", view.FlashError})
        }
        lisTipo,  err        = model.Tipos()
        if err   != nil {
             sess.AddFlash(view.Flash{"No hay tipos ", view.FlashError})
        }
	err = (&cuot).CuotById()
	if err != nil { // Si no existe cuota
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No esta cuota.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    :=  view.New(r)
	v.Name                =  "cuota/cuotupdate"
	v.Vars["token"]       =  csrfbanana.Token(w, r, sess)
	v.Vars["Title"]       =  "Actualizar Cuota"
        v.Vars["Action"]      =  "/cuota/update"
        v.Vars["Cuot"]        =  cuot
        v.Vars["LisApt"]      =  lisApts
        v.Vars["LisTip"]      =  lisTipo
        v.Render(w)
   }

// ---------------------------------------------------
 func   getCuotFormUp(c1,c2 model.CuotaN, r * http.Request)(stUp string){
        var sf string
	var sup  []string
        formato        :=  "2006/01/02"
	if c1.ApartaId != c2.ApartaId {
             sf  =  fmt.Sprintf( " aparta_id = %d ", c1.ApartaId )
	     sup = append(sup, sf)
	}

	if c1.TipoId != c2.TipoId {
             sf  =  fmt.Sprintf( " tipo_id = %d ", c1.TipoId )
	     sup = append(sup, sf)
	}

	if c1.Fecha.Format(formato) != c2.Fecha.Format(formato) {
             sf  =  fmt.Sprintf( " fecha = '%s' ", c1.Fecha.Format(formato) )
	     sup = append(sup, sf)
	}

	if c1.Amount != c2.Amount {
             sf  =  fmt.Sprintf( " amount = %d ", c1.Amount )
	     sup = append(sup, sf)
	}
       lon := len(sup)
       if lon  > 0 {
	    now         := time.Now()
	    sf           =  fmt.Sprintf( " , updated_at = '%s' ", now.Format(formato) )
            sini        :=  "update cuotas set "
            stUp         =  strings.Join(sup, ", ")
	    sr          :=  fmt.Sprintf(" where cuotas.id = %d ", c1.Id)

            stUp         = sini + stUp + sf +  sr
       }
         return
  }
  // ---------------------------------------------------
// CuotUpPOST procesa la forma enviada con los datos
func CuotUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var c, cuot model.CuotaN
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        cuot.Id      = Id
        c.Id         = Id
        path        :=  "/cuota/list"
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
	    err = (&cuot).CuotById()
	    if err != nil { // Si no existe cuota
                  sess.AddFlash(view.Flash{"Es raro. No esta cuota.", view.FlashError})
            }
	    getFormCuot(&c,r, false)
            st          :=  getCuotFormUp(c, cuot, r)
            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No actualizacion solicitada", view.FlashSuccess})
            } else {
             err   =  cuot.CuotUpdate(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Cuota actualizada exitosamente : " , view.FlashSuccess})
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
// CuotLisGET displays the cuot page
func CuotLisGET(w http.ResponseWriter, r *http.Request) {
	sess            := model.Instance(r)
        lisPeriod,err   := model.Periods()
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Obteniendo Periodos.", view.FlashError})
            sess.Save(r, w)
         }
	v                   := view.New(r)
	v.Name               = "cuota/cuotper"
	v.Vars["token"]      = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     =  "Listar"
        v.Vars["Action"]    =  "/cuota/list"
        v.Vars["LisPeriod"]  = lisPeriod
	v.Render(w)
 }

//------------------------------------------------
// CuotLis displays the cuot page
func CuotLisPOST(w http.ResponseWriter, r *http.Request) {
        var Id  uint32
	var per model.Periodo
	sess            := model.Instance(r)
        lisPeriod,err    := model.Periods()
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Obteniendo Periodos.", view.FlashError})
            sess.Save(r, w)
         }
        Id,_             = atoi32(r.FormValue("id"))
	per.Id           = Id
	err              = (&per).PeriodById()
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error con Periodo.", view.FlashError})
            sess.Save(r, w)
         }

        lisCuot, err         := model.CuotLim(Id)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Cuotas.", view.FlashError})
            sess.Save(r, w)
         }
	v                   := view.New(r)
	v.Name               = "cuota/cuotlis"
	v.Vars["token"]      = csrfbanana.Token(w, r, sess)
	v.Vars["Per"]        = per
        v.Vars["LisPeriod"]    = lisPeriod
        v.Vars["LisCuot"]    = lisCuot
        v.Vars["Level"]      =  sess.Values["level"]
	v.Render(w)
 }

//------------------------------------------------
// CuotDeleteGET handles the note deletion
 func CuotDeleteGET(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        var cuot model.CuotaN
        var params httprouter.Params
        params       = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
	id,_        := atoi32(SId)
        path        :=  "/cuota/list"
        cuot.Id      = id
	err         := (&cuot).CuotById()
	if err != nil { // Si no existe cuota
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No esta cuota.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "cuota/cuotdelete"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     =  "Eliminar Cuota"
        v.Vars["Action"]    =  "/cuota/delete"
        v.Vars["Cuot"]        = cuot
	v.Render(w)
  }

// ---------------------------------------------------
// CuotDeletePOST procesa la forma enviada con los datos
func CuotDeletePOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var cuot model.Cuota
	sess := model.Instance(r)
        var params httprouter.Params
        params       = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        cuot.Id      = Id
        path        :=  "/cuota/list"
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
             err = cuot.CuotDelete()
             if err != nil {
                 log.Println(err)
                 sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
              } else {
                  sess.AddFlash(view.Flash{"Cuota borrado!", view.FlashSuccess})
              }
              sess.Save(r, w)
        }
	http.Redirect(w, r, path, http.StatusFound)
 }
// ---------------------------------------------------
