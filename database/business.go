package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tsthght/backup/config"
	"github.com/tsthght/backup/secret"
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

	gettaskUUIDasignedtomachine = "select task_id from bk_machine_info where ip = ?"
	gettasktypebyUUID = "select task_type from bk_task_info where uuid = ?"

	getmachinestagebyip = "select stage from bk_machine_info where ip = ?"
	setmachinestagebyip = "update bk_machine_info set stage = ? where ip = ?"

	settaskstatebyuuid = "update bk_task_info set state = ? where uuid = ?"
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

/*
 * 作用：获取当前机器（ip）被分配的任务ID
 * 返回值：-1 没有被分配任务
 *        [0...] 被分配的任务号
 */
func GetTaskUUIDAsignedToMachine(db *sql.DB, ip string) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return -1, errors.New("call GetTaskUUIDAsignedToMachine: tx Begin failed: " + err.Error())
	}
	rows, err := tx.Query(gettaskUUIDasignedtomachine, ip)
	if err != nil {
		return -1, errors.New("call GetTaskUUIDAsignedToMachine: tx Query failed: " + err.Error())
	}
	uuid := -1
	for rows.Next() {
		err := rows.Scan(&uuid)
		if err != nil {
			rows.Close()
			tx.Rollback()
			return uuid, errors.New("call GetTaskUUIDAsignedToMachine: tx scan failed: " + err.Error())
		}
		rows.Close()
		break
	}
	tx.Commit()
	return uuid, nil
}

/*
 * 作用：获取当前机器（ip）的阶段 stage
 * 返回值：idle 表示还没有开始任务，需要开始
 *        xxx  表示各个阶段
 */
func GetMachineStageByIp(db *sql.DB, ip string) (string, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", errors.New("call GetMachineStageById: tx Begin failed: " + err.Error())
	}
	rows, err := tx.Query(getmachinestagebyip, ip)
	if err != nil {
		return "", errors.New("call GetMachineStageById: tx Query failed: " + err.Error())
	}
	stage := ""
	for rows.Next() {
		err := rows.Scan(&stage)
		if err != nil {
			rows.Close()
			tx.Rollback()
			return stage, errors.New("call GetMachineStageById: tx scan failed: " + err.Error())
		}
		rows.Close()
		break
	}
	tx.Commit()
	return stage, nil
}

/*
 * 作用：设置当前机器（ip）的阶段 stage
 * 返回值：error
 */
func SetMachineStageByIp(db *sql.DB, ip ,state string) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.New("call SetMachineStageByIp: tx Begin failed: " + err.Error())
	}
	stmt, err := tx.Prepare(setmachinestagebyip)
	if err != nil {
		tx.Rollback()
		return errors.New("call SetMachineStageByIp: tx Prepare failed")
	}
	_, err = stmt.Exec(state, ip)
	if err != nil {
		tx.Rollback()
		return errors.New("call AssignFromCmdb: tx Exec failed")
	}
	tx.Commit()
	return nil
}

/*
 * 作用：获取当前任务的类型type
 * 返回值："" or "schema" or "full" or "all"
 */
func GetTaskTypeByUUID(db *sql.DB, uuid int) (string, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", errors.New("call GetTaskTypeByUUID: tx Begin failed: " + err.Error())
	}
	rows, err := tx.Query(gettasktypebyUUID, uuid)
	if err != nil {
		return "", errors.New("call GetTaskTypeByUUID: tx Query failed: " + err.Error())
	}
	tp := ""
	for rows.Next() {
		err := rows.Scan(&tp)
		if err != nil {
			rows.Close()
			tx.Rollback()
			return tp, errors.New("call GetTaskTypeByUUID: tx scan failed: " + err.Error())
		}
		rows.Close()
		break
	}
	rows.Close()
	tx.Commit()
	return tp, nil
}

/*
 * 作用：设置任务的状态 state
 * 返回值：error
 */
func SetTaskStageByUUID(db *sql.DB, uuid int ,state string) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.New("call SetMachineStageByIp: tx Begin failed: " + err.Error())
	}
	stmt, err := tx.Prepare(settaskstatebyuuid)
	if err != nil {
		tx.Rollback()
		return errors.New("call SetMachineStageByIp: tx Prepare failed")
	}
	_, err = stmt.Exec(state, uuid)
	if err != nil {
		tx.Rollback()
		return errors.New("call AssignFromCmdb: tx Exec failed")
	}
	tx.Commit()
	return nil
}

/*
 * 作用：获得集群的基本信息
 */

func GetCluserBasicInfo(db *sql.DB, uuid int, cfg config.BkConfig, tp int) (BladeInfo, error) {
	bi := BladeInfo{}
	bi.User = cfg.Blade.BladeUser
	//需要获取appkey
	bi.Password = secret.GetValueByeKey(cfg.Blade.BladeAk, bi.User)
	fmt.Printf("bladeinfo: %v\n", bi)
	//获取tidb节点

	return bi, nil
}