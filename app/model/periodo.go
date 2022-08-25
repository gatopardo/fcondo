package model

import (
        "database/sql"
	"time"
//	"fmt"
//        "log"


)

// *****************************************************************************
// Periodo
// *****************************************************************************

// Periods table contains the information for each period
type Periodo struct {
	Id            uint32        `db:"id" bson:"id,omitempty"`
	Inicio        time.Time     `db:"inicio" bson:"inicio"`
	Final         time.Time     `db:"final" bson:"final"`
	CreatedAt     time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at" bson:"updated_at"`
}


// --------------------------------------------------------

// PeriodById tenemos el period dado id
func (period * Periodo)PeriodById() (err error) {
        stq  :=   "SELECT id, inicio, final, created_at, updated_at FROM periods WHERE id=$1"
	err = Db.QueryRow(stq, &period.Id).Scan(&period.Id, &period.Inicio, &period.Final, &period.CreatedAt, &period.UpdatedAt)

	return  standardizeError(err)
}

// --------------------------------------------------------
// PeriodByCode tenemos el period dado inicio
func (period * Periodo)PeriodByCode() (err error) {
        stq  :=   "SELECT id, inicio, final, created_at, updated_at FROM periods WHERE inicio = $1"
	err = Db.QueryRow(stq, &period.Inicio).Scan(&period.Id, &period.Inicio, &period.Final, &period.CreatedAt, &period.UpdatedAt)

	return  standardizeError(err)
}

// --------------------------------------------------------
// PeriodByFec tenemos el period dado id
func (period * Periodo)PeriodByFec(fec time.Time) (err error) {
        stq  :=   "SELECT id, inicio, final, created_at, updated_at FROM periods WHERE $1 BETWEEN inicio AND final"
	err = Db.QueryRow(stq, fec).Scan(&period.Id, &period.Inicio, &period.Final, &period.CreatedAt, &period.UpdatedAt)

	return  standardizeError(err)
}

// -----------------------------------------------------
// PeriodCreate crear period
func (period *Periodo)PeriodCreate() error {
         var err error
         var stmt  *sql.Stmt
         stq := "INSERT INTO periods ( inicio, final, created_at, updated_at ) VALUES ($1,$2,$3, $4) returning id"

	now  := time.Now()

            if stmt, err = Db.Prepare(stq ); err != nil  {
	          return standardizeError(err)
              }
             defer stmt.Close()
             var id uint32
             err = stmt.QueryRow(  &period.Inicio, &period.Final,  now, now ).Scan(&id)
             if err == nil {
                period.Id = id
             }
// fmt.Println("Period Creado")
	return standardizeError(err)
}

// -----------------------------------------------------
 func  (period * Periodo)PeriodDeleteById()( err error){
         stqd :=  "DELETE FROM periods where id = $1"
           _, err = Db.Exec(stqd, period.Id)
         return
       }

// -----------------------------------------------------
// Delete period from databa
func (period *Periodo)PeriodDelete() (err error) {
	statement := "delete from periods where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(period.Id)
	return
}

// -----------------------------------------------------
// Actualizar informacion de period en la database
func (period *Periodo)PeriodUpdate(stq string) (err error) {
        _, err = Db.Exec(stq )
        return standardizeError(err)
}

// -----------------------------------------------------
// Delete all period from database
func PeriodDeleteAll() (err error) {
	statement := "delete from periods"
	_, err = Db.Exec(statement)
	return
}

// -------------------------------------------------------------
// Get number of records in periods
  func PeriodCount( ) ( count int) {
        stq :=  "SELECT COUNT(*) as count FROM periods "
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
  func PeriodLim(lim int , offs int) ( periods []Periodo, err error) {
        var period Periodo
        stq :=   "SELECT id, inicio, final, created_at, updated_at FROM periods order by inicio LIMIT $1 OFFSET $2"

	rows, err := Db.Query(stq, lim, offs)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
                period = Periodo{}
		if err = rows.Scan(&period.Id,  &period.Inicio, &period.Final, &period.CreatedAt, &period.UpdatedAt); err != nil {
			return
		}
		periods = append(periods, period)
	}
	return
 }
// -------------------------------------------------------------
// Get all periods in the database and returns the list
  func Periods() (periods []Periodo, err error) {
        var period Periodo
        stq :=   "SELECT id,  inicio, final, created_at, updated_at FROM periods order by  inicio desc"
	rows, err := Db.Query(stq)
	if err != nil {

		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&period.Id, &period.Inicio, &period.Final, &period.CreatedAt, &period.UpdatedAt); err != nil {
		return
		}
		periods = append(periods, period)
	}
	return
 }
// -------------------------------------------------------------
// Get all periods with ingress in the database and returns the list
  func PeriodsI() (periods []Periodo, err error) {
        var period Periodo
	stq := "select p.id, p.inicio, p.final, p.created_at, p.updated_at  from periods p join (select i.period_id, sum(i.amount) from ingresos i group by i.period_id having sum(i.amount) > 0) b on p.id  = b.period_id order by p.inicio desc"
	rows, err := Db.Query(stq)
	if err != nil {

		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&period.Id, &period.Inicio, &period.Final, &period.CreatedAt, &period.UpdatedAt); err != nil {
		return
		}
		periods = append(periods, period)
	}
	return
 }
// -------------------------------------------------------------
// dbcondo=# select  i.fecha, i.amount, i.description from ingresos i join (select p.id, sum(i.amount) from periods p join ingresos i on  p.id = i.period_id  group by p.id having sum(i.amount) > 0 ) b on i.period_id = b.id ;

