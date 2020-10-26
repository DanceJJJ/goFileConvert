package response

import (
	"encoding/json"
	"net/http"
)

type RetData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func MJson(w http.ResponseWriter, code int, msg string, data interface{}) {

	//设置返回头信息
	w.Header().Set("Content-Type", "application/json;charset=utf8")

	res := RetData{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	retJson, _ := json.Marshal(res)

	w.Write(retJson)
	return
}
