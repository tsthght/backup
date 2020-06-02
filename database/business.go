package database

import (
	"database/sql"
	"errors"
)

const (
	register_sql = "replace into bk_machine_info (ip) values (?)"
)

func RegisterToCmdb(db *sql.DB, ip string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, errors.New("tx Begin failed")
	}
	stmt, err := tx.Prepare(register_sql)
	if err != nil {
		return 0, errors.New("tx Prepare failed")
	}
	res, err := stmt.Exec(ip)
	if err != nil {
		return 0, errors.New("tx Exec failed")
	}
	tx.Commit()
	return res.RowsAffected()
}
