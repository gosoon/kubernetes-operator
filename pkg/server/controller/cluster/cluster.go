package cluster

import (
	"encoding/json"
	"net/http"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/server/controller"

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

	router.Methods("POST").Path("/region/{region}/namespace/{namespace}/cluster/{name}/create").HandlerFunc(
		(c.createCluster))
	router.Methods("DELETE").Path("/region/{region}/namespace/{namespace}/cluster/{name}/delete").HandlerFunc(
		(c.deleteCluster))
	router.Methods("POST").Path("/region/{region}/namespace/{namespace}/cluster/{name}/create/callback").HandlerFunc(
		(c.createClusterCallback))
	router.Methods("POST").Path("/region/{region}/namespace/{namespace}/cluster/{name}/delete/callback").HandlerFunc(
		(c.deleteClusterCallback))
}

func (c *cluster) createCluster(w http.ResponseWriter, r *http.Request) {
	kubernetesCluster := &ecsv1.KubernetesCluster{}
	err := json.NewDecoder(r.Body).Decode(kubernetesCluster)
	if err != nil {
		errMsg := "Data format error,post data unable to decode."
		controller.FailedResponse(w, r, errMsg, int(500))
		return
	}
	//region := mux.Vars(r)["region"]
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.CreateCluster(namespace, name, kubernetesCluster)
	if err != nil {
		controller.FailedResponse(w, r, err, int(500))
		return
	}
	controller.SuccessResponse(w, r, err, int(200))
}

func (c *cluster) deleteCluster(w http.ResponseWriter, r *http.Request) {
	cluster := &ecsv1.KubernetesCluster{}
	err := json.NewDecoder(r.Body).Decode(cluster)
	if err != nil {
		errMsg := "Data format error,post data unable to decode."
		controller.FailedResponse(w, r, errMsg, int(500))
		return
	}
	//region := mux.Vars(r)["region"]
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.DeleteCluster(namespace, name)
	if err != nil {
		controller.FailedResponse(w, r, err, int(500))
		return
	}
	controller.SuccessResponse(w, r, `ok`, int(200))
}

func (c *cluster) createClusterCallback(w http.ResponseWriter, r *http.Request) {

}

func (c *cluster) deleteClusterCallback(w http.ResponseWriter, r *http.Request) {

}
