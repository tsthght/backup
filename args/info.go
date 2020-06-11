package args

type Arguments struct {
	Cfg *string
}

type ClientArgs struct {
	CfgFile *string
	Src *string
	Dst *string
	Type *string
	Db *string
	Operator *string
	Role *string
	UUID *int
}