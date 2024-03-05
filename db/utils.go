package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/X3NOOO/maud/types"
)

func (db *DB) getSwitch(switch_id int64) (*types.Switch, *types.RequestError) {
	var response types.Switch

	var recipients_json string
	err := db.conn.QueryRow("SELECT content, run_after, recipients FROM `Switches` WHERE id = ?", switch_id).Scan(&response.Content, &response.Run_after, &recipients_json)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &types.RequestError{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("switch not found"),
			}
		}
		return nil, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to get switch data"),
		}
	}

	err = json.Unmarshal([]byte(recipients_json), &response.Recipients)
	if err != nil {
		return nil, &types.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to unmarshal recipients"),
		}
	}

	response.Id = switch_id

	return &response, nil
}

