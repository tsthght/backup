package database

import (
	"database/sql"
	"errors"
	"fmt"
)

const (
	register_sql = "insert into bk_machine_info (ip) values (?) on duplicate key update update_time = now()"
	status_sql = "update bk_machine_info set cpu_physic_core_num = ?, cpu_logic_core_num = ?, cpu_percent = ?," +
		"mem_total = ?, mem_used = ?, mem_used_percent = ?, " +
		"disk_path = ?, disk_total = ?, disk_free = ?, disk_used_percent = ? " +
		"where ip = ?"
	getTask_sql = "select uuid from bk_task_info where state = 'todo' order by priority desc, uuid desc limit 1 for update"
	assignTask_sql = "update bk_machine_info set task_id = ?, stage = 'todo' where ip = ?"
	getStatus_sql = "select stage from bk_machine_info where ip = ?"
	doingTask_ssql = "update bk_task_info set state = 'doing' where uuid = ?"
)

func RegisterToCmdb(db *sql.DB, ip string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, errors.New("call RegisterToCmdb: tx Begin failed: " + err.Error())
	}
	stmt, err := tx.Prepare(register_sql)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("call RegisterToCmdb: tx Prepare failed: " + err.Error())
	}
	res, err := stmt.Exec(ip)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("call RegisterToCmdb: tx Exec failed: " + err.Error())
	}
	tx.Commit()
	return res.RowsAffected()
}

func StatusUpdateToCmdb(db *sql.DB, ip string, info CPUInfo, mem MEMInfo, dsk DiskInfo, path string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, errors.New("call StatusUpdateToCmdb: tx Begin failed: " + err.Error())
	}
	stmt, err := tx.Prepare(status_sql)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("call StatusUpdateToCmdb: tx Prepare failed: " + err.Error())
	}
	res, err := stmt.Exec(info.PhysicCoreNum, info.LogicCoreNum, info.Percent,
		mem.TotalSize, mem.Available, mem.UsedPercent,
		path, dsk.TotalSize, dsk.Free, dsk.UsedPercent, ip)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("call StatusUpdateToCmdb: tx Exec failed: " + err.Error())
	}
	tx.Commit()
	return res.RowsAffected()
}

func AssignFromCmdb(db *sql.DB, ip string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, errors.New("call AssignFromCmdb: tx Begin failed")
	}
	rows, err := tx.Query(getTask_sql)
	if err != nil {
		return 0, errors.New("call AssignFromCmdb: tx query failed")
	}
	var uuid int64  = -1
	for rows.Next() {
		err := rows.Scan(&uuid)
		if err != nil {
			rows.Close()
			tx.Rollback()
			return 0, errors.New("call AssignFromCmdb: tx scan failed")
		}
		if uuid == -1 {
			rows.Close()
			tx.Commit()
			return uuid, nil
		}
		break
	}
	rows.Close()
	if uuid == -1 {
		tx.Commit()
		return uuid, nil
	}

	stmt, err := tx.Prepare(assignTask_sql)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("call AssignFromCmdb: tx Prepare failed")
	}
	res, err := stmt.Exec(uuid, ip)
	fmt.Printf("assign to:  %d, %s\n", uuid, ip)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("call AssignFromCmdb: tx Exec failed")
	}

	stmt, err = tx.Prepare(doingTask_ssql)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("call AssignFromCmdb: tx Prepare failed")
	}

	res, err = stmt.Exec(uuid)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("call AssignFromCmdb: tx Exec failed")
	}

	tx.Commit()
	return res.RowsAffected()
}

func GetStatusFromCmdb(db *sql.DB, ip string) (string, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", errors.New("call GetStatusFromCmdb: tx Begin failed: " + err.Error())
	}
	if err != nil {
		tx.Rollback()
		return "", errors.New("call GetStatusFromCmdb: tx Prepare failed:" + err.Error())
	}
	rows, err := tx.Query(getStatus_sql, ip)
	if err != nil {
		return "", errors.New("call GetStatusFromCmdb: tx query failed:" + err.Error())
	}
	stage := ""
	for rows.Next() {
		err := rows.Scan(&stage)
		if err != nil {
			rows.Close()
			tx.Rollback()
			return "", errors.New("call GetStatusFromCmdb: tx scan failed:" + err.Error())
		}
		rows.Close()
		break
	}
	tx.Commit()
	return stage, nil
}