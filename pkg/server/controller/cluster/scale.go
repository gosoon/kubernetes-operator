package cluster

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gosoon/kubernetes-operator/pkg/server/controller"
	"github.com/gosoon/kubernetes-operator/pkg/types"
)

func (c *cluster) scaleUpCluster(w http.ResponseWriter, r *http.Request) {
	cluster := &types.EcsClient{}
	err := json.NewDecoder(r.Body).Decode(cluster)
	if err != nil {
		errMsg := "Data format error,post data unable to decode."
		controller.FailedResponse(w, r, errMsg, int(500))
		return
	}
	//region := mux.Vars(r)["region"]
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.ScaleUp(namespace, name, cluster)
	if err != nil {
		controller.FailedResponse(w, r, err, int(500))
		return
	}
	controller.SuccessResponse(w, r, err, int(200))
}

func (c *cluster) scaleDownCluster(w http.ResponseWriter, r *http.Request) {
	cluster := &types.EcsClient{}
	err := json.NewDecoder(r.Body).Decode(cluster)
	if err != nil {
		errMsg := "Data format error,post data unable to decode."
		controller.FailedResponse(w, r, errMsg, int(500))
		return
	}
	//region := mux.Vars(r)["region"]
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.CreateCluster(namespace, name, cluster)
	if err != nil {
		controller.FailedResponse(w, r, err, int(500))
		return
	}
	controller.SuccessResponse(w, r, err, int(200))
}
