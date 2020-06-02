package database

import (
	"database/sql"
	"errors"
)

const (
	register_sql = "insert into bk_machine_info (ip) values (?) on duplicate key update update_time = now()"
	status_sql = "insert into bk_machine_info( ip, " +
		"cpu_physic_core_num, cpu_logic_core_num, cpu_percent," +
		"mem_total, mem_used, mem_used_percent, " +
		"disk_path, disk_total, disk_free, disk_used_percent) values (?,   ?, ?, ?,    ?, ?, ?   ?, ?, ?, ?)"
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

func StatusUpdateToCmdb(db *sql.DB, ip string, info CPUInfo, mem MEMInfo, dsk DiskInfo, path string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, errors.New("tx Begin failed")
	}
	stmt, err := tx.Prepare(status_sql)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("tx Prepare failed")
	}
	res, err := stmt.Exec(ip, info.PhysicCoreNum, info.LogicCoreNum, info.Percent,
		mem.TotalSize, mem.Available, mem.UsedPercent,
		path, dsk.TotalSize, dsk.Free, dsk.UsedPercent)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("tx Exec failed")
	}
	tx.Commit()
	return res.RowsAffected()
}