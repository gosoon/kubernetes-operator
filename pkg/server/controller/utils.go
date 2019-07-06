package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gosoon/glog"
)

// Response : http response func (no return http code)
func Response(w http.ResponseWriter, r *http.Request, data interface{}, httpCode int) {
	jsonByte, err := json.Marshal(data)
	if err != nil {
		glog.Errorf("marshal [%v] failed with err [%v]", data, err)
	}
	_, err = r.Cookie("WriteHeader")
	// if no WriteHeader
	if err != nil {
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpCode)
		w.Write(jsonByte)
	}
}

// SuccessResponse : http success response func
func SuccessResponse(w http.ResponseWriter, r *http.Request, data interface{}, httpCode int) {
	resp := make(map[string]interface{})
	resp["msg"] = "200"
	switch data.(type) {
	case error:
		resp["data"] = data.(error).Error()
	default:
		resp["data"] = data
	}
	jsonByte, err := json.Marshal(resp)
	if err != nil {
		glog.Errorf("marshal [%v] failed with err [%v]", data, err)
	}
	_, err = r.Cookie("WriteHeader")
	// if no WriteHeader
	if err != nil {
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpCode)
		w.Write(jsonByte)
	}
}

// FailedResponse : http failed response func
func FailedResponse(w http.ResponseWriter, r *http.Request, data interface{}, httpCode int) {
	resp := make(map[string]interface{})
	resp["msg"] = "500"
	switch data.(type) {
	case error:
		resp["data"] = data.(error).Error()
	case []error:
		errs := data.([]error)
		var s []string
		for _, err := range errs {
			s = append(s, fmt.Sprintf("%v", err))
		}
		resp["data"] = s
	default:
		resp["data"] = data
	}

	jsonByte, err := json.Marshal(resp)
	if err != nil {
		glog.Errorf("marshal [%v] failed with err [%v]", data, err)
	}
	_, err = r.Cookie("WriteHeader")
	// if no WriteHeader
	if err != nil {
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpCode)
		w.Write(jsonByte)
	}
}
