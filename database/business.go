package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

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
	gettasktypebyUUID = "select task_type, stage from bk_task_info where uuid = ?"

	getmachinestagebyip = "select stage from bk_machine_info where ip = ?"
	setmachinestagebyip = "update bk_machine_info set stage = ? where ip = ?"

	settaskstatebyuuid = "update bk_task_info set state = ? where uuid = ?"

	gettasksrcdstbyuuid = "select src, dst from bk_task_info where uuid = ?"

	getsqlnode = "select hostname from bladecmdb.blade_sql where physical_cluster_name = ?"

	getrootnode = "select hostname from bladecmdb.blade_root where physical_cluster_name = ?"

	getdbinfobyuuid = "select dbinfo from bk_task_info where uuid = ?"

	settaskstateandmessagebyuuid = "update bk_task_info set state = ?, stage = ?, error_message = ?, pos = ? where uuid = ?"

	setgclifetime = "update mysql.tidb set VARIABLE_VALUE= ? where VARIABLE_NAME='tikv_gc_life_time'"

	getgclifetime = "select VARIABLE_VALUE from mysql.tidb where VARIABLE_NAME = 'tikv_gc_life_time'"

	getmachinenum = "select ip from bk_machine_info where task_id = ? and stage = ?"

	getsrcclustername = "select src from bk_task_info where uuid = ?"

	isbinlogopen = "show variables like 'log_bin'"

	getmaxexecutetime = "show variables like 'max_execution_time'"

	setmaxexecutetime = "set global max_execution_time = ?"

	setatask = "insert into bk_task_info (src, dst, task_type, dbinfo) values (?, ?, ?, ?)"
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
func GetTaskTypeByUUID(db *sql.DB, uuid int) (string, string, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", "", errors.New("call GetTaskTypeByUUID: tx Begin failed: " + err.Error())
	}
	rows, err := tx.Query(gettasktypebyUUID, uuid)
	if err != nil {
		return "", "", errors.New("call GetTaskTypeByUUID: tx Query failed: " + err.Error())
	}
	tp := ""
	stage := ""
	for rows.Next() {
		err := rows.Scan(&tp, &stage)
		if err != nil {
			rows.Close()
			tx.Rollback()
			return tp, stage, errors.New("call GetTaskTypeByUUID: tx scan failed: " + err.Error())
		}
		rows.Close()
		break
	}
	rows.Close()
	tx.Commit()
	return tp, stage, nil
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

func GetCluserBasicInfo(db *sql.DB, uuid int, cfg config.BkConfig, tp int) (*BladeInfo, error) {
	bi := BladeInfo{}
	bi.Port = strconv.Itoa(cfg.Blade.BladePort)
	bi.User = cfg.Blade.BladeUser
	//需要获取appkey
	bi.Password = secret.GetValueByeKey(cfg.Blade.BladeAk, bi.User)
	//获取SQL节点
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.New("call GetTaskTypeByUUID: tx Begin failed: " + err.Error())
	}
	rows, err := tx.Query(gettasksrcdstbyuuid, uuid)
	if err != nil {
		return nil, errors.New("call GetTaskTypeByUUID: tx Query failed: " + err.Error())
	}
	var src, dst, cluster string
	for rows.Next() {
		err := rows.Scan(&src, &dst)
		if err != nil {
			rows.Close()
			tx.Rollback()
			return nil, errors.New("call GetTaskTypeByUUID: tx scan failed: " + err.Error())
		}
		rows.Close()
		break
	}
	rows.Close()

	if tp == DownStream {
		cluster = dst
	} else {
		cluster = src
	}

	rows, err = tx.Query(getsqlnode, cluster)
	if err != nil {
		return nil, errors.New("call GetTaskTypeByUUID: tx Query failed: " + err.Error())
	}

	var sqlnode string
	for rows.Next() {
		err := rows.Scan(&sqlnode)
		if err != nil {
			rows.Close()
			tx.Rollback()
			return nil, errors.New("call GetTaskTypeByUUID: tx scan failed: " + err.Error())
		}
		bi.Hosts = append(bi.Hosts, sqlnode)
	}
	rows.Close()

	rows, err = tx.Query(getrootnode, cluster)
	fmt.Printf("get root, cluster : %s\n", cluster)
	if err != nil {
		return nil, errors.New("call GetTaskTypeByUUID: tx Query failed: " + err.Error())
	}

	var rootnode string
	for rows.Next() {
		err := rows.Scan(&rootnode)
		if err != nil {
			rows.Close()
			tx.Rollback()
			return nil, errors.New("call GetTaskTypeByUUID: tx scan failed: " + err.Error())
		}
		bi.ROOT = append(bi.ROOT, rootnode)
	}
	rows.Close()

	tx.Commit()
	return &bi, nil
}

/*
 * 作用：获取任务的
 */
func GetDBInfoByUUID(db *sql.DB, uuid int) (string, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", errors.New("call GetDBInfoByUUID: tx Begin failed: " + err.Error())
	}
	rows, err := tx.Query(getdbinfobyuuid, uuid)
	if err != nil {
		return "", errors.New("call GetDBInfoByUUID: tx Query failed: " + err.Error())
	}
	dbinfo := ""
	for rows.Next() {
		err := rows.Scan(&dbinfo)
		if err != nil {
			rows.Close()
			tx.Rollback()
			return dbinfo, errors.New("call GetDBInfoByUUID: tx scan failed: " + err.Error())
		}
		rows.Close()
		break
	}
	rows.Close()
	tx.Commit()
	return dbinfo, nil
}

/*
 * 作用：修改任务的状态
 */
func SetTaskStateAndMessageByUUID(db *sql.DB, uuid int, state, stage, message, pos string) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.New("call SetTaskStateByUUID: tx Begin failed: " + err.Error())
	}
	stmt, err := tx.Prepare(settaskstateandmessagebyuuid)
	if err != nil {
		tx.Rollback()
		return errors.New("call SetTaskStateAndMessageByUUID: tx Prepare failed")
	}
	_, err = stmt.Exec(state, stage, message, pos, uuid)
	if err != nil {
		tx.Rollback()
		return errors.New("call SetTaskStateAndMessageByUUID: tx Exec failed")
	}
	tx.Commit()
	return nil
}

/*
 * 作用：设置GC时间
 */
func SetGCTimeByUUID(db *sql.DB, gc string) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.New("call SetGCTimeByUUID: tx Begin failed: " + err.Error())
	}
	stmt, err := tx.Prepare(setgclifetime)
	if err != nil {
		tx.Rollback()
		return errors.New("call SetTaskStateAndMessageByUUID: tx Prepare failed")
	}
	_, err = stmt.Exec(gc)
	if err != nil {
		tx.Rollback()
		return errors.New("call SetTaskStateAndMessageByUUID: tx Exec failed")
	}
	tx.Commit()
	return nil
}

/*
 * 作用：查找GC时间
 */
func GetGCTimeByUUID(db *sql.DB) (error, string) {
	tx, err := db.Begin()
	if err != nil {
		return errors.New("call SetGCTimeByUUID: tx Begin failed: " + err.Error()), ""
	}
	rows, err := tx.Query(getgclifetime)
	if err != nil {
		return errors.New("call GetDBInfoByUUID: tx Query failed: " + err.Error()), ""
	}
	gc := ""
	for rows.Next() {
		err := rows.Scan(&gc)
		if err != nil {
			rows.Close()
			tx.Rollback()
			return errors.New("call GetDBInfoByUUID: tx scan failed: " + err.Error()), gc
		}
		rows.Close()
		break
	}
	rows.Close()
	tx.Commit()
	return nil, gc
}

/*
 * 根据UUID判断执行的数量
 */
func GetMachineNum(db *sql.DB, uuid int, stage string) (error, int){
	tx, err := db.Begin()
	if err != nil {
		return errors.New("call SetGCTimeByUUID: tx Begin failed: " + err.Error()), 0
	}
	rows, err := tx.Query(getmachinenum, uuid, stage)
	if err != nil {
		return errors.New("call GetDBInfoByUUID: tx Query failed: " + err.Error()), 0
	}
	num := 0
	for rows.Next() {
		num ++
	}
	rows.Close()
	tx.Commit()
	return nil, num
}
/*
 * 根据UUID获取pump的数量
 */
func GetMachinePumpIpByPump(db *sql.DB, uuid int, stage string) (error, []string){
	tx, err := db.Begin()
	if err != nil {
		return errors.New("call SetGCTimeByUUID: tx Begin failed: " + err.Error()), nil
	}
	rows, err := tx.Query(getmachinenum, uuid, stage)
	if err != nil {
		return errors.New("call GetDBInfoByUUID: tx Query failed: " + err.Error()), nil
	}
	var hosts []string
	host := ""
	for rows.Next() {
		err := rows.Scan(&host)
		if err != nil {
			rows.Close()
			tx.Rollback()
			return errors.New("call GetSrcClusterName: tx scan failed: " + err.Error()), hosts
		}
		hosts = append(hosts, host)
	}
	rows.Close()
	tx.Commit()
	return nil, hosts
}
/*
 * 根据UUID获取集群名
 */
func GetSrcClusterName(db *sql.DB, uuid int) (error, string) {
	tx, err := db.Begin()
	if err != nil {
		return errors.New("call GetSrcClusterName: tx Begin failed: " + err.Error()), ""
	}
	rows, err := tx.Query(getsrcclustername, uuid)
	if err != nil {
		return errors.New("call GetSrcClusterName: tx Query failed: " + err.Error()), ""
	}
	src := ""
	for rows.Next() {
		err := rows.Scan(&src)
		if err != nil {
			rows.Close()
			tx.Rollback()
			return errors.New("call GetSrcClusterName: tx scan failed: " + err.Error()), src
		}
		rows.Close()
		break
	}
	rows.Close()
	tx.Commit()
	return nil, src
}
/*
 * 获取是否打开binlog
 */
func IsBinlogOpen(db *sql.DB) (error, int) {
	tx, err := db.Begin()
	if err != nil {
		return errors.New("call IsBinlogOpen: tx Begin failed: " + err.Error()), 0
	}
	rows, err := tx.Query(isbinlogopen)
	if err != nil {
		return errors.New("call IsBinlogOpen: tx Query failed: " + err.Error()), 0
	}
	key := ""
	value := 0
	for rows.Next() {
		err := rows.Scan(&key, &value)
		if err != nil {
			rows.Close()
			tx.Rollback()
			return errors.New("call IsBinlogOpen: tx scan failed: " + err.Error()), value
		}
		rows.Close()
		break
	}
	rows.Close()
	tx.Commit()
	return nil, value
}

/*
 * 获取max execute time
 */
func GetMaxExecuteTime(db *sql.DB) (error, int){
	tx, err := db.Begin()
	if err != nil {
		return errors.New("call GetMaxExecuteTime: tx Begin failed: " + err.Error()), 0
	}
	rows, err := tx.Query(getmaxexecutetime)
	if err != nil {
		return errors.New("call GetMaxExecuteTime: tx Query failed: " + err.Error()), 0
	}
	key := ""
	value := 0
	for rows.Next() {
		err := rows.Scan(&key, &value)
		if err != nil {
			rows.Close()
			tx.Rollback()
			return errors.New("call GetMaxExecuteTime: tx scan failed: " + err.Error()), value
		}
		rows.Close()
		break
	}
	rows.Close()
	tx.Commit()
	return nil, value
}
/*
 * 设置max execute time
 */
func SetMaxExecuteTime(db *sql.DB, exetime int) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.New("call SetMaxExecuteTime: tx Begin failed: " + err.Error())
	}
	stmt, err := tx.Prepare(setmaxexecutetime)
	if err != nil {
		tx.Rollback()
		return errors.New("call SetMaxExecuteTime: tx Prepare failed")
	}
	_, err = stmt.Exec(exetime)
	if err != nil {
		tx.Rollback()
		return errors.New("call SetMaxExecuteTime: tx Exec failed")
	}
	tx.Commit()
	return nil
}

/*
 * 设置任务
 */
func SetATask(db *sql.DB, src, dst, tp, dt string) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.New("call SetATask: tx Begin failed: " + err.Error())
	}
	stmt, err := tx.Prepare(setatask)
	if err != nil {
		tx.Rollback()
		return errors.New("call SetATask: tx Prepare failed")
	}
	_, err = stmt.Exec(src, dst, tp, dt)
	if err != nil {
		tx.Rollback()
		return errors.New("call SetATask: tx Exec failed")
	}
	tx.Commit()
	return nil
}