package args

import "flag"

func InitArgs(args *Arguments) {
	args.Cfg = flag.String("file", "../config/config.toml", "config file")
	flag.Parse()
}

func InitClientArgs(args *ClientArgs) {
	args.CfgFile = flag.String("file", "../config/client.toml", "config file")
	args.Src = flag.String("src", "", "src cluster name")
	args.Dst = flag.String("dst", "", "dst cluster name")
	args.Type = flag.String("type", "full", "support: schema, full")
	args.Db = flag.String("db", "", "db:tb,tb2")
	args.Operator = flag.String("ops", "create", "create/show")
	args.Role = flag.String("role", "task", "task/machine")
	args.UUID = flag.Int("uuid", -1, "task uuid")
	flag.Parse()
}