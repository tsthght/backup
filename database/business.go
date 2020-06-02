package database

import (
	"database/sql"
	"errors"
)

const (
	register_sql = "insert into bk_machine_info (ip) values (?) on duplicate key update update_time = now()"
	status_sql = "replace into bk_machine_info(ip, cpu_physic_core_num, cpu_logic_core_num, cpu_percent) values (?, ?, ?, ?)"
)

func RegisterToCmdb(db *sql.DB, ip string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, errors.New("tx Begin failed")
	}
	stmt, err := tx.Prepare(register_sql)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("tx Prepare failed")
	}
	res, err := stmt.Exec(ip)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("tx Exec failed")
	}
	tx.Commit()
	return res.RowsAffected()
}

func StatusUpdateToCmdb(db *sql.DB, ip string, info CPUInfo) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, errors.New("tx Begin failed")
	}
	stmt, err := tx.Prepare(status_sql)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("tx Prepare failed")
	}
	res, err := stmt.Exec(ip, info.PhysicCoreNum, info.LogicCoreNum, info.Percent)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("tx Exec failed")
	}
	tx.Commit()
	return res.RowsAffected()
}