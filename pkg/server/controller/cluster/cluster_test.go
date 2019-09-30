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

package cluster_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gosoon/kubernetes-operator/pkg/server/controller"
	"github.com/gosoon/kubernetes-operator/pkg/server/controller/cluster"
	"github.com/gosoon/kubernetes-operator/pkg/test/mock_service"
	"github.com/gosoon/kubernetes-operator/pkg/types"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// create
func TestCreateCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockInterface(ctrl)
	mockService.EXPECT().CreateCluster(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	router := mux.NewRouter()
	//pathPrefix := "/api/v1"
	cluster.New(&controller.Options{Service: mockService}).Register(router.PathPrefix("").Subrouter())
	testCases := []struct {
		method       string
		url          string
		body         string
		expectStatus int
		expectBody   string
	}{
		{
			method:       "POST",
			url:          "/api/v1/region/default/cluster/test-cluster",
			body:         "",
			expectStatus: http.StatusOK,
		},
	}

	client := &types.EcsClient{}
	body, _ := json.Marshal(client)

	for _, test := range testCases {
		req := httptest.NewRequest(test.method, test.url, strings.NewReader(string(body)))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		//t.Log(req)
		//t.Log(w)
		if !assert.Equal(t, test.expectStatus, w.Code) {
			t.Errorf("expect %v, got %v", test.expectStatus, w.Code)
		}
	}
}

// delete
func TestDeleteCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockInterface(ctrl)
	mockService.EXPECT().DeleteCluster(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	router := mux.NewRouter()
	cluster.New(&controller.Options{Service: mockService}).Register(router.PathPrefix("").Subrouter())
	testCases := []struct {
		method       string
		url          string
		body         string
		expectStatus int
		expectBody   string
	}{
		// delete
		{
			method:       "DELETE",
			url:          "/api/v1/region/default/cluster/test-cluster",
			expectStatus: http.StatusOK,
		},
	}

	client := &types.EcsClient{}
	body, _ := json.Marshal(client)

	for _, test := range testCases {
		req := httptest.NewRequest(test.method, test.url, strings.NewReader(string(body)))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if !assert.Equal(t, test.expectStatus, w.Code) {
			t.Errorf("expect %v, got %v", test.expectStatus, w.Code)
		}
	}
}

// scale up
func TestScaleUpCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockInterface(ctrl)
	mockService.EXPECT().ScaleUp(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	router := mux.NewRouter()
	cluster.New(&controller.Options{Service: mockService}).Register(router.PathPrefix("").Subrouter())
	testCases := []struct {
		method       string
		url          string
		body         string
		expectStatus int
		expectBody   string
	}{
		{
			method:       "POST",
			url:          "/api/v1/region/default/cluster/test-cluster/scaleup",
			expectStatus: http.StatusOK,
		},
	}

	client := &types.EcsClient{}
	body, _ := json.Marshal(client)

	for _, test := range testCases {
		req := httptest.NewRequest(test.method, test.url, strings.NewReader(string(body)))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if !assert.Equal(t, test.expectStatus, w.Code) {
			t.Errorf("expect %v, got %v", test.expectStatus, w.Code)
		}
	}
}

// scale down
func TestScaleDownCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockInterface(ctrl)
	mockService.EXPECT().ScaleDown(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	router := mux.NewRouter()
	cluster.New(&controller.Options{Service: mockService}).Register(router.PathPrefix("").Subrouter())
	testCases := []struct {
		method       string
		url          string
		body         string
		expectStatus int
		expectBody   string
	}{
		{
			method:       "POST",
			url:          "/api/v1/region/default/cluster/test-cluster/scaledown",
			expectStatus: http.StatusOK,
		},
	}

	client := &types.EcsClient{}
	body, _ := json.Marshal(client)

	for _, test := range testCases {
		req := httptest.NewRequest(test.method, test.url, strings.NewReader(string(body)))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if !assert.Equal(t, test.expectStatus, w.Code) {
			t.Errorf("expect %v, got %v", test.expectStatus, w.Code)
		}
	}
}

// create callback
func TestCreateClusterCallback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockInterface(ctrl)
	mockService.EXPECT().CreateClusterCallback(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	router := mux.NewRouter()
	cluster.New(&controller.Options{Service: mockService}).Register(router.PathPrefix("").Subrouter())
	testCases := []struct {
		method       string
		url          string
		body         string
		expectStatus int
		expectBody   string
	}{
		{
			method:       "POST",
			url:          "/api/v1/region/default/cluster/test-cluster/create/callback",
			expectStatus: http.StatusOK,
		},
	}

	callback := &types.Callback{}
	body, _ := json.Marshal(callback)

	for _, test := range testCases {
		req := httptest.NewRequest(test.method, test.url, strings.NewReader(string(body)))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if !assert.Equal(t, test.expectStatus, w.Code) {
			t.Errorf("expect %v, got %v", test.expectStatus, w.Code)
		}
	}
}

// delete callback
func TestDeleteClusterCallback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockInterface(ctrl)
	mockService.EXPECT().DeleteClusterCallback(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	router := mux.NewRouter()
	cluster.New(&controller.Options{Service: mockService}).Register(router.PathPrefix("").Subrouter())
	testCases := []struct {
		method       string
		url          string
		body         string
		expectStatus int
		expectBody   string
	}{
		{
			method:       "POST",
			url:          "/api/v1/region/default/cluster/test-cluster/delete/callback",
			expectStatus: http.StatusOK,
		},
	}

	callback := &types.Callback{}
	body, _ := json.Marshal(callback)

	for _, test := range testCases {
		req := httptest.NewRequest(test.method, test.url, strings.NewReader(string(body)))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if !assert.Equal(t, test.expectStatus, w.Code) {
			t.Errorf("expect %v, got %v", test.expectStatus, w.Code)
		}
	}
}

// scale up callback
func TestScaleUpClusterCallback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockInterface(ctrl)
	mockService.EXPECT().ScaleUpCallback(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	router := mux.NewRouter()
	cluster.New(&controller.Options{Service: mockService}).Register(router.PathPrefix("").Subrouter())
	testCases := []struct {
		method       string
		url          string
		body         string
		expectStatus int
		expectBody   string
	}{
		{
			method:       "POST",
			url:          "/api/v1/region/default/cluster/test-cluster/scaleup/callback",
			expectStatus: http.StatusOK,
		},
	}

	callback := &types.Callback{}
	body, _ := json.Marshal(callback)

	for _, test := range testCases {
		req := httptest.NewRequest(test.method, test.url, strings.NewReader(string(body)))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if !assert.Equal(t, test.expectStatus, w.Code) {
			t.Errorf("expect %v, got %v", test.expectStatus, w.Code)
		}
	}
}

// scale down callback
func TestScaleDownClusterCallback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_service.NewMockInterface(ctrl)
	mockService.EXPECT().ScaleDownCallback(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	router := mux.NewRouter()
	cluster.New(&controller.Options{Service: mockService}).Register(router.PathPrefix("").Subrouter())
	testCases := []struct {
		method       string
		url          string
		body         string
		expectStatus int
		expectBody   string
	}{
		{
			method:       "POST",
			url:          "/api/v1/region/default/cluster/test-cluster/scaledown/callback",
			expectStatus: http.StatusOK,
		},
	}

	callback := &types.Callback{}
	body, _ := json.Marshal(callback)

	for _, test := range testCases {
		req := httptest.NewRequest(test.method, test.url, strings.NewReader(string(body)))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if !assert.Equal(t, test.expectStatus, w.Code) {
			t.Errorf("expect %v, got %v", test.expectStatus, w.Code)
		}
	}
}
