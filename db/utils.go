package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/X3NOOO/maud/types"
)

func (db *DB) getSwitch(switch_id int64, authorization_token string) (*types.Switch, *types.RequestError) {
	var response types.Switch
	var recipients_json string
	err := db.conn.QueryRow("SELECT account_id, content, run_after, recipients FROM `Switches` WHERE id = ? AND account_id = (SELECT a.id FROM `Accounts` a WHERE a.authorization_token = ?)", switch_id, authorization_token).Scan(&response.AccountId, &response.Content, &response.Run_after, &recipients_json)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &types.RequestError{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("switch not found"),
			}
		}
		log.Println(err)
		return nil, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to get switch data"),
		}
	}

	if recipients_json != "" {
		err = json.Unmarshal([]byte(recipients_json), &response.Recipients)
		if err != nil {
			log.Println(err)
			return nil, &types.RequestError{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("failed to unmarshal recipients"),
			}
		}
	}

	response.Id = switch_id

	return &response, nil
}
