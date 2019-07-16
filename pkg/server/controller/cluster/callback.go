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
		controller.BadRequest(w, r, err)
		return
	}
	region := mux.Vars(r)["region"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.CreateClusterCallback(region, namespace, name, callback)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}

func (c *cluster) scaleUpCallback(w http.ResponseWriter, r *http.Request) {
	callback := &types.CallBack{}
	err := json.NewDecoder(r.Body).Decode(callback)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	region := mux.Vars(r)["region"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.ScaleUpCallback(region, namespace, name, callback)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}

func (c *cluster) scaleDownCallback(w http.ResponseWriter, r *http.Request) {
	callback := &types.CallBack{}
	err := json.NewDecoder(r.Body).Decode(callback)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	region := mux.Vars(r)["region"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.ScaleDownCallback(region, namespace, name, callback)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}

func (c *cluster) deleteClusterCallback(w http.ResponseWriter, r *http.Request) {
	callback := &types.CallBack{}
	err := json.NewDecoder(r.Body).Decode(callback)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	region := mux.Vars(r)["region"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.DeleteClusterCallback(region, namespace, name, callback)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}
