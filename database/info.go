package database

type MGRInfo struct {
	Hosts []string
	Port int
	WriteIndex int
}

type TiDBInfo struct {
	Hosts []string
	Port int
}

type UserInfo struct {
	Username string
	Password string
	Port string
	Database string
}