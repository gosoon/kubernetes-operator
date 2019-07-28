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
	router.Methods("POST").Path("/region/{region}/cluster/{name}").HandlerFunc(
		(c.createCluster))
	// scale up
	router.Methods("POST").Path("/region/{region}/cluster/{name}/scaleup").HandlerFunc(
		(c.scaleUpCluster))
	// scale down
	router.Methods("POST").Path("/region/{region}/cluster/{name}/scaledown").HandlerFunc(
		(c.scaleDownCluster))
	// scale up callback
	router.Methods("POST").Path("/region/{region}/cluster/{name}/scaleup/callback").HandlerFunc(
		(c.scaleUpCallback))
	// scale down callback
	router.Methods("POST").Path("/region/{region}/cluster/{name}/scaledown/callback").HandlerFunc(
		(c.scaleDownCallback))
	// delete
	router.Methods("DELETE").Path("/region/{region}/cluster/{name}").HandlerFunc(
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

// createCluster
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

// deleteCluster
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

// scaleUpCluster
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

// scaleDownCluster
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

// createClusterCallback
func (c *cluster) createClusterCallback(w http.ResponseWriter, r *http.Request) {
	callback := &types.Callback{}
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

// scaleUpCallback
func (c *cluster) scaleUpCallback(w http.ResponseWriter, r *http.Request) {
	callback := &types.Callback{}
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

// scaleDownCallback
func (c *cluster) scaleDownCallback(w http.ResponseWriter, r *http.Request) {
	callback := &types.Callback{}
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

// deleteClusterCallback
func (c *cluster) deleteClusterCallback(w http.ResponseWriter, r *http.Request) {
	callback := &types.Callback{}
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

// getClusterOperationLogs
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
