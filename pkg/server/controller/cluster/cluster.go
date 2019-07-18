package cluster

import (
	"encoding/json"
	"net/http"

	"github.com/gosoon/kubernetes-operator/pkg/server/controller"
	"github.com/gosoon/kubernetes-operator/pkg/types"

	"github.com/gorilla/mux"
)

const (
	namespace = "default"
)

type cluster struct {
	opt *controller.Options
}

func New(opt *controller.Options) controller.Controller {
	return &cluster{opt: opt}
}

func (c *cluster) Register(router *mux.Router) {
	router = router.PathPrefix("/api/v1").Subrouter()

	// create
	router.Methods("POST").Path("/region/{region}/cluster/{name}/create").HandlerFunc(
		(c.createCluster))
	// scale up
	router.Methods("POST").Path("/region/{region}/cluster/{name}/scale/up").HandlerFunc(
		(c.scaleUpCluster))
	// scale down
	router.Methods("POST").Path("/region/{region}/cluster/{name}/scale/down").HandlerFunc(
		(c.scaleDownCluster))
	// scale up callback
	router.Methods("POST").Path("/region/{region}/cluster/{name}/scale/up/callback").HandlerFunc(
		(c.scaleUpCallback))
	// scale down callback
	router.Methods("POST").Path("/region/{region}/cluster/{name}/scale/down/callback").HandlerFunc(
		(c.scaleDownCallback))
	// delete
	router.Methods("DELETE").Path("/region/{region}/cluster/{name}/delete").HandlerFunc(
		(c.deleteCluster))
	// create callback
	router.Methods("POST").Path("/region/{region}/cluster/{name}/create/callback").HandlerFunc(
		(c.createClusterCallback))
	// delete callback
	router.Methods("POST").Path("/region/{region}/cluster/{name}/delete/callback").HandlerFunc(
		(c.deleteClusterCallback))

	// get current operate logs
	router.Methods("GET").Path("/region/{region}/cluster/{name}/logs").HandlerFunc(
		(c.getClusterOperationLogs))

}

func (c *cluster) createCluster(w http.ResponseWriter, r *http.Request) {
	cluster := &types.EcsClient{}
	err := json.NewDecoder(r.Body).Decode(cluster)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	region := mux.Vars(r)["region"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.CreateCluster(region, namespace, name, cluster)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}

func (c *cluster) deleteCluster(w http.ResponseWriter, r *http.Request) {
	cluster := &types.EcsClient{}
	err := json.NewDecoder(r.Body).Decode(cluster)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	region := mux.Vars(r)["region"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.DeleteCluster(region, namespace, name, cluster)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}
