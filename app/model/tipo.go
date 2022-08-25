package model

import (
        "database/sql"
	"time"
//	"fmt"
//        "log"


)

// *****************************************************************************
// Tipo
// *****************************************************************************

// Tipo table contains the information for each category
type Tipo struct {
	Id            uint32        `db:"id" bson:"id,omitempty"`
	Codigo        string        `db:"codigo" bson:"codigo"`
	Aplica        string        `db:"aplica" bson:"aplica"`
	Descripcion   string        `db:"descricpion" bson:"descripcion"`
	CreatedAt     time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at" bson:"updated_at"`
}

type Aplic struct {
	Id            uint32        `db:"id" bson:"id,omitempty"`
	Aplica        string        `db:"aplica" bson:"aplica"`
}

// --------------------------------------------------------

func  ApliLis()(  ptApli []Aplic){
      ptApli = make([]Aplic, 2)
      ptApli[0]. Id = 0
      ptApli[0].Aplica = "DB"
      ptApli[1]. Id = 0
      ptApli[1].Aplica = "CR"
      return
}

// TipoById tenemos el category dado id
func (tip * Tipo)TipoById() (err error) {
        stq  :=   "SELECT id, codigo, aplica, descripcion, created_at, updated_at FROM tipos WHERE id=$1"
	err = Db.QueryRow(stq, &tip.Id).Scan(&tip.Id, &tip.Codigo, &tip.Descripcion, &tip.CreatedAt, &tip.UpdatedAt)

	return  standardizeError(err)
}

// --------------------------------------------------------

// TipById tenemos el category dado id
func (tip * Tipo)TipoByCode() (err error) {
        stq  :=   "SELECT id, codigo, aplica, descripcion, created_at, updated_at FROM tipos WHERE codigo=$1"
	err = Db.QueryRow(stq, &tip.Codigo).Scan(&tip.Id, &tip.Codigo,&tip.Aplica, &tip.Descripcion, &tip.CreatedAt, &tip.UpdatedAt)

	return  standardizeError(err)
}

// -----------------------------------------------------
// TipCreate crear category
func (tip *Tipo)TipoCreate() error {
         var err error
         var stmt  *sql.Stmt
         stq := "INSERT INTO tipos ( codigo, aplica, descripcion, created_at, updated_at ) VALUES ($1,$2,$3, $4) returning id"

	now  := time.Now()

            if stmt, err = Db.Prepare(stq ); err != nil  {
	          return standardizeError(err)
              }
             defer stmt.Close()
             var id uint32
             err = stmt.QueryRow(  &tip.Codigo, &tip.Descripcion,  now, now ).Scan(&id)
             if err == nil {
                tip.Id = id
             }
	return standardizeError(err)
}

// -----------------------------------------------------
 func  (tip * Tipo)TipoDeleteById()( err error){
         stqd :=  "DELETE FROM tipos where id = $1"
           _, err = Db.Exec(stqd, tip.Id)
         return
       }

// Delete category from databa
func (tip *Tipo)Delete() (err error) {
	statement := "delete from tipos where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(tip.Id)
	return
}

// -----------------------------------------------------
// Actualizar informacion de category en la database
func (tip *Tipo)Update(stq string) (err error) {
	stmt, err := Db.Prepare(stq)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(tip.Id, tip.Codigo, tip.Descripcion)
	return
}

// -----------------------------------------------------
// Delete all category from database
func TipoDeleteAll() (err error) {
	statement := "delete from tipos"
	_, err = Db.Exec(statement)
	return
}

// -------------------------------------------------------------
// Get number of records in tipos
  func TipoCount( ) ( count int) {
        stq :=  "SELECT COUNT(*) as count FROM tipos "
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
  func TipoLim(lim int , offs int) ( tipos []Tipo, err error) {
        var tip Tipo
        stq :=   "SELECT id, codigo, aplica, descripcion, created_at, updated_at FROM tipos order by codigo LIMIT $1 OFFSET $2"
	rows, err := Db.Query(stq, lim, offs)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
                tip = Tipo{}
		if err = rows.Scan(&tip.Id,  &tip.Codigo,&tip.Aplica, &tip.Descripcion, &tip.CreatedAt, &tip.UpdatedAt); err != nil {
			return
		}
		tipos = append(tipos, tip)
	}
	return
 }
// -------------------------------------------------------------
// Get all tipos in the database and returns the list
  func Tipos() (tipos []Tipo, err error) {
        var tip Tipo

        stq :=   "SELECT id,  codigo,aplica, descripcion, created_at, updated_at FROM tipos   order by  codigo"
	rows, err := Db.Query(stq)
	if err != nil {

		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&tip.Id, &tip.Codigo,&tip.Aplica, &tip.Descripcion, &tip.CreatedAt, &tip.UpdatedAt); err != nil {
		return
		}
		tipos = append(tipos, tip)
	}
	return
 }
// -------------------------------------------------------------
// Get all tipos in the database and returns the list
  func TiposI() (tipos []Tipo, err error) {
        var tip Tipo
	stq := " select id,codigo, aplica, descripcion, created_at, updated_at  from tipos where  aplica = 'CR' and not  (codigo = 'NC' or codigo = 'IC') order by codigo"
	rows, err := Db.Query(stq)
	if err != nil {

		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&tip.Id, &tip.Codigo,&tip.Aplica, &tip.Descripcion, &tip.CreatedAt, &tip.UpdatedAt); err != nil {
		return
		}
		tipos = append(tipos, tip)
	}
	return
 }
// -------------------------------------------------------------
