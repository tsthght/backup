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
	flag.Parse()
}