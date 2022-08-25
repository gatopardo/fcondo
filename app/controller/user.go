package controller

import (
	"log"
	"net/http"
         "strings"
         "fmt"
         "path"
          "os"
          "io"
	  "math/rand"
	  "time"
          "mime/multipart"

	"github.com/gatopardo/acad/ch14/fcondo/app/model"
	"github.com/gatopardo/acad/ch14/fcondo/app/shared/passhash"
	"github.com/gatopardo/acad/ch14/fcondo/app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )

 const (
	 maxUploadSize = 2 << 20 // 2 mb
         path_img = "/static/favicons/icon_"
 )

 type Tipo struct{
     Code string
     Stip  string
 }

 var LisTip =  [...]Tipo{{"D","Owner"},{"I","Tenant"},{"A","Admin"},{"R","Reference"}}
 var LisIcon = [...]string{"a", "b", "c", "d"}
// ---------------------------------------------------

  func saveImage(r *http.Request, sname, spath string)( err error){
         var file multipart.File
         var fileHeader *multipart.FileHeader
         r.ParseMultipartForm(maxUploadSize)
         file, fileHeader, err = r.FormFile(sname)
         if err != nil {
            log.Println(err)
             return
         }
         defer file.Close()

	 str :=  "."+path.Dir(spath)+"/"+fileHeader.Filename
         f, err := os.OpenFile(str, os.O_WRONLY|os.O_CREATE, 0666)
         if err != nil {
            log.Println(err)
             return
         }
         defer f.Close()
            var  nr int64
            nr, err =     io.Copy(f, file)
            if err != nil {
                  log.Println("No Copiando", nr, err)
                return
         }
         return
   }

   func getFileName(r * http.Request, sname string)(stf string, err error){
         var fileHeader *multipart.FileHeader
         r.ParseMultipartForm(maxUploadSize)
         _, fileHeader, err = r.FormFile(sname)
         if err != nil {
                  log.Println(err)
             return
         }

         sr  := rand.NewSource(time.Now().UnixNano())
         rn  := rand.New(sr)
	 ir  :=  rn.Intn(len(LisIcon))
         s2  :=  LisIcon[ir] + "/"
         stf   =  path_img+s2+fileHeader.Filename
	 return
   }

//----------------------------------------------------
 func getFormPers(p *model.Person, r * http.Request){
	 p.ApartaId,_     = atoi32(r.FormValue("aptId"))
         p.Fname          =  r.FormValue("fname")
         p.Lname          =  r.FormValue("lname")
         p.Email          =  r.FormValue("email")
         p.Address        =  r.FormValue("address")
         p.Tele           =  r.FormValue("tele")
         p.Mobil          =  r.FormValue("mobil")
         p.Tipo           =  r.FormValue("tipo")
         rPhoto,err       :=  getFileName(r, "photo")
	 if err == nil{
                 p.Photo          =  rPhoto
         }
   }

// --------------------------------------------------------
// RegisterGET despliega la pagina del usuario
 func RegisterGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess    := model.Instance(r)
	lsApts, err  := model.Apts()
	if err != nil {
//	   fmt.Println(err)
	   log.Println(err)
	}
	// Display the view
	v := view.New(r)
	v.Name = "register/register"
        v.Vars["LisApts"] = lsApts
        v.Vars["LisTips"] = LisTip
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
 }
// ---------------------------------------------------
	// Refill any form fields
//	view.Repopulate([]string{"cuenta", "password", "level"}, r.Form, v.Vars)
// ---------------------------------------------------

// RegisterPOST procesa la forma enviada con los datos
func RegisterPOST(w http.ResponseWriter, r *http.Request) {
        var user model.User
        var person model.Person
        var shad model.Shadow
	sess       := model.Instance(r)
        action     := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
           if validate, missingField := view.Validate(r, []string{"cuenta", "nivel", "password"}); !validate {
               sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
               sess.Save(r, w)
//               RegisterGET(w, r)
               return
	   }
           stUuid := model.CreateUUID()
           // Get form values
           rPasswd         := r.FormValue("password")
	   user.Cuenta      = r.FormValue("cuenta")
           vPasswd         := r.FormValue("password_verify")
           user.Nivel, _    = atoi32( r.FormValue("nivel"))
           user.Uuid        = stUuid
           if strings.Compare(rPasswd, vPasswd) != 0{
		log.Println(rPasswd, " * ", vPasswd)
		sess.AddFlash(view.Flash{"Claves distintas no posible", view.FlashError})
		sess.Save(r, w)
                RegisterGET(w, r)
//		http.Redirect(w, r, "/user/register", http.StatusFound)
		return
           }
            pass, errp := passhash.HashString(rPasswd)
	   if errp != nil {
		log.Println(errp)
                sess.AddFlash(view.Flash{"Problema encriptando clave.", view.FlashError})
		sess.Save(r, w)
//              RegisterGET(w, r)
		http.Redirect(w, r, "/user/register", http.StatusFound)
		return
	   }
           user.Password = pass
	   err := (&user).UserByCuenta()
           if err == model.ErrNoResult { // Exito:  no hay usuario creado aun 
                getFormPers(&person, r)
 // fmt.Println("Person ", person.Fname, person.Lname)
                err = (&person).PersonCreate()
                if err != nil {
                   log.Println(err)
                   sess.AddFlash(view.Flash{"Error guardando Person.", view.FlashError})
                   sess.Save(r, w)
		  http.Redirect(w, r, "/user/register", http.StatusFound)
                   return
                }
                saveImage(r, "photo", person.Photo )

                user.PersonId =  person.Id
                ex := (&user).UserCreate()
	        if ex != nil {  // uyy como fue esto ? 
                   log.Println(ex)
                   sess.AddFlash(view.Flash{"Error guardando User.", view.FlashError})
                   sess.Save(r, w)
		   http.Redirect(w, r, "/user/register", http.StatusFound)
                   return
	        } else {  // todo bien
                   shad.UserId    = user.Id
                   shad.Uuid      = stUuid
                   shad.Password  = rPasswd
                   if  err = (&shad).ShadCreate() ; err != nil{
                       sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
                       log.Println( err)
                       sess.Save(r, w)
		       http.Redirect(w, r, "/user/register", http.StatusFound)
                       return
                   }
                   sess.AddFlash(view.Flash{"Creando: " +user.Cuenta, view.FlashSuccess})
                   sess.Save(r, w)
	        }
           }
         }
	http.Redirect(w, r, "/user/list", http.StatusFound)
  }
// ---------------------------------------------------
// ---------------------------------------------------
// RegisUpGET despliega la pagina del usuario
func RegisUpGET(w http.ResponseWriter, r *http.Request) {
        var user model.User
        var pers model.Person
	sess := model.Instance(r)
	lsApts, err  := model.Apts()
	if err != nil {
//	   fmt.Println(err)
	   log.Println(err)
	}
        // necesitamos user id
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        id,_ := atoi32(params.ByName("id"))
        user.Id = id
        // Obtener usuario dado id
        err  = (&user).UserById()
        if err != nil { // Si no existe el usuario
            log.Println(err)
            sess.AddFlash(view.Flash{"Es raro. No usuario.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, "/user/list", http.StatusFound)
            return
        }
        pers.Id = user.PersonId
        err = (&pers).PersonById()
        if err != nil { // Si no existe el usuario
            log.Println(err)
            sess.AddFlash(view.Flash{"No hay atributos.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, "/user/list", http.StatusFound)
            return
        }
// fmt.Println("RegisUpGET uid ", id, " pid " , pers.Id," ",pers.Photo, " apt ", pers.ApartaId, " ", pers.Lname)
	// Display the view
	v := view.New(r)
	v.Name               = "register/regisupdate"
	v.Vars["token"]      = csrfbanana.Token(w, r, sess)
        v.Vars["User"]       = user
        v.Vars["Person"]     = pers
        v.Vars["LisTips"]    = LisTip
        v.Vars["LisApts"]    = lsApts
//    Refill any form fields
//	view.Repopulate([]string{"cuenta", "level"}, r.Form, v.Vars)
        v.Render(w)
   }
//---------------------------------------------------------------
//---------------------------------------------------------------
    func getPersFormUp(p1, p2 model.Person,r *http.Request )(stUp string){
      var sform string
      var  sArrSup []string
      if p1.ApartaId != p2.ApartaId  {
           sform = fmt.Sprintf(" %s = %d ","aparta_id", p2.ApartaId)
           sArrSup = append(sArrSup, sform)
      }
      if strings.Trim(p1.Fname, " ") !=  strings.Trim(p2.Fname," ") {
           sform = fmt.Sprintf(" %s  = '%s' ","fname",  strings.Trim(p2.Fname," ") )
           sArrSup = append(sArrSup, sform)
      }
      if strings.Trim(p1.Lname, " ") !=  strings.Trim(p2.Lname," ") {
           sform = fmt.Sprintf(" %s  = '%s' ","lname",  strings.Trim(p2.Lname," ") )
           sArrSup = append(sArrSup, sform)
      }
      if strings.Trim(p1.Address, " ") !=  strings.Trim(p2.Address," ") {
           sform = fmt.Sprintf(" %s  = '%s' ","address",  strings.Trim(p2.Address," ") )
           sArrSup = append(sArrSup, sform)
      }
      if strings.Trim(p1.Email, " ") !=  strings.Trim(p2.Email," ") {
           sform = fmt.Sprintf(" %s  = '%s' ","email",  strings.Trim(p2.Email," ") )
           sArrSup = append(sArrSup, sform)
      }
      if strings.Trim(p1.Tele, " ") !=  strings.Trim(p2.Tele," ") {
           sform = fmt.Sprintf(" %s  = '%s' ","tele",  strings.Trim(p2.Tele," ") )
           sArrSup = append(sArrSup, sform)
      }
      if strings.Trim(p1.Mobil, " ") !=  strings.Trim(p2.Mobil," ") {
           sform = fmt.Sprintf(" %s  = '%s' ","mobil",  strings.Trim(p2.Mobil," ") )
           sArrSup = append(sArrSup, sform)
      }
      if strings.Trim(p1.Tipo, " ") !=  strings.Trim(p2.Tipo," ") {
           sform = fmt.Sprintf(" %s  = '%s' ","tipo",  strings.Trim(p2.Tipo," ") )
           sArrSup = append(sArrSup, sform)
      }
      if (strings.Trim(p1.Photo, " ") !=  strings.Trim(p2.Photo," ")) &&
             (len(strings.Trim(p2.Photo, " ")) > 0){
           sform = fmt.Sprintf(" %s  = '%s' ","photo",  strings.Trim(p2.Photo," ") )
           sArrSup = append(sArrSup, sform)
	   saveImage(r, "photo",  p2.Photo)
      }

       lon := len(sArrSup)
       if lon  > 0 {
            now         := time.Now()
	    sf          :=  fmt.Sprintf( " updated_at = '%s' ", now.Format(layout) )
            sini        :=  "update persons set "
            stUp =  strings.Join(sArrSup, ", ")
            sr          :=  fmt.Sprintf(" where persons.id = %d ", p1.Id)
             stUp = sini + stUp +sf + sr
       }
         return
    }
//----------------------------------------------------
// RegisUpPOST procesa la forma enviada con los datos
func RegisUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var user model.User
        var pers, p model.Person
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
        action        := r.FormValue("action")
        path :=  "/user/list/"
        if (strings.Compare(action,"Cancelar") == 0) {
               sess.Save(r, w)
               http.Redirect(w, r, path, http.StatusFound)
               return
        }
	user.Id, _ = atoi32(params.ByName("id"))
        // Obtener usuario dado id
        err  = (&user).UserById()
        if err != nil { // Si no existe el usuario
            log.Println(err)
            sess.AddFlash(view.Flash{"Raro No usuario.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, "/user/list", http.StatusFound)
            return
        }
        pers.Id = user.PersonId
        err     = (&pers).PersonById()
        if err != nil { // Si no existe persona
            log.Println(err)
            sess.AddFlash(view.Flash{"No hay atributos.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, "/user/list", http.StatusFound)
            return
        }
	user.Cuenta     = r.FormValue("cuenta")
	getFormPers(&p,r)
        st          :=  getPersFormUp(pers,p, r)
//	fmt.Println(" RegisUpPOST ", st)
        if len(st) == 0{
            sess.AddFlash(view.Flash{"No actualizacion solicitada", view.FlashSuccess})
        } else {
            err =  pers.Update(st)
            if err == nil{
                 sess.AddFlash(view.Flash{"Persona actualizado: " +user.Cuenta, view.FlashSuccess})
            } else       {
		log.Println(err)
		sess.AddFlash(view.Flash{"Un error actualizando.", view.FlashError})
	    }
          }
		sess.Save(r, w)
	http.Redirect(w, r, path, http.StatusFound)
     }

// ---------------------------------------------------
// RegisSearchPOST procesa la forma enviada con los datos
func RegisSearchPOST(w http.ResponseWriter, r *http.Request) {
	sess       := model.Instance(r)
        rSearch    := r.FormValue("search")
        if rSearch == ""{
//           fmt.Println("Nada a buscar")
	   log.Println("Nada encontrado")
           return
         }
        lisUsers, err := model.SUsers(rSearch)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Usuarios.", view.FlashError})
            sess.Save(r, w)
         }
	// Display the view
	v := view.New(r)
	v.Name = "register/regislis"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
        v.Vars["Level"]    =  sess.Values["level"]
        v.Vars["LisRegis"] = lisUsers
	// Refill any form fields
	v.Render(w)
     }
// ---------------------------------------------------
//------------------------------------------------
// RegisLisGET displays the register page
func RegisLisGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
        lisUsers, err := model.Users()
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Usuarios.", view.FlashError})
            sess.Save(r, w)
         }
//	 fmt.Println(lisUsers)
	// Display the view
	v := view.New(r)
	v.Name = "register/regislis"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
        v.Vars["Level"]    =  sess.Values["level"]
        v.Vars["LisRegis"] = lisUsers
	v.Render(w)
}

// PersRegisGET displays the register page
func PersRegisGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "person/persreg"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	// Refill any form fields
	view.Repopulate([]string{"cuenta", "lname", "fname","email","address","tele", "mobil"  }, r.Form, v.Vars)
	v.Render(w)
}

// PersRegisPOST maneja registro personas  con  form submission
func  PersRegisPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
        var user model.User
        var person model.Person
	sess := model.Instance(r)
	// Se previene ataque de fuerza bruta con entradas invalidas múltiples :-)
	if sess.Values["register_attempt"] != nil && sess.Values["register_attempt"].(int) >= 5 {
		log.Println("Detenido intento en persona entrada repetida multiples veces")
		http.Redirect(w, r, "/regispers", http.StatusFound)
		sess.AddFlash(view.Flash{"Faltan Campos", view.FlashError})
		return
	}
	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"cuenta", "lname", "fname", "email", "address", "tele", "mobil" }); !validate {
		sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
		sess.Save(r, w)
		RegisterGET(w, r)
		return
	}
	// Get form values
	person.Fname   =  r.FormValue("fname")
	person.Lname   =  r.FormValue("lname")
	person.Email   =  r.FormValue("email")
	person.Address =  r.FormValue("address")
	person.Tele    =  r.FormValue("tele")
	person.Mobil   =  r.FormValue("mobil")
	// Busquemos datos de cuenta
	 err := user.UserByCuenta()
	if err == nil { // Exito  usuario con esa cuenta
	    ex := (&person).PersonCreate() // Si hay error es por el query
            if ex == nil { // todo bien creando persona
		sess.AddFlash(view.Flash{"Persona creada exitosamente para: " +user.Cuenta, view.FlashSuccess})
                sess.Save(r, w)
                http.Redirect(w, r, "/regispers", http.StatusFound)
		return
             }else{ // persona ya existía
                    sess.AddFlash(view.Flash{"La cuenta ya existe para: " + user.Cuenta, view.FlashError})
                sess.Save(r, w)
             }
        } else { // no hay usuario para esta persona    
	           log.Println(err)
		   sess.AddFlash(view.Flash{"Vamos a Crear Usuario.", view.FlashError})
		   sess.Save(r, w)
	}
//	if err == model.ErrNoResult  // If success (no usuario con esa cuenta)
	// Display the page
	PersRegisGET(w, r)
}

// PersLisGET displays the register page
func PersLisGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
        lisPers, err := model.Persons()
        if err != nil {
	    log.Println(err)
	    sess.AddFlash(view.Flash{"Error ocurrio Listando Personas.", view.FlashError})
             sess.Save(r, w)
         }
	// Display the view
	v := view.New(r)
	v.Name = "person/perslis"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
        v.Vars["lisPers"] = lisPers
	v.Render(w)
}
//
// UserDelGET handles the user deletion
 func RegisDelGET(w http.ResponseWriter, r *http.Request) {
	// Get session
        sess := model.Instance(r)
        var params httprouter.Params
        params =  context.Get(r, "params").(httprouter.Params)
        Id,_   := atoi32(params.ByName("id"))
        path   := "/user/list"
        var user model.User
        var pers model.Person
	var apt  model.Aparta
        user.Id = Id
        err := (&user).UserById()
        if err != nil {
            log.Println(err)
//            fmt.Println(err)
            sess.AddFlash(view.Flash{"Error Usuario no hallado.", view.FlashError})
             http.Redirect(w, r, path, http.StatusFound)
            return
        }
        pers.Id = user.PersonId
        err = (&pers).PersonById()
        if err != nil { // Si no existe el usuario
            log.Println(err)
            sess.AddFlash(view.Flash{"No hay atributos.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, "/user/list", http.StatusFound)
            return
        }
	apt.Id  = pers.ApartaId
        err =   (&apt).AptById()
        if err != nil {
            log.Println(err)
//            fmt.Println(err)
            sess.AddFlash(view.Flash{"Error Usuario no hallado.", view.FlashError})
             http.Redirect(w, r, path, http.StatusFound)
            return
        }

	v := view.New(r)
        v.Name = "register/regisdelete"
	v.Name               = "register/regisupdate"
	v.Vars["token"]      = csrfbanana.Token(w, r, sess)
        v.Vars["User"]       = user
        v.Vars["aptCodigo"]  = apt.Codigo
        v.Vars["Person"]     = pers
        v.Render(w)
  }

// PersDelPOST handles the user deletion
 func RegisDelPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var pers model.Person
	sess := model.Instance(r)
        var params httprouter.Params
        params       = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        pers.Id      = Id
        path        :=  "/user/list"
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
             err = pers.PersDeleteById()
             if err != nil {
                 log.Println(err)
                 sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
              } else {
                  sess.AddFlash(view.Flash{"Persona borrada!", view.FlashSuccess})
              }
              sess.Save(r, w)
        }
	http.Redirect(w, r, path, http.StatusFound)

 }




