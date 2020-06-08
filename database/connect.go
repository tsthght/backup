package database

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	primary_uuid = "SHOW STATUS LIKE 'group_replication_primary_member'"
	current_uuid = "SHOW VARIABLES LIKE 'server_uuid'"
)

func GetMGRConnection(cluster *MGRInfo, userinfo UserInfo, writenode bool) *sql.DB {
	ips := cluster.Hosts
	l := 0
	if l = len(ips); l == 0 {
		return nil
	}

	//get connect by random
	if !writenode {
		index := rand.Intn(l - 1)
		for i := 0; i < l; i++ {
			ip := ips[index%l]
			ref := strings.Join([]string{userinfo.Username, ":", userinfo.Password, "@tcp(",ip, ":", userinfo.Port, ")/", userinfo.Database, "?charset=utf8"}, "")
			db, _ := sql.Open("mysql", ref)
			if db == nil {
				fmt.Printf("open db connection failed, ref: %s, index: %d\n", ref, index)
				index ++
				continue
			}
			if err := db.Ping(); err != nil {
				index ++
				db.Close()
				continue
			} else {
				fmt.Printf("ref: %v\n", ref)
				return db
			}
		}
		return nil
	}

	//must return primary node
	index := cluster.WriteIndex

	var err error
	var pu, cu string
	for i := 0; i< l; i++ {
		ip := ips[index%l]
		ref := strings.Join([]string{userinfo.Username, ":", userinfo.Password, "@tcp(",ip, ":", userinfo.Port, ")/", userinfo.Database, "?charset=utf8"}, "")
		db, _ := sql.Open("mysql", ref)
		if len(pu) == 0 {
			err, pu = getPrimaryUUID(db)
			if err != nil {
				db.Close()
				index ++
				continue
			}
		}
		err, cu = getCurrentUUID(db)
		if strings.EqualFold(pu, cu) {
			cluster.WriteIndex = index
			return db
		}
		db.Close()
		index ++
	}
	return nil
}

func getPrimaryUUID(db *sql.DB) (error, string) {
	rows, err := db.Query(primary_uuid)
	if err != nil {
		return nil, ""
	}
	for rows.Next() {
		Variable_name := ""
		Value := ""
		err := rows.Scan(&Variable_name, &Value)
		if err != nil {
			rows.Close()
			return err, ""
		}
		rows.Close()
		return nil, Value
	}
	return errors.New("unexpected error when call GetPrimaryUUID"), ""
}

func getCurrentUUID(db *sql.DB) (error, string) {
	rows, err := db.Query(current_uuid)
	if err != nil {
		return nil, ""
	}
	for rows.Next() {
		Variable_name := ""
		Value := ""
		err := rows.Scan(&Variable_name, &Value)
		if err != nil {
			rows.Close()
			return err, ""
		}
		rows.Close()
		return nil, Value
	}
	return errors.New("unexpected error when call GetPrimaryUUID"), ""
}

func GetTiDBConnection(cluster *BladeInfo) *sql.DB {
	hosts := cluster.Hosts
	l := 0
	if l = len(hosts); l == 0 {
		return nil
	}
	index := rand.Intn(l - 1)
	for i := 0; i < l; i++ {
		host := hosts[index%l]
		addr, err := net.LookupAddr(host)
		if err != nil {
			fmt.Printf("call LookupAddr failed, err : %s", err.Error())
			return nil
		}
		if len(addr) == 0 {
			fmt.Printf("addr is nil")
			return nil
		}
		ref := strings.Join([]string{cluster.User, ":", cluster.Password, "@tcp(",addr[0], ":", cluster.Port, ")/", cluster.Database, "?charset=utf8"}, "")
		db, _ := sql.Open("mysql", ref)
		if err := db.Ping(); err != nil {
			index ++
			continue
		} else {
			return db
		}
		index ++
	}
	return nil
}