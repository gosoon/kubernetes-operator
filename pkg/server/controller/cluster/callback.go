package cluster

import (
	"encoding/json"
	"net/http"

	"github.com/gosoon/kubernetes-operator/pkg/server/controller"
	"github.com/gosoon/kubernetes-operator/pkg/types"

	"github.com/gorilla/mux"
)

func (c *cluster) createClusterCallback(w http.ResponseWriter, r *http.Request) {
	callback := &types.CallBack{}
	err := json.NewDecoder(r.Body).Decode(callback)
	if err != nil {
		errMsg := "Data format error,post data unable to decode."
		controller.FailedResponse(w, r, errMsg, int(500))
		return
	}
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.CreateClusterCallback(namespace, name, callback)
	if err != nil {
		controller.FailedResponse(w, r, err, int(500))
		return
	}
	controller.SuccessResponse(w, r, `ok`, int(200))
}

func (c *cluster) scaleUpCallback(w http.ResponseWriter, r *http.Request) {
	callback := &types.CallBack{}
	err := json.NewDecoder(r.Body).Decode(callback)
	if err != nil {
		errMsg := "Data format error,post data unable to decode."
		controller.FailedResponse(w, r, errMsg, int(500))
		return
	}
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.ScaleUpCallback(namespace, name, callback)
	if err != nil {
		controller.FailedResponse(w, r, err, int(500))
		return
	}
	controller.SuccessResponse(w, r, `ok`, int(200))
}

func (c *cluster) scaleDownCallback(w http.ResponseWriter, r *http.Request) {
	callback := &types.CallBack{}
	err := json.NewDecoder(r.Body).Decode(callback)
	if err != nil {
		errMsg := "Data format error,post data unable to decode."
		controller.FailedResponse(w, r, errMsg, int(500))
		return
	}
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.ScaleDownCallback(namespace, name, callback)
	if err != nil {
		controller.FailedResponse(w, r, err, int(500))
		return
	}
	controller.SuccessResponse(w, r, `ok`, int(200))
}

func (c *cluster) deleteClusterCallback(w http.ResponseWriter, r *http.Request) {
	callback := &types.CallBack{}
	err := json.NewDecoder(r.Body).Decode(callback)
	if err != nil {
		errMsg := "Data format error,post data unable to decode."
		controller.FailedResponse(w, r, errMsg, int(500))
		return
	}
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.DeleteClusterCallback(namespace, name, callback)
	if err != nil {
		controller.FailedResponse(w, r, err, int(500))
		return
	}
	controller.SuccessResponse(w, r, `ok`, int(200))
}
