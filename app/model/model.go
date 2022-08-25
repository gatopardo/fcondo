package model

import (
	"database/sql"
	"errors"
 //	"encoding/json"
	"fmt"
	"log"
        "os"
        "regexp"
//	"time"

      _ "github.com/lib/pq"   // PostgreSQL driver
      _ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/mgo.v2"

)

var (
	// ErrCode is a config or an internal error
	ErrCode = errors.New("Sentencia Case en codigo no es correcta.")
	// ErrNoResult is a not results error
	ErrNoResult = errors.New("Resultado  no encontrado.")
	// ErrUnavailable is a database not available error
	ErrUnavailable = errors.New("Database no disponible.")
	// ErrUnauthorized is a permissions violation
	ErrUnauthorized = errors.New("Usuario sin permiso para realizar esta operacion.")
)

var (
	// SQL wrapper
	SQL *sqlx.DB
        // Postgresql wrapper
         Db *sql.DB

	// Database info
	databases Info
	formato = "2006-01-02"

)

// Type is the type of database from a Type* constant
type Type string

const (
	// TypeMySQL is MySQL
	TypeMySQL Type = "MySQL"
	TypePostgreSQL Type = "PostgreSQL"
)

// Info contains the database configurations
type Info struct {
         Remote  bool
	// Database type
	Type Type
	// MySQL info if used
	MySQL MySQLInfo
	// PostgreSQL info if used
	PostgreSQL  PostgreSQLInfo
}

// MySQLInfo is the details for the database connection
type MySQLInfo struct {
	Username  string
	Password  string
	Name      string
	Hostname  string
	Port      int
	Parameter string
}

type PostgreSQLInfo struct {
	Username  string
	Password  string
	Name      string
	Hostname  string
	Port      int
	Parameter string
}

// standardizeErrors returns the same error regardless of the database used
func standardizeError(err error) error {
	if err == sql.ErrNoRows || err == mgo.ErrNotFound {
		return ErrNoResult
	}

	return err
}

// DSN returns the Data Source Name
func MyDSN(ci MySQLInfo) string {
	// Example: root:@tcp(localhost:3306)/test
	return ci.Username + ":" + ci.Password + "@tcp(" +
		ci.Hostname + ":" + fmt.Sprintf("%d", ci.Port) + ")/" +
		ci.Name + ci.Parameter
}

  func PgDNS(ci PostgreSQLInfo  ) string {
         return   fmt.Sprintf("user=%s dbname=%s port=%d sslmode=%s",ci.Username, ci.Name, ci.Port, ci.Parameter)
     }
// Connect to the database
func Connect(d Info) {
	var err error
	// Store the config
	databases = d
        if d.Remote  {
          regex := regexp.MustCompile("(?i)^postgres://(?:([^:@]+):([^@]*)@)?([^@/:]+):(\\d+)/(.*)$")
          matches := regex.FindStringSubmatch(os.Getenv("DATABASE_URL"))
	  if matches == nil {
		log.Fatalf("DATABASE_URL variable must look like: postgres://username:password@hostname:port/dbname (not '%v')", os.Getenv("DATABASE_URL"))
	   }

	sslmode := os.Getenv("PGSSL")
	if sslmode == "" {
		sslmode = "disable"
	}
//	spec := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", matches[1], matches[2], matches[3], matches[4], matches[5], sslmode)

//          Db, err = sql.Open("postgres", spec)
            Db, err =  sql.Open("postgres", os.Getenv("DATABASE_URL"))
          if err != nil {
                 log.Println(err)
              }


        } else {

          if Db,err  = sql.Open("postgres", PgDNS(d.PostgreSQL)); err !=  nil{
			log.Println("SQL Driver Error", err)
                        log.Fatal("Connection to database error")
           }
		if err = Db.Ping(); err != nil {
			log.Println("Database Error", err)
		}
	}
}


// ReadConfig returns the database information
func ReadConfig() Info {
	return databases
}

