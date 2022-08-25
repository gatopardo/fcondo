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
  var formato         =  "2006/01/02"
  var formato2        =  "2006-01-02"
//------------------------------------------------
 func JIngreGET(w http.ResponseWriter, r *http.Request) {
	var periodo model.Periodo
	var lisIngre  []model.IngresoJ
	var  ingresoL  model.IngresoL
        var js []byte
        var params httprouter.Params
	sess := model.Instance(r)
        params      = context.Get(r, "params").(httprouter.Params)
	sfec       :=  params.ByName("fec")[:10]
	dtfec,err  :=  time.Parse(layout, sfec)
        if err != nil {
//	        fmt.Println(err)
	        log.Println(err)
	}else{
//        fmt.Println(" JIngreGET fec:",sfec, " - ", dtfec )
        dtfec       =  time.Date(dtfec.Year(), dtfec.Month(),dtfec.Day(), 0, 0, 0, 0, time.Local)
//        fmt.Println(" JIngreGET fec:",sfec, " - ", dtfec )
        err         = (&periodo).PeriodByFec(dtfec)
        if err     != nil {
//	        fmt.Println(err)
	        log.Println(err)
        }else{
          lisIngre, err  = model.IngresoJPer( periodo.Id )
	  var  num int64 = 0
          if err != nil {
             sess.AddFlash(view.Flash{"No ingreso periodo ", view.FlashError})
            log.Println(err)
          }else {
            if len(lisIngre) == 0 {
	        var dfec time.Time
                dfec,_ = time.Parse(formato , "1900/01/01")
		ingresoj := model.IngresoJ{Tipo:"X" , Fecha: dfec, Amount: num, Descripcion: ""}
                lisIngre =  append( lisIngre, ingresoj)
            }
            ingresoL.Period  =  periodo.Inicio
            ingresoL.LisIngre =  lisIngre
//        fmt.Println(" JIngreGET ",ingresoL )
            js, err =  json.Marshal(ingresoL)
            if err == nil{
               w.Header().Set("Content-Type", "application/json")
               w.Write(js)
	       return
            }
           }
          }
          }
          log.Println("JIngre  ", err)
          http.Error(w, err.Error(), http.StatusInternalServerError)
          return
 }
// ---------------------------------------------------

// ---------------------------------------------------
// IngrePerGET despliega formulario escoger periodo
func IngrePerGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriod, err := model.Periods()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay periodos ", view.FlashError})
         }
	v                  := view.New(r)
	v.Name              = "ingreso/ingresoper"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     =  "Ingreso"
        v.Vars["Action"]    =  "/ingreso/periodo/register"
        v.Vars["LisPeriod"] = lisPeriod
	v.Render(w)
 }
// ---------------------------------------------------
// IngrePerPOST procesa la forma enviada con periodo
func IngrePerPOST(w http.ResponseWriter, r *http.Request) {
        var ingres model.IngresoN
        var period model.Periodo
        var  err  error
	sess          := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            ingres.PeriodId,  _   =  atoi32(r.FormValue("id"))
            period.Id            =  ingres.PeriodId
            _                    =  (&period).PeriodById()
            ingres.Period         =  period.Inicio
            var lisTipo []model.Tipo
            var lisIngre []model.IngresoN
            lisTipo,  err        = model.Tipos()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay tipos ", view.FlashError})
            }
            lisIngre, _          = model.IngresLim(ingres.PeriodId)
	    v                   := view.New(r)
	    v.Name               = "ingreso/ingresoreg"
            v.Vars["token"]      = csrfbanana.Token(w, r, sess)
	    v.Vars["Title"]     =  "Crear Ingreso"
            v.Vars["Action"]    =  "/ingreso/register"
            v.Vars["Ingreso"]    = ingres
            v.Vars["LisTip"]     = lisTipo
            v.Vars["LisIngres"]  = lisIngre
            v.Render(w)
        }
	http.Redirect(w, r, "/ingreso/list", http.StatusFound)
 }
// ---------------------------------------------------
 func getFormIngre(ing *  model.IngresoN, r *http.Request, b bool)(err error){
           ing.TipoId, _    =  atoi32(r.FormValue("tipId"))
           ing.PeriodId, _  =  atoi32(r.FormValue("periodId"))
           ing.Period, _    =  time.Parse(layout,r.FormValue("period"))
	   if b{
               ing.Fecha, _ =  time.Parse(formato2,r.FormValue("fecha"))
           }else{
               ing.Fecha, _ =  time.Parse(formato,r.FormValue("fecha"))
           }
	   ing.Descripcion   =  r.FormValue("descripcion")
	   var nro int64
	   ramount         :=  r.FormValue("amount")
	   samount         :=   strings.ReplaceAll(ramount, ",","")
           nro, err       = money2int64(samount)
           if err == nil {
		   ing.Amount   =  nro
            }
           ing.Descripcion   =  r.FormValue("descripcion")
       return
   }
// ---------------------------------------------------
// IngreRegPOST despliega formulario crear ingreso
func IngreRegPOST(w http.ResponseWriter, r *http.Request) {
        var ingres   model.IngresoN
        var period  model.Periodo
        var err  error
	sess   := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
           getFormIngre(&ingres, r, true)
           period.Id            =  ingres.PeriodId
        _                    =  (&period).PeriodById()

         err                 =  (&ingres).IngresCreate()
         if err != nil {  // uyy como fue esto ? 
               log.Println(err)
//               fmt.Println(err)
               sess.AddFlash(view.Flash{"Error guardando Ingreso.", view.FlashError})
	       http.Redirect(w, r, "/ingreso/list", http.StatusFound)
         } else {  // todo bien
                sess.AddFlash(view.Flash{"Ingreso. creado: " , view.FlashSuccess})
         }

            var lisTipo []model.Tipo
            var lisIngre []model.IngresoN
            lisTipo, err  = model.Tipos()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay tipos ", view.FlashError})
            }
            lisIngre, _           = model.IngresLim(period.Id)
            v                   := view.New(r)
            v.Name               = "ingreso/ingresoreg"
            v.Vars["token"]      = csrfbanana.Token(w, r, sess)
            v.Vars["Ingreso"]     = ingres
            v.Vars["LisTip"]     = lisTipo
            v.Vars["LisIngres"]   = lisIngre
	    v.Render(w)
        }
	http.Redirect(w, r, "/ingreso/list", http.StatusFound)
 }
// ---------------------------------------------------
// IngreUpGET despliega la pagina del usuario
func IngreUpGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var ingres model.IngresoN
	var params httprouter.Params
	params  = context.Get(r, "params").(httprouter.Params)
	id,_   := atoi32(params.ByName("id"))
        path   :=  "/ingreso/list"
        ingres.Id = id
	lisTipo,  err        := model.Tipos()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay tipos ", view.FlashError})
            }
	err = (&ingres).IngresById()
	if err != nil { // Si no existe el ingreso
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No esta ingreso.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "ingreso/ingresoupdate"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     =  "Actualizar Ingreso"
        v.Vars["Action"]    =  "/ingreso/update"
        v.Vars["Ingre"]       = ingres
        v.Vars["LisTip"]      = lisTipo
        v.Render(w)
   }

// ---------------------------------------------------
 func   getIngreFormUp(i1, i2 model.IngresoN ,r * http.Request)(stUp string){
        var sf string
        var sup []string

	if i1.TipoId != i2.TipoId {
             sf  =  fmt.Sprintf( " tipo_id = %d ", i1.TipoId )
	     sup = append(sup, sf)
	}
	if i1.Amount  != i2.Amount {
             sf  =  fmt.Sprintf( " amount = '%d' ", i1.Amount )
	     sup = append(sup, sf)
	}
        if i1.Fecha.Format(formato) != i2.Fecha.Format(formato) {
             sf  =  fmt.Sprintf( " fecha = '%s' ", i1.Fecha.Format(layout) )
	     sup = append(sup, sf)
	}
	if i1.Descripcion != i2.Descripcion {
             sf  =  fmt.Sprintf( " descripcion = '%s' ", i1.Descripcion )
	     sup = append(sup, sf)
	}
        lon := len(sup)
        if lon  > 0 {
            sini :=  "update ingresos set "
	    now         := time.Now()
	    sf           =  fmt.Sprintf( " , updated_at = '%s' ", now.Format(layout) )
            stUp  =  strings.Join(sup, ", ")
            sr   :=  fmt.Sprintf(" where ingresos.id = %d ", i1.Id)
            stUp = sini + stUp + sf + sr
       }
//fmt.Println(stUp)
       return
  }
// ---------------------------------------------------
// IngreUpPOST procesa la forma enviada con los datos
func IngreUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var ing, ingres model.IngresoN
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        ingres.Id    = Id
        ing.Id       = Id
        path        :=  "/ingreso/list"
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            err  = (&ingres).IngresById()
	    if err != nil { // Si no existe ingreso
                  sess.AddFlash(view.Flash{"Es raro. No esta ingreso.", view.FlashError})
            }
	    getFormIngre(&ing,r,false)
            st          :=  getIngreFormUp(ing, ingres, r)

            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No actualizacion solicitada", view.FlashSuccess})
            } else {
             err   =  ingres.IngresUpdate(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Ingreso actualizada exitosamente : " , view.FlashSuccess})
             } else       {
		log.Println(err)
		sess.AddFlash(view.Flash{"Un error ocurrio actualizando.", view.FlashError})
	     }
		sess.Save(r, w)
           }
        }
	http.Redirect(w, r, path, http.StatusFound)
 }
// ---------------------------------------------------
// IngreLisGET despliega formulario escoger periodo
func IngreLisGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriod, err := model.Periods()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay periodos ", view.FlashError})
         }
	v                  := view.New(r)
	v.Name              = "ingreso/ingresoper"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["LisPeriod"] = lisPeriod
        v.Vars["Title"]     =  "Listar"
        v.Vars["Action"]    =  "/ingreso/list"
	v.Render(w)
 }
//------------------------------------------------
// IngreLis displays the ingres page
func IngreLisPOST(w http.ResponseWriter, r *http.Request) {
        var Id  uint32
	var per model.Periodo
	sess            := model.Instance(r)
        lisPeriod,err    := model.Periods()
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error con Lista Periodos.", view.FlashError})
            sess.Save(r, w)
         }
         Id,_             = atoi32(r.FormValue("id"))
	per.Id                 = Id
	err = (&per).PeriodById();
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error con  Periodo.", view.FlashError})
            sess.Save(r, w)
         }
        lisIngre, err         := model.IngresLim(Id)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Ingresos.", view.FlashError})
            sess.Save(r, w)
         }
	v                   := view.New(r)
	v.Name               = "ingreso/ingresolis"
	v.Vars["token"]      = csrfbanana.Token(w, r, sess)
	v.Vars["Per"]        = per
        v.Vars["LisPeriod"]  = lisPeriod
        v.Vars["LisIngre"]    = lisIngre
        v.Vars["Level"]      =  sess.Values["level"]
	v.Render(w)
 }

//------------------------------------------------
// IngreDeleteGET handles the note deletion
 func IngreDeleteGET(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        var ingres model.IngresoN
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
	id,_        := atoi32(SId)
        path        :=  "/ingreso/list"
        ingres.Id   = id
	err         := (&ingres).IngresById()
	if err != nil { // Si no existe ingreso
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No hay ingreso.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "ingreso/ingresodelete"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     =  "Eliminar Ingreso"
        v.Vars["Action"]    =  "/ingreso/delete"
        v.Vars["Ingre"]        = ingres
        v.Vars["Level"]       =  sess.Values["level"]
	v.Render(w)
  }

// ---------------------------------------------------
// IngreUpPOST procesa la forma enviada con los datos
func IngreDeletePOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var ingres model.Ingreso
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        ingres.Id    = Id
        path        :=  "/ingreso/list"
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
             err = ingres.IngresDelete()
             if err != nil {
                 log.Println(err)
                 sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
              } else {
                  sess.AddFlash(view.Flash{"Ingreso borrado!", view.FlashSuccess})
              }
              sess.Save(r, w)
        }
	http.Redirect(w, r, path, http.StatusFound)
 }
// ---------------------------------------------------











