package cluster

import (
	"net/http"

	"github.com/gosoon/kubernetes-operator/pkg/server/controller"

	"github.com/gorilla/mux"
)

func (c *cluster) getClusterOperationLogs(w http.ResponseWriter, r *http.Request) {

	region := mux.Vars(r)["region"]
	name := mux.Vars(r)["name"]

	logs, err := c.opt.Service.GetClusterOperationLogs(region, namespace, name)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, logs)
}
