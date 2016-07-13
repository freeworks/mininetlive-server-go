package common

type Resp struct {
	Ret  int64       `form:"ret" json:"ret"`
	Msg  string      `form:"msg" json:"msg"`
	Data interface{} `form:"data" json:"data"`
}
