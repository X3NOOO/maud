package db

import (
	"log"

	"github.com/X3NOOO/maud/types"
)

func (db *DB) GetSwitchesToFire() ([]types.Switch, error) {
	var switches []types.Switch
	rows, err := db.conn.Query("SELECT s.id, a.authorization_token FROM `Switches` s JOIN `Accounts` a ON DATE_ADD(a.alive, INTERVAL s.run_after DAY) <= CURRENT_DATE()")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var sw *types.Switch
		var id int64
		var authorization_token string

		err = rows.Scan(&id, &authorization_token)
		if err != nil {
			log.Println(err)
			continue
		}
		
		sw, rerr := db.getSwitch(id, authorization_token)
		if rerr != nil {
			continue // somethings fucked up idek sql hard
			//return nil, rerr.Err
		}

		switches = append(switches, *sw)
	}
	rows.Close()

	return switches, nil
}