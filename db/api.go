package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	_ "embed"

	"github.com/alexedwards/argon2id"
	_ "github.com/go-sql-driver/mysql"

	"github.com/X3NOOO/maud/crypto"
	"github.com/X3NOOO/maud/types"
)

//go:embed maud.sql
var SQL_INIT string

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
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)

	// workaround over not using multiple statements in one exec
	for _, s := range strings.Split(SQL_INIT, ";") {
		if strings.TrimSpace(s) == "" {
			continue
		}
		_, err = conn.Exec(s + ";")
		if err != nil {
			return nil, err
		}
	}

	db.conn = conn
	return &db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Authorize(authorization_token string) (bool, *types.RequestError) {
	var result bool
	err := db.conn.QueryRow("SELECT COUNT(*) FROM `Accounts` WHERE authorization_token = ?", authorization_token).Scan(&result)
	if err != nil {
		return false, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to verify login"),
		}
	}

	return result, nil
}

func (db *DB) Register(account types.RegisterPOST) (*types.RegisterResponse, *types.RequestError) {
	var response types.RegisterResponse

	password_hash, err := argon2id.CreateHash(account.Password, argon2id.DefaultParams)
	if err != nil {
		return &response, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to generate the hash of the password"),
		}
	}

	authorization_token, err := crypto.RandomString(128)
	if err != nil {
		return &response, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to generate authorization token"),
		}
	}

	insert, err := db.conn.Exec("INSERT INTO `Accounts` (nick, `password`, authorization_token) VALUES (?, ?, ?)", account.Nick, password_hash, authorization_token)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062 ") {
			return &response, &types.RequestError{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("this account already exists"),
			}
		}
		return &response, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to create account"),
		}
	}

	id, err := insert.LastInsertId()
	if err != nil {
		return &response, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to create account"),
		}
	}

	err = db.conn.QueryRow("SELECT nick, authorization_token FROM `Accounts` WHERE id = ?", id).Scan(&response.Nick, &response.AuthorizationToken)
	if err != nil {
		return &response, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to verify account creation"),
		}
	}

	return &response, nil
}

func (db *DB) Login(account types.LoginPOST) (*types.LoginResponse, *types.RequestError) {
	var response types.LoginResponse
	var password_hash_db string

	err := db.conn.QueryRow("SELECT password, authorization_token FROM `Accounts` WHERE nick = ?", account.Nick).Scan(&password_hash_db, &response.AuthorizationToken)
	if err != nil {
		return &response, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to get account details"),
		}
	}

	match, err := argon2id.ComparePasswordAndHash(account.Password, password_hash_db)
	if err != nil {
		log.Println(err.Error())
		return &response, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to compare hashes"),
		}
	}

	if !match {
		return nil, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("passwords don't match"),
		}
	}

	return &response, nil
}

func (db *DB) Status(authorization_token string) (*types.StatusResponse, *types.RequestError) {
	var result types.StatusResponse
	err := db.conn.QueryRow("SELECT alive FROM `Accounts` WHERE authorization_token = ?", authorization_token).Scan(&result.Date)
	if err != nil {
		return nil, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to get status from the database"),
		}
	}

	return &result, nil
}

func (db *DB) UpdateAlive(authorization_token string) (*types.AliveResponse, *types.RequestError) {
	var result types.AliveResponse

	_, err := db.conn.Exec("UPDATE `Accounts` SET alive = CURRENT_DATE() WHERE authorization_token = ?", authorization_token)
	if err != nil {
		return nil, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to update alive status"),
		}
	}

	err = db.conn.QueryRow("SELECT alive FROM `Accounts` WHERE authorization_token = ?", authorization_token).Scan(&result.Date)
	if err != nil {
		return nil, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to confirm new alive status"),
		}
	}

	return &result, nil
}

func (db *DB) AddSwitch(authorization_token string, switch_body types.SwitchesPOST) (*types.Switch, *types.RequestError) {
	var response *types.Switch
	var account_id uint

	err := db.conn.QueryRow("SELECT id FROM `Accounts` WHERE authorization_token = ?", authorization_token).Scan(&account_id)
	if err != nil {
		return nil, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to get the account id"),
		}
	}

	recipients_json, err := json.Marshal(switch_body.Recipients)
	if err != nil {
		return nil, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to parse recipients"),
		}
	}

	var switch_id int64

	err = db.conn.QueryRow("SELECT MAX(id) FROM `Switches` WHERE account_id = ?", account_id).Scan(&switch_id)
	if err != nil {
		if strings.HasSuffix(err.Error(), "converting NULL to int64 is unsupported") {
			switch_id = 0
		} else {
			return nil, &types.RequestError{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("failed to get known rows"),
			}
		}
	}
	switch_id++

	_, err = db.conn.Exec("INSERT INTO `Switches` (account_id, id, content, subject, run_after, recipients) VALUES (?, ?, ?, ?, ?, ?)", account_id, switch_id, switch_body.Content, switch_body.Subject, switch_body.Run_after, string(recipients_json))
	if err != nil {
		return nil, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to create switch"),
		}
	}

	response, rerr := db.getSwitch(switch_id, authorization_token)
	if rerr != nil {
		return nil, rerr
	}

	return response, nil
}

/*
If switch_id == -1 will return an array of all the switches
*/
func (db *DB) GetSwitch(authorization_token string, switch_id int64) (*types.SwitchesResponse, *types.RequestError) {
	var response types.SwitchesResponse
	var switches []int64

	if switch_id == -1 {
		rows, err := db.conn.Query("SELECT s.id FROM `Switches` s WHERE s.account_id = (SELECT a.id FROM `Accounts` a WHERE a.authorization_token = ?)", authorization_token)
		if err != nil {
			return nil, &types.RequestError{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("failed to get switches"),
			}
		}

		var sw int64
		for rows.Next() {
			rows.Scan(&sw)
			switches = append(switches, sw)
		}
		err = rows.Close()
		if err != nil {
			return nil, &types.RequestError{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("failed to close rows"),
			}
		}
	} else {
		switches = []int64{switch_id}
	}

	for _, item := range switches {
		sw, rerr := db.getSwitch(item, authorization_token)
		if rerr != nil {
			return nil, rerr
		}
		response.Switches = append(response.Switches, *sw)
	}

	return &response, nil
}

func (db *DB) DeleteSwitch(authorization_token string, switch_id int64) (*types.Switch, *types.RequestError) {
	var response *types.Switch
	response, rerr := db.getSwitch(switch_id, authorization_token)
	if rerr != nil {
		return nil, rerr
	}

	_, err := db.conn.Exec("DELETE FROM `Switches` WHERE id = ? AND account_id = (SELECT id from `Accounts` WHERE authorization_token = ?)", switch_id, authorization_token)
	if err != nil {
		return nil, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to remove the requested switch"),
		}
	}

	_, rerr = db.getSwitch(switch_id, authorization_token)
	if rerr.Error() != "switch not found" {
		return nil, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to verify the deletion of the switch"),
		}
	}

	return response, nil
}

func (db *DB) UpdateSwitch(authorization_token string, switch_id int64, switch_body types.SwitchesPATCH) (*types.Switch, *types.RequestError) {
	var response *types.Switch

	old_switch, rerr := db.getSwitch(switch_id, authorization_token)
	if rerr != nil {
		return nil, rerr
	}

	new_switch := old_switch

	j, err := json.Marshal(switch_body)
	if err != nil {
		return nil, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to create new switch structure"),
		}
	}
	err = json.Unmarshal(j, new_switch)
	if err != nil {
		return nil, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to create new switch structure"),
		}
	}

	new_recipients_json, err := json.Marshal(new_switch.Recipients)
	if err != nil {
		return nil, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to marshal recipients"),
		}
	}

	_, err = db.conn.Exec("UPDATE `Switches` SET content = ?, subject = ?, run_after = ?, recipients = ? WHERE id = ? AND account_id = (SELECT id from `Accounts` WHERE authorization_token = ?)", new_switch.Content, new_switch.Subject, new_switch.Run_after, new_recipients_json, switch_id, authorization_token)
	if err != nil {
		return nil, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to update switch"),
		}
	}

	response, rerr = db.getSwitch(switch_id, authorization_token)
	return response, rerr
}
