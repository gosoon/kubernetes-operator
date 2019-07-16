package cluster

import (
	"encoding/json"
	"net/http"

	"github.com/gosoon/kubernetes-operator/pkg/server/controller"
	"github.com/gosoon/kubernetes-operator/pkg/types"

	"github.com/gorilla/mux"
)

type cluster struct {
	opt *controller.Options
}

func New(opt *controller.Options) controller.Controller {
	return &cluster{opt: opt}
}

//
func (c *cluster) Register(router *mux.Router) {
	router = router.PathPrefix("/api/v1").Subrouter()

	// create
	router.Methods("POST").Path("/region/{region}/namespace/{namespace}/cluster/{name}/create").HandlerFunc(
		(c.createCluster))
	// scale up
	router.Methods("POST").Path("/region/{region}/namespace/{namespace}/cluster/{name}/scale/up").HandlerFunc(
		(c.scaleUpCluster))
	// scale down
	router.Methods("POST").Path("/region/{region}/namespace/{namespace}/cluster/{name}/scale/down").HandlerFunc(
		(c.scaleDownCluster))
	// scale up callback
	router.Methods("POST").Path("/region/{region}/namespace/{namespace}/cluster/{name}/scale/up/callback").HandlerFunc(
		(c.scaleUpCallback))
	// scale down callback
	router.Methods("POST").Path("/region/{region}/namespace/{namespace}/cluster/{name}/scale/down/callback").HandlerFunc(
		(c.scaleDownCallback))
	// delete
	router.Methods("DELETE").Path("/region/{region}/namespace/{namespace}/cluster/{name}/delete").HandlerFunc(
		(c.deleteCluster))
	// create callback
	router.Methods("POST").Path("/region/{region}/namespace/{namespace}/cluster/{name}/create/callback").HandlerFunc(
		(c.createClusterCallback))
	// delete callback
	router.Methods("POST").Path("/region/{region}/namespace/{namespace}/cluster/{name}/delete/callback").HandlerFunc(
		(c.deleteClusterCallback))
}

func (c *cluster) createCluster(w http.ResponseWriter, r *http.Request) {
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

func (c *cluster) deleteCluster(w http.ResponseWriter, r *http.Request) {
	// TODO: generate simple ecsv1.KubernetesCluster for web
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

	err = c.opt.Service.DeleteCluster(namespace, name, cluster)
	if err != nil {
		controller.FailedResponse(w, r, err, int(500))
		return
	}
	controller.SuccessResponse(w, r, `ok`, int(200))
}
