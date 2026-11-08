package jsonify

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
)

type Response struct {
	Error  string `json:"error"`
	Status string `json:"status"`
	Resp   string `json:"resp"`
}

// WriteJSON 以 application/json 返回统一响应结构
func WriteJSON(w http.ResponseWriter, httpStatus int, status string, errStr string, resp string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatus)
	body, err := json.Marshal(Response{Error: errStr, Status: status, Resp: resp})
	if err != nil {
		slog.Error("[jsonify/writer.go] An error occurred", "error", err)
		return
	}
	_, _ = w.Write(append(body, '\n'))
}

// WriteOK 返回 status=ok 的 200 响应
func WriteOK(w http.ResponseWriter, resp string) {
	WriteJSON(w, http.StatusOK, "ok", "", resp)
}

// WriteError 根据错误返回对应的 HTTP 状态码和 status=bad 响应
func WriteError(w http.ResponseWriter, err error) {
	WriteJSON(w, ErrorStatus(err), "bad", err.Error(), "")
}

// WriteErrorMessage 返回指定消息的 status=bad 响应
func WriteErrorMessage(w http.ResponseWriter, httpStatus int, errStr string) {
	WriteJSON(w, httpStatus, "bad", errStr, "")
}

// ErrorStatus 将错误映射为 HTTP 状态码
func ErrorStatus(err error) int {
	if errors.Is(err, os.ErrNotExist) {
		return http.StatusNotFound
	}
	if err.Error() == "Game Over" {
		return http.StatusConflict
	}
	return http.StatusBadRequest
}
