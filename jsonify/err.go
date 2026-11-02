package jsonify

import (
	"encoding/json"
	"log/slog"
)

type errStruct struct {
	Error  string `json:"error"`
	Status string `json:"status"`
	Resp   string `json:"resp"`
}

// 将错误转换为 JSON, 以交给 HTTP 请求,返回 json: error & status:bad
func Grateful_Error(errStr string) string {
	var theErr errStruct
	theErr.Status = "bad"
	theErr.Resp = ""
	theErr.Error = errStr
	jsonContByte, err := json.Marshal(theErr)
	if err != nil {
		slog.Error("[err.go] An error occurred", "error", err)
	}
	return string(jsonContByte)
}
