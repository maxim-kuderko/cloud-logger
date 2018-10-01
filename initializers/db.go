package initializers

import (
	"database/sql"
	"log"
)

type Db struct {
	conn *sql.DB
}

func NewDb(DSN string) *Db {
	db, err := sql.Open("mysql", DSN)
	db.SetMaxIdleConns(1000)
	db.SetMaxOpenConns(1000)
	if err != nil {
		log.Fatal(err)
	}
	return &Db{
		conn: db,
	}
}

func (db *Db) Prepare(q string) (*sql.Stmt,error){
	return db.conn.Prepare(q)
}