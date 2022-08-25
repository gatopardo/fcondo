package model

import (
        "database/sql"
	"time"
        "log"
//	"fmt"
)

// *****************************************************************************
// Ingreso
// *****************************************************************************.

// Ingreso table contains the information for each ingres
type Ingreso struct {
	Id               uint32     `db:"id" bson:"id,omitempty"`
	PeriodId         uint32     `db:"periodid" bson:"periodid,omitempty"`
	TipoId           uint32     `db:"tipoid" bson:"tipoid,omitempty"`
        Fecha       time.Time       `db:"fecha" bson:"fecha"`
        Amount           int64     `db:"amount" bson:"amount"`
        Descripcion      string      `db:"dscripcion" bson:"dscripcion"`
	CreatedAt   time.Time       `db:"created_at" bson:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" bson:"updated_at"`
}

type IngresoN struct {
	Id               uint32     `db:"id" bson:"id,omitempty"`
	PeriodId         uint32     `db:"periodid" bson:"periodid,omitempty"`
        Period      time.Time       `db:"period" bson:"period"`
	TipoId           uint32     `db:"tipoid" bson:"tipoid,omitempty"`
	Tipo             string     `db:"tcodigo" bson:"tcodigo,omitempty"`
        Fecha       time.Time       `db:"fecha" bson:"fecha"`
        Amount           int64     `db:"amount" bson:"amount"`
        Descripcion      string      `db:"dscripcion" bson:"dscripcion"`
	CreatedAt time.Time        `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time        `db:"updated_at" bson:"updated_at"`
}

type IngresoJ struct {
	Tipo             string     `db:"tcodigo" bson:"tcodigo,omitempty"`
        Fecha       time.Time       `db:"fecha" bson:"fecha"`
        Amount           int64     `db:"amount" bson:"amount"`
        Descripcion      string      `db:"dscripcion" bson:"dscripcion"`
}

type IngresoL struct {
        Period      time.Time
	LisIngre     []IngresoJ
}

// -----------------------------------------
// IngresById tenemos el ingreso dado id
func (ingres * IngresoN)IngresById() (err error) {
	stq :=  " SELECT i.id, i.period_id, p.inicio, i.tipo_id, t.codigo, i.fecha, i.amount,  i.description, i.created_at, i.updated_at  FROM ingresos i JOIN periods p ON i.period_id = p.id  JOIN  tipos t ON i.tipo_id = t.id WHERE  i.id=$1 "

		err = Db.QueryRow(stq, &ingres.Id). Scan(&ingres.Id, &ingres.PeriodId,&ingres.Period,  &ingres.TipoId, &ingres.Tipo, &ingres.Fecha, &ingres.Amount, &ingres.Descripcion,  &ingres.CreatedAt, &ingres.UpdatedAt)

	return  standardizeError(err)
}

// -----------------------------------------------------
// IngresCreate crear ingreso
func (ing * IngresoN)IngresCreate() error {
         var err error
         var stmt  *sql.Stmt
         stq := "INSERT INTO ingresos ( period_id, tipo_id, fecha, amount, description, created_at, updated_at ) VALUES ($1,$2,$3,$4, $5, $6, $7) returning id"
	 now  := time.Now()
         if stmt, err = Db.Prepare(stq ); err != nil  {
	          return standardizeError(err)
         }
         defer stmt.Close()
         var id uint32
         err = stmt.QueryRow(&ing.PeriodId, &ing.TipoId,&ing.Fecha, &ing.Amount, &ing.Descripcion,  now, now ).Scan(&id)
         if err == nil {
              ing.Id = id
         }
	 return standardizeError(err)
  }

// -----------------------------------------------------
 func  (ingres * Ingreso)IngresDeleteById()( err error){
         stqd :=  "DELETE FROM ingresos where id = $1"
           _, err = Db.Exec(stqd, ingres.Id)
        return
       }

// -----------------------------------------------------
// Delete ingres from databa
func (ingres *Ingreso) IngresDelete() (err error) {
	statement := "delete from ingresos where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(ingres.Id)
	return
}

// -----------------------------------------------------
// Actualizar informacion de ingres en la database
func (ingres *IngresoN)IngresUpdate(stq string) (err error) {
        _, err = Db.Exec(stq )
        return standardizeError(err)
}

// -----------------------------------------------------
// Delete all ingresos from database
func IngresDeleteAll() (err error) {
	statement := "delete from ingresos"
	_, err = Db.Exec(statement)
	return
}

// -------------------------------------------------------------
// Get number of records in ingresos
  func IngresCount( ) ( count int) {
        stq :=  "SELECT COUNT(*) as count FROM ingresos "
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
// Get ingresos from a period 
  func IngresLim(id uint32 ) (ingresos []IngresoN, err error) {
      stq :=  "SELECT i.id, i.period_id, p.inicio,  i.tipo_id, t.codigo, i.fecha, i.amount, i.description, i.created_at, i.updated_at FROM ingresos i JOIN  periods p ON i.period_id = p.id JOIN tipos t ON i.tipo_id = t.id WHERE   p.id = $1 ORDER BY p.inicio, i.fecha "

	rows, err := Db.Query(stq, id)
	if err != nil {
            return
	}
        defer rows.Close()
        for rows.Next() {
           ingres := IngresoN{}
           if err = rows.Scan(&ingres.Id,&ingres.PeriodId,&ingres.Period, &ingres.TipoId, &ingres.Tipo, &ingres.Fecha, &ingres.Amount, &ingres.Descripcion,  &ingres.CreatedAt, &ingres.UpdatedAt); err != nil {
                  return
            }
           ingresos = append(ingresos, ingres)
         }
       return
 }
// -------------------------------------------------------------
// Get ingresos from a period - json 
  func IngresoJPer(id uint32 ) (ingresos []IngresoJ, err error) {
	stLayout := "2006-01-02"
        stq :=   "SELECT t.codigo, e.fecha, e.amount, e.description FROM ingresos e, periods p,  tipos t where e.period_id = p.id  and e.tipo_id = t.id and p.id = $1 order by p.inicio, e.fecha "
	rows, err := Db.Query(stq, id)
	if err != nil {
            return
	}
        defer rows.Close()
        for rows.Next() {
	    var  sqFec  sql.NullTime
	    var  sqAmt  sql.NullInt64

           ingres := IngresoJ{}
           if err = rows.Scan( &ingres.Tipo, &sqFec, &sqAmt, &ingres.Descripcion); err != nil {
		   log.Println(err)
                  return
            }
           if sqFec.Valid{
		 ingres.Fecha = sqFec.Time
	    }else{
		  ingres.Fecha, _ = time.Parse(stLayout, "1900-01-01")
	  }
	    if sqAmt.Valid{
		    ingres.Amount = sqAmt.Int64
	    }else{
		    ingres.Amount = 0
	    }
           ingresos = append(ingresos, ingres)
         }
       return
 }

// -------------------------------------------------------------
// -------------------------------------------------------------
// Get all ingresos per a period in the database and returns the list
  func (ingre * IngresoN)IngresPer() (ingresos []IngresoN, err error) {
        stq :=   "SELECT i.id, i.period_id, p.inicio,  i.tipo_id, t.codigo, i.fecha, i.amount, i.description, i.created_at, i.updated_at FROM ingresos i, periods p,  tipos t where i.period_id = p.id and i.tipo_id = t.id and i.period_id = $1 order by p.inicio"
	rows, err := Db.Query(stq, ingre.PeriodId )
	if err != nil {
            return
	}
	defer rows.Close()
	for rows.Next() {
            i := IngresoN{}
           if err = rows.Scan(&i.Id, &i.PeriodId, &i.Period, &i.TipoId, &i.Tipo, &i.Fecha, &i.Amount, &i.Descripcion, &i.CreatedAt, &i.UpdatedAt); err != nil {
                  return
             }
             ingresos = append(ingresos, i)
	}
        return
 }
// -------------------------------------------------------------
// Get all ingresos in the database and returns the list
  func Ingress() (ingresos []IngresoN, err error) {
        stq :=   "SELECT i.id, i.period_id, p.inicio,  a.codigo, i.tipo_id, t.codigo, i.fecha, i.amount, i.description, i.created_at, i.updated_at FROM ingresos i, periods p,  tipos t where i.period_id = p.id  and i.Tipo_id = t.id order by p.inicio"
	rows, err := Db.Query(stq)
	if err != nil {
            return
	}
	defer rows.Close()
	for rows.Next() {
            ingres := IngresoN{}
           if err = rows.Scan(&ingres.Id,&ingres.PeriodId,&ingres.Period, &ingres.TipoId, &ingres.Tipo,  &ingres.Fecha,&ingres.Amount, &ingres.Descripcion,  &ingres.CreatedAt, &ingres.UpdatedAt); err != nil {
                  return
             }
             ingresos = append(ingresos, ingres)
	}
        return
 }
// -------------------------------------------------------------
