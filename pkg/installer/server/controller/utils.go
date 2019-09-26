/*
 * Copyright 2019 gosoon.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package controller

import (
	"encoding/json"
	"net/http"

	"k8s.io/klog"
)

type commResp struct {
	Code    string      `json:"code"`
	Message interface{} `json:"message"`
}

// OK reply
func OK(w http.ResponseWriter, r *http.Request, message string) {
	Response(w, r, http.StatusOK, message)
}

// ResourceNotFound will return an error message indicating that the resource is not exist
func ResourceNotFound(w http.ResponseWriter, r *http.Request, message string) {
	Response(w, r, http.StatusNotFound, message)
}

// BadRequest will return an error message indicating that the request is invalid
func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	Response(w, r, http.StatusBadRequest, err.Error())
}

// Forbidden will block user access the resource, not authorized
func Forbidden(w http.ResponseWriter, r *http.Request, err error) {
	Response(w, r, http.StatusForbidden, err.Error())
}

// Unauthorized will block user access the api, not login
func Unauthorized(w http.ResponseWriter, r *http.Request, err error) {
	Response(w, r, http.StatusUnauthorized, err.Error())
}

// InternalError will return an error message indicating that the something is error inside the controller
func InternalError(w http.ResponseWriter, r *http.Request, err error) {
	Response(w, r, http.StatusInternalServerError, err.Error())
}

// ServiceUnavailable will return an error message indicating that the service is not available now
func ServiceUnavailable(w http.ResponseWriter, r *http.Request, err error) {
	Response(w, r, http.StatusServiceUnavailable, err.Error())
}

// Conflict xxx
func Conflict(w http.ResponseWriter, r *http.Request, err error) {
	Response(w, r, http.StatusConflict, err.Error())
}

// NotAcceptable xxx
func NotAcceptable(w http.ResponseWriter, r *http.Request, err error) {
	Response(w, r, http.StatusNotAcceptable, err.Error())
}

// Response : http response func (no return http code)
func Response(w http.ResponseWriter, r *http.Request, httpCode int, message interface{}) {
	resp := commResp{
		Code:    http.StatusText(httpCode),
		Message: message,
	}

	jsonByte, err := json.Marshal(resp)
	if err != nil {
		klog.Errorf("marshal [%v] failed with err [%v]", resp, err)
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
