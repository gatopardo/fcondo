package model

import (
        "database/sql"
	"time"
        "log"
//        "fmt"

)

// *****************************************************************************
// Cuota
// *****************************************************************************.

// Cuota table contains the information for each cuot
type Cuota struct {
	Id               uint32     `db:"id" bson:"id,omitempty"`
	PeriodId         uint32     `db:"periodid" bson:"periodid,omitempty"`
	ApartaId         uint32     `db:"apartaid" bson:"apartaid,omitempty"`
	TipoId           uint32     `db:"tipoid" bson:"tipoid,omitempty"`
        Fecha       time.Time       `db:"fecha" bson:"fecha"`
        Amount           int64     `db:"amount" bson:"amount"`
	CreatedAt   time.Time       `db:"created_at" bson:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" bson:"updated_at"`
}

type CuotaN struct {
	Id               uint32     `db:"id" bson:"id,omitempty"`
	PeriodId         uint32     `db:"periodid" bson:"periodid,omitempty"`
        Period        time.Time     `db:"period" bson:"period"`
	ApartaId         uint32     `db:"apartaid" bson:"apartaid,omitempty"`
	Apto             string     `db:"acodigo" bson:"acodigo,omitempty"`
	TipoId           uint32     `db:"tipoid" bson:"tipoid,omitempty"`
	Tipo             string     `db:"tdescripcion" bson:"tdescripcion,omitempty"`
        Fecha         time.Time     `db:"fecha" bson:"fecha"`
        Amount           int64     `db:"amount" bson:"amount"`
	CreatedAt     time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at" bson:"updated_at"`
}

  type  CuotApt  struct {
	 Id               uint32
         Inicio        time.Time
	 Cuota             int64
         Fecha         time.Time
	 Amount             int64
	 Balance           int64
  }

  type  AmtCond  struct {
	 Codigo        string
	 Amount        int64
	 Atraso        int64
	 Mora          int64
  }

  type AmtPerCond struct {
          Fecha         time.Time
	  LisAmt        []AmtCond
  }

   type  AmtAptTot struct{
         Inicio  time.Time
	 Final    time.Time
	 Tp      string
	 Ap      string
	 Cuota   int64
         Fecha  time.Time
	 Monto   int64
	 Dif     int64
}

// "SELECT c.id, c.period_id,p.inicio,c.aparta_id,a.codigo,c.tipo_id, t.descripcion, c.fecha, c.amount, c.created_at, c.updated_at FROM cuotas c, periods p, apartas a, tipos t  WHERE c.period_id = p.id and c.aparta_id = a.id AND c.tipo_id = t.id  AND c.id= $1 "

// -----------------------------------------
// CuotById tenemos la cuota dado id
func (cuot * CuotaN)CuotById() (err error) {
        stq  := " SELECT c.id, c.period_id,p.inicio,c.aparta_id,a.codigo,c.tipo_id," +
	        " t.descripcion, c.fecha, c.amount, c.created_at, c.updated_at " +
		" FROM cuotas c JOIN periods p ON c.period_id = p.id " +
		" JOIN  apartas a ON c.aparta_id = a.id " +
		" JOIN tipos t  ON c.tipo_id = t.id  WHERE c.id= $1 "

	err = Db.QueryRow(stq, &cuot.Id). Scan(&cuot.Id, &cuot.PeriodId,&cuot.Period,&cuot.ApartaId, &cuot.Apto, &cuot.TipoId, &cuot.Tipo, &cuot.Fecha, &cuot.Amount,  &cuot.CreatedAt, &cuot.UpdatedAt)

	return  standardizeError(err)
}

// -----------------------------------------------------
// CuotCreate crear cuota
func (c * CuotaN)CuotCreate() error {
         var err error
         var stmt  *sql.Stmt
         stq := "INSERT INTO cuotas ( period_id,aparta_id, tipo_id, fecha, amount,  created_at, updated_at ) VALUES ($1,$2,$3,$4, $5, $6, $7) returning id"
	 now  := time.Now()
         if stmt, err = Db.Prepare(stq ); err != nil  {
	          return standardizeError(err)
         }
         defer stmt.Close()
         var id uint32
         err = stmt.QueryRow(&c.PeriodId,&c.ApartaId, &c.TipoId,&c.Fecha, &c.Amount,  now, now ).Scan(&id)
         if err == nil {
              c.Id = id
         }
	 return standardizeError(err)
 }

// -----------------------------------------------------
 func  (cuot * Cuota)CuotDeleteById()( err error){
         stqd :=  "DELETE FROM cuotas where id = $1"
           _, err = Db.Exec(stqd, cuot.Id) 
         return
       }

// -----------------------------------------------------
// Delete cuot from databa
func (cuot *Cuota) CuotDelete() (err error) {
	statement := "delete from cuotas where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(cuot.Id)
	return
}

// -----------------------------------------------------
// Actualizar informacion de cuota en la database
func (cuot *CuotaN)CuotUpdate(stq string) (err error) {
        _, err = Db.Exec(stq )
        return standardizeError(err)
}

// -----------------------------------------------------
// Delete all cuotas from database
func CuotDeleteAll() (err error) {
	statement := "delete from cuotas"
	_, err = Db.Exec(statement)
	return
}

// -------------------------------------------------------------
// Get number of records in cuotas
  func CuotsCount( ) ( count int) {
        stq :=  "SELECT COUNT(*) as count FROM cuotas "
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

//	"SELECT c.id, c.period_id, p.inicio, c.aparta_id, a.codigo, c.tipo_id, t.descripcion, c.fecha, c.amount, c.created_at, c.updated_at FROM cuotas c, periods p, apartas a, tipos t where c.period_id = p.id and c.aparta_id = a.id and c.tipo_id = t.id and p.id = $1 order by p.inicio, c.fecha "
// -------------------------------------------------------------
// Get cuotas from a period 
  func CuotLim(id uint32 ) (cuotas []CuotaN, err error) {
        stq :=  " SELECT c.id, c.period_id, p.inicio, c.aparta_id, a.codigo, c.tipo_id, t.descripcion, " +
	        " c.fecha, c.amount, c.created_at, c.updated_at FROM cuotas c " +
		" JOIN  periods p ON c.Period_id = p.id JOIN  apartas a ON c.aparta_id = a.id " +
		" JOIN  tipos t ON c.tipo_id = t.id  WHERE   p.id = $1 ORDER BY p.inicio, c.fecha" 

	rows, err := Db.Query(stq, id)
	if err != nil {
		log.Println(err)
            return
	}
        defer rows.Close()
        cuot := CuotaN{}
        for rows.Next() {
           if err = rows.Scan(&cuot.Id,&cuot.PeriodId,&cuot.Period, &cuot.ApartaId, &cuot.Apto, &cuot.TipoId, &cuot.Tipo, &cuot.Fecha, &cuot.Amount,  &cuot.CreatedAt, &cuot.UpdatedAt); err != nil {
		log.Println(err)
                 return
            }
           cuotas = append(cuotas, cuot)
         }
       return
 }



// -------------------------------------------------------------
// Get all cuotas in the database and returns the list
  func Cuots() (cuotas []CuotaN, err error) {
        stq :=   "SELECT c.id, c.period_id, p.inicio, c.aparta_id, a.codigo, c.tipo_id, t.descripcion, c.fecha, c.amount, c.created_at, c.updated_at FROM cuotas c, periods p, aparta a, tipos t where c.period_id = p.id and c.aparta_id = a.id and c.Tipo_id = t.id order by p.inicio"
	rows, err := Db.Query(stq)
	if err != nil {
            return
	}
	defer rows.Close()
	for rows.Next() {
            cuot := CuotaN{}
           if err = rows.Scan(&cuot.Id,&cuot.PeriodId,&cuot.Period, &cuot.ApartaId, &cuot.Apto, &cuot.TipoId, &cuot.Tipo,  &cuot.Fecha,&cuot.Amount,  &cuot.CreatedAt, &cuot.UpdatedAt); err != nil {
                  return
             }
             cuotas = append(cuotas, cuot)
	}
        return
 }
// -------------------------------------------------------------
// Get all cuotas per a period in the database and returns the list
  func (cuot * CuotaN)CuotsPer() (cuotas []CuotaN, err error) {
        stq :=   "SELECT c.id, c.period_id, p.inicio, c.aparta_id, a.codigo, c.tipo_id, t.descripcion, c.fecha, c.amount, c.created_at, c.updated_at FROM cuotas c, periods p, apartas a, tipos t where c.period_id = p.id and c.aparta_id = a.id and c.tipo_id = t.id and c.period_id = $1 order by p.inicio"
	rows, err := Db.Query(stq, cuot.PeriodId )
	if err != nil {
            return
	}
	defer rows.Close()
	for rows.Next() {
            c := CuotaN{}
           if err = rows.Scan(&c.Id,&c.PeriodId,&c.Period, &c.ApartaId, &c.Apto, &c.TipoId, &c.Tipo,  &c.Fecha,&c.Amount,  &c.CreatedAt, &c.UpdatedAt); err != nil {
                  return
             }
             cuotas = append(cuotas, c)
	}
        return
 }
// -------------------------------------------------------------
// Get all cuotas per a period in the database and returns the list
  func CuotsPeri(id uint32) (cuotas []CuotaN, err error) {
        stq :=   "SELECT c.id, c.period_id, p.inicio, c.aparta_id, a.codigo, c.tipo_id, t.descripcion, c.fecha, c.amount, c.created_at, c.updated_at FROM cuotas c, periods p, apartas a, tipos t where c.period_id = p.id and c.aparta_id = a.id and c.tipo_id = t.id and c.period_id = $1 order by p.inicio"
	rows, err := Db.Query(stq, id )
	if err != nil {
            return
	}
	defer rows.Close()
	for rows.Next() {
            c := CuotaN{}
           if err = rows.Scan(&c.Id,&c.PeriodId,&c.Period, &c.ApartaId, &c.Apto, &c.TipoId, &c.Tipo,  &c.Fecha,&c.Amount,  &c.CreatedAt, &c.UpdatedAt); err != nil {
                  return
             }
             cuotas = append(cuotas, c)
	}
        return
 }
// -------------------------------------------------------------
// -------------------------------------------------------------
//  merge to merge 2 arrays
    func merge( cuots []CuotApt, fec time.Time) ( []CuotApt){
	    sum := int64(0)
	    for i, _ :=  range cuots{
	        sum              =   sum + cuots[i].Cuota - cuots[i].Amount
		if i > 0 && cuots[i].Id == cuots[i -1].Id{
                       sum  = sum - cuots[i].Cuota
		}
	        cuots[i].Balance = sum
	    }
	    var i  int
	    for i,_ =  range cuots {
		 feci  := cuots[i].Inicio
		 fecj  :=  time.Date(feci.Year(), feci.Month(), feci.Day(), 0, 0, 0, 0, time.Local)
                 if fec.Equal(fecj) || fec.Before(fecj) {
			 break
		 }
	    }
	    return cuots[i:]
    }

// -------------------------------------------------------------
// Payments of a given apt up to a period
   func Payments(aid uint32, fecf time.Time, feci time.Time )(cuotas []CuotApt, err error){
        var rows * sql.Rows
        stq  := "SELECT p.id,  p.inicio, b.cuota, c1.fecha, c1.amount " +
	        " FROM balances b JOIN periods p ON b.period_id = p.id " +
		" left JOIN " +
		" (SELECT p.id AS id, c.fecha AS fecha,c.aparta_id, c.amount AS amount " +
		" FROM cuotas c JOIN periods p ON p.id = c.period_id " +
		" JOIN apartas a ON c.aparta_id = a.id JOIN persons p1 ON p1.aparta_id = a.id " +
		" WHERE c.aparta_id = $1 )  c1 ON p.id = c1.id "  +
		" WHERE p.inicio <= $2  ORDER BY p.inicio"

	rows, err = Db.Query(stq, aid, fecf)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()
	stLayout := "2006-01-02"
        for rows.Next() {
            c := CuotApt{}
	    var  sqFec  sql.NullTime
	    var  sqAmt  sql.NullInt64
            err = rows.Scan(&c.Id, &c.Inicio,  &c.Cuota, &sqFec, &sqAmt )
	    if err != nil {
		log.Println(err)
	         return
	    }
	    if sqFec.Valid{
		 c.Fecha = sqFec.Time
	    }else{
		  c.Fecha, _ = time.Parse(stLayout, "1900-01-01")
	  }
	    if sqAmt.Valid{
		    c.Amount = sqAmt.Int64
	    }else{
		    c.Amount = 0
	    }
            cuotas = append(cuotas, c)
        }
        cuotas = merge(cuotas, feci)
	return
   }
// -------------------------------------------------------------
// Amounts paid in Condo up to a period
   func Amounts( id uint32 )(amts []AmtCond, err error){
    stq  := "SELECT a.codigo, sum(c.amount) AS monto FROM cuotas c JOIN apartas a ON c.aparta_id = a.id JOIN periods p ON c.period_id = p.id WHERE c.fecha <= (SELECT p1.final FROM periods p1 WHERE p1.id =  $1 )GROUP BY a.codigo ORDER BY monto"

        var rows  * sql.Rows
	rows, err = Db.Query(stq,  id)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()
        for rows.Next() {
            c := AmtCond{}
            err = rows.Scan(&c.Codigo, &c.Amount)
	    if err != nil {
		log.Println(err)
	         return
	    }
	    amts = append(amts,c)
        }

	stq = "SELECT sum(b.cuota) AS cuota FROM balances b JOIN periods p ON b.period_id = p.id WHERE p.inicio <= (SELECT p1.final FROM periods p1 WHERE p1.id  = $1) "

	rows, err = Db.Query(stq,  id)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()
        var  a  int64
        if rows.Next() {
            err = rows.Scan(&a)
	    if err != nil {
		log.Println(err)
	         return
	    }
        }
	 for i,_ := range amts {
	      dife               := a - amts[i].Amount
	      mora               :=  dife + dife/ int64(5)
	      amts[i].Atraso      = dife
	      amts[i].Mora        =  mora
	 }

	 return
    }
// -------------------------------------------------------------
//  MoneyFlow per period: cuotas, ingresos, egresos and payments
    func MoneyFlow( id uint32)(cuots []CuotaN,ingresos []IngresoN, egresos []EgresoN, Amts []AmtCond, err error){
	 cuots, err   =  CuotsPeri(id)
	 if err != nil{
             return
	 }
         ingresos, err = IngresLim(id)
	 if err != nil{
             return
	 }
	 egresos,err =  EgresLim(id)
	 if err != nil{
             return
	 }
         Amts, err = Amounts(id)
	 if err != nil{
             return
	 }
         return
    }

// -------------------------------------------------------------
// 
 func  AptDetails(apt string, dIni time.Time, dFini time.Time)(amtTots []AmtAptTot, err error){
	 stq := "select p.inicio,p.final,  coalesce(x.tp, ' '), coalesce(x.ap, $1) as ap, b.cuota,  coalesce(x.fecha, p.inicio) as fecha, coalesce(x.monto,0) as monto  from periods p  join balances b on b.period_id = p.id  left join (   select t.codigo as tp, a.codigo as ap,  c.period_id,  c.fecha,  c.amount as monto  from cuotas c join apartas a on a.id = c.aparta_id  join tipos t on c.tipo_id = t.id where a.codigo = $1 ) x on p.id = x.period_id where p.inicio >= $2 and p.final <= $3 order by  p.inicio, x.fecha "
        var rows * sql.Rows
	rows, err = Db.Query(stq,  apt, dIni, dFini )
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()
        for rows.Next() {
            c := AmtAptTot{}
            err = rows.Scan(&c.Inicio,&c.Final, &c.Tp, &c.Ap,&c.Cuota, &c.Fecha, &c.Monto)
	    if err != nil {
		log.Println(err)
	         return
	    }
	    amtTots = append(amtTots,c)
        }
    return
 }
// -------------------------------------------------------------

