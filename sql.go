package main

import (
    "os"
    "fmt"
    "log"
    "regexp"
    "time"

    "database/sql"

    _ "github.com/lib/pq"
    // "github.com/coopernurse/gorp"
    "github.com/go-gorp/gorp"
)

/*
    Typical .bash_profile

## default vars                                                                 
export NODE_ENV=development
export ENV=development
export POSTGRESQL_LOCAL_URL="postgres://misrab: @localhost:5432/misrab"
*/



/*
    Models
*/

type Amounts struct {
    Id int64 `db:"id"`
    Created int64 
    Updated int64

    Json string
}
func (i *Amounts) PreInsert(s gorp.SqlExecutor) error {
    i.Created = time.Now().UnixNano()
    i.Updated = i.Created
    return nil
}
func (i *Amounts) PreUpdate(s gorp.SqlExecutor) error {
    i.Updated = time.Now().UnixNano()
    return nil
}


/*
    Setup
*/


func SetupDB() *gorp.DbMap {
    db := pgConnect()

    // construct a gorp DbMap
    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

    // add a table, setting the table name to 'posts' and
    // specifying that the Id property is an auto incrementing PK
    dbmap.AddTableWithName(Amounts{}, "amounts").SetKeys(true, "Id")


    // drop all tables for testing
    // log.Println("DROPPING TABLES!")
    // err1 := dbmap.DropTablesIfExists()
    // if err1 != nil { panic(err1) }

    err2 := dbmap.CreateTablesIfNotExists()
    if err2 != nil {
        panic(err2)
    }

    // set logging for development
    env := os.Getenv("ENV")
    if env == "" || env == "development" {
        dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds)) 
    }
    
    return dbmap
}


func pgConnect() *sql.DB {
    // Connect to Postgres database
    env := os.Getenv("ENV")
    regex := regexp.MustCompile("(?i)^postgres://(?:([^:@]+):([^@]*)@)?([^@/:]+):(\\d+)/(.*)$")
    var connection string
    switch env {
    //case "staging":
    //case "production":
    // default to development
    default:
        connection = os.Getenv("POSTGRESQL_LOCAL_URL")
    }
    matches := regex.FindStringSubmatch(connection)
    sslmode := "disable"
    spec := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", matches[1], matches[2], matches[3], matches[4], matches[5], sslmode)

    db, err := sql.Open("postgres", spec)
    //PanicIf(err)
    if err != nil {
        panic(err)
    }

    log.Printf("Connected to %s\n", connection)

    return db
}