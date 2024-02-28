package db

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/X3NOOO/maud/types"
)

type DB struct {
	conn *sql.DB
}

/*
Initialises the database connection.
*/
func InitDatabase(dsn string) (*DB, error) {
	var db DB
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.conn = conn
	db.conn.SetConnMaxLifetime(time.Minute * 3)
	db.conn.SetMaxOpenConns(10)
	db.conn.SetMaxIdleConns(10)

	return &db, nil
}

func (db *DB) Register(account *types.RegisterPOST) (*types.RegisterResponse, *types.RequestError) {
	return nil, &types.RequestError{
		StatusCode: http.StatusNotImplemented,
		Err:        errors.New("not implemented yet"),
	}
}

func (db *DB) Login(account *types.LoginPOST) (*types.LoginResponse, *types.RequestError) {
	return nil, &types.RequestError{
		StatusCode: http.StatusNotImplemented,
		Err:        errors.New("not implemented yet"),
	}
}

func (db *DB) Alive(account *types.LoginPOST, authorization_token string) *types.RequestError {
	return &types.RequestError{
		StatusCode: http.StatusNotImplemented,
		Err:        errors.New("not implemented yet"),
	}
}
