package pg

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PgDB struct {
	db *sql.DB
}

func NewDB(connString string) *PgDB {
    db, err := sql.Open("postgres", connString)
    if err != nil {
        panic(err)
    }

    db.SetMaxOpenConns(100)
    db.SetConnMaxIdleTime(0)
    db.SetMaxIdleConns(0)

    return &PgDB{
        db: db,
    }
}
