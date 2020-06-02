package args

import "flag"

func InitArgs(args *Arguments) {
	args.Cfg = flag.String("file", "../config/config.toml", "config file")
	flag.Parse()
}
