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
		controller.BadRequest(w, r, err)
		return
	}
	region := mux.Vars(r)["region"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.ScaleUp(region, namespace, name, cluster)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}

func (c *cluster) scaleDownCluster(w http.ResponseWriter, r *http.Request) {
	cluster := &types.EcsClient{}
	err := json.NewDecoder(r.Body).Decode(cluster)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	region := mux.Vars(r)["region"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.ScaleDown(region, namespace, name, cluster)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}
