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
			url:          "/api/v1/region/default/cluster/test-cluster/create",
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
			url:          "/api/v1/region/default/cluster/test-cluster/delete",
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
			url:          "/api/v1/region/default/cluster/test-cluster/scale/up",
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
			url:          "/api/v1/region/default/cluster/test-cluster/scale/down",
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

	callback := &types.CallBack{}
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

	callback := &types.CallBack{}
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
			url:          "/api/v1/region/default/cluster/test-cluster/scale/up/callback",
			expectStatus: http.StatusOK,
		},
	}

	callback := &types.CallBack{}
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
			url:          "/api/v1/region/default/cluster/test-cluster/scale/down/callback",
			expectStatus: http.StatusOK,
		},
	}

	callback := &types.CallBack{}
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
