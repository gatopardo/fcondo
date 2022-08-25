package model

import (
        "database/sql"
	"time"
//	"fmt"
        "log"


)

// *****************************************************************************
// Aparta
// *****************************************************************************

// Apartas table contains the information for each apartment
type Aparta struct {
	Id            uint32        `db:"id" bson:"id,omitempty"`
	Codigo        string        `db:"codigo" bson:"codigo"`
	Descripcion   string        `db:"descricpion" bson:"descripcion"`
	CreatedAt     time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at" bson:"updated_at"`
}

   type  Apt string

// --------------------------------------------------------

// AptById tenemos el apartamento dado id
func (apt * Aparta)AptById() (err error) {
        stq  :=   "SELECT id, codigo, descripcion, created_at, updated_at FROM apartas WHERE id=$1"
	err = Db.QueryRow(stq, &apt.Id).Scan(&apt.Id, &apt.Codigo, &apt.Descripcion, &apt.CreatedAt, &apt.UpdatedAt)

	return  standardizeError(err)
}

// --------------------------------------------------------

// AptById tenemos el apartamento dado id
func (apt * Aparta)AptByCode() (err error) {
        stq  :=   "SELECT id, codigo, descripcion, created_at, updated_at FROM apartas WHERE codigo=$1"
	err = Db.QueryRow(stq, &apt.Codigo).Scan(&apt.Id, &apt.Codigo, &apt.Descripcion, &apt.CreatedAt, &apt.UpdatedAt)

	return  standardizeError(err)
}

// -----------------------------------------------------
// AptCreate crear apartamento
func (apt *Aparta)AptCreate() error {
         var err error
         var stmt  *sql.Stmt
         stq := "INSERT INTO apartas ( codigo, descripcion, created_at, updated_at ) VALUES ($1,$2,$3, $4) returning id" 

	now  := time.Now()

            if stmt, err = Db.Prepare(stq ); err != nil  {
	          return standardizeError(err)
              }
             defer stmt.Close()
             var id uint32
             err = stmt.QueryRow(  &apt.Codigo, &apt.Descripcion,  now, now ).Scan(&id)
             if err == nil {
                apt.Id = id
             }
// fmt.Println("Apt Creado")
	return standardizeError(err)
}

// -----------------------------------------------------
 func  (apt * Aparta)AptDeleteById()( err error){
         stqd :=  "DELETE FROM apartas where id = $1"
           _, err = Db.Exec(stqd, apt.Id)
         return
       }

// Delete apartamento from databa
func (apt *Aparta)AptDelete() (err error) {
	statement := "delete from apartas where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(apt.Id)
	return
}

// -----------------------------------------------------
// Actualizar informacion de apartamento en la database
func (apt *Aparta)AptUpdate(stq string) (err error) {
        _, err = Db.Exec(stq )
        return standardizeError(err)
 }
// -----------------------------------------------------
// Delete all apartamento from database
func AptDeleteAll() (err error) {
	statement := "delete from apartas"
	_, err = Db.Exec(statement)
	return
}

// -------------------------------------------------------------
// Get number of records in apartas
  func AptCount( ) ( count int) {
        stq :=  "SELECT COUNT(*) as count FROM apartas "
	rows, err := Db.Query(stq)
	if err != nil {
		return
	}
	defer rows.Close()
        for rows.Next() {
            err = rows.Scan(&count)
	    if err != nil {
	         return
	    }
        }
	return
 }
// -------------------------------------------------------------
// Get limit records from offset
  func AptLim(lim int , offs int) ( apartas []Aparta, err error) {
        var apt Aparta
        stq :=   "SELECT id, codigo, descripcion, created_at, updated_at FROM apartas order by codigo LIMIT $1 OFFSET $2"

	rows, err := Db.Query(stq, lim, offs)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
                apt = Aparta{}
		if err = rows.Scan(&apt.Id,  &apt.Codigo, &apt.Descripcion, &apt.CreatedAt, &apt.UpdatedAt); err != nil {
			return
		}
		apartas = append(apartas, apt)
	}
	return
 }
// -------------------------------------------------------------
// Get all apartas in the database and returns the list
  func Apts() (apartas []Aparta, err error) {
        var apt Aparta
        stq :=   "SELECT id,  codigo, descripcion, created_at, updated_at FROM apartas order by  codigo"
	rows, err := Db.Query(stq)
	if err != nil {

		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&apt.Id, &apt.Codigo, &apt.Descripcion, &apt.CreatedAt, &apt.UpdatedAt); err != nil {
		return
		}
		apartas = append(apartas, apt)
	}
	return
 }
// -------------------------------------------------------------
// Get all apartas in the database and returns index and code 
  func IApts() (apartas []Aparta, err error) {
        var apt Aparta
        stq :=   "SELECT id,  codigo FROM apartas order by  codigo"
	rows, err := Db.Query(stq)
	if err != nil {

		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&apt.Id, &apt.Codigo) ; err != nil {
		return
		}
		apartas = append(apartas, apt)
	}
	return
 }
// -------------------------------------------------------------
// Given user.Id get aparta
    func ApartaByUserId(id uint32)(p Person, a Aparta, err error){

        stq := "SELECT a.id, a.codigo, a.descripcion, p.fname, p.lname, p.email, p.mobil, p.tipo,p.photo from users u join persons p on u.person_id = p.id join apartas a on p.aparta_id = a.id  where u.id = $1"

	row := Db.QueryRow(stq, id)
        err = row.Scan(&a.Id, &a.Codigo, &a.Descripcion, &p.Fname, &p.Lname, &p.Email,&p.Mobil, &p.Tipo, &p.Photo)
	if err != nil {
                  if err ==  sql.ErrNoRows {
                     log.Println("No rows were returned!")
		   }else{
                     log.Println(err)
		   }
             }
     return
    }
// -------------------------------------------------------------

 func LisApts()(lsApts []Apt,err error){
        stq := "select codigo from Apartas order by codigo"
        var rows * sql.Rows
	rows, err = Db.Query(stq)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()
        for rows.Next() {
            var c   Apt
            err = rows.Scan(&c)
	    lsApts = append(lsApts, c)
        }
	return
  }
// -------------------------------------------------------------

