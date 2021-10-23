package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type RestEndpoint struct {
	http.Handler
	*http.ServeMux
}

func NewRestEndpoint(funcs map[string]http.HandlerFunc) *RestEndpoint {
	ret := &RestEndpoint{
		ServeMux: http.NewServeMux(),
	}
	ret.ServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		funcs[r.Method](w, r)
	})
	return ret
}

func (e *RestEndpoint) Json(w http.ResponseWriter, status int, v interface{}) {
	var jsonData []byte
	var err error

	// if config.Get().API.CompressResponses {
	// jsonData, err = json.Marshal(v)
	// } else {
	jsonData, err = json.MarshalIndent(v, "", "  ")
	// }
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err.Error())
	} else {
		w.WriteHeader(status)
		fmt.Fprintf(w, "%s", jsonData)
	}
}

func (e *RestEndpoint) Error(w http.ResponseWriter, status int, code, message string) {
	e.Json(w, status, ErrorResponse{
		Code:    code,
		Message: message,
	})
}

func (e *RestEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.ServeMux.ServeHTTP(w, r)
}
