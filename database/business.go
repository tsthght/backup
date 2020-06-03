package database

import (
	"database/sql"
	"errors"
	"fmt"
)

const (
	register_sql = "insert into bk_machine_info (ip) values (?) on duplicate key update update_time = now()"
	status_sql = "replace into bk_machine_info( ip, " +
		"cpu_physic_core_num, cpu_logic_core_num, cpu_percent," +
		"mem_total, mem_used, mem_used_percent, " +
		"disk_path, disk_total, disk_free, disk_used_percent) values (?,   ?, ?, ?,    ?, ?, ?   ,?, ?, ?, ?)"
	getTask_sql = "select uuid from bk_task_info where state = 'todo' order by priority desc, uuid desc limit 1 for update"
	assignTask_sql = "update bk_machine_info set task_id = ? and stage = 'todo' where ip = ?"
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

func AssignFromCmdb(db *sql.DB, ip string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, errors.New("tx Begin failed")
	}
	rows, err := tx.Query(getTask_sql)
	if err != nil {
		return 0, errors.New("tx query failed")
	}
	uuid := -1
	for rows.Next() {
		err := rows.Scan(&uuid)
		if err != nil {
			rows.Close()
			tx.Rollback()
			return 0, errors.New("tx scan failed")
		}
		break
	}
	rows.Close()

	stmt, err := tx.Prepare(assignTask_sql)
	fmt.Printf("assign %s\n", ip)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("tx Prepare failed")
	}
	res, err := stmt.Exec(uuid, ip)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("tx Exec failed")
	}
	tx.Commit()
	return res.RowsAffected()
}