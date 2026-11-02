package jsonify

import "encoding/json"

type okStruct struct {
	Error  string `json:"error"`
	Status string `json:"status"`
	Resp   string `json:"resp"`
}

func GoodResp(resp string) string {
	var jsonRespStru okStruct
	jsonRespStru.Error = ""
	jsonRespStru.Status = "ok"
	jsonRespStru.Resp = resp
	result, _ := json.Marshal(jsonRespStru)
	return string(result)

}
