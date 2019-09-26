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

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/installer/server/controller"
	"github.com/gosoon/kubernetes-operator/pkg/installer/server/middleware"

	"github.com/gorilla/mux"
)

// cluster implements the controller interface.
type cluster struct {
	opt *controller.Options
}

// New is create a cluster object.
func New(opt *controller.Options) controller.Controller {
	return &cluster{opt: opt}
}

// Register is register the routes to router
func (c *cluster) Register(router *mux.Router) {
	router = router.PathPrefix("/api/v1").Subrouter()

	// create
	router.Methods("POST").Path("/region/{region}/cluster/{name}").HandlerFunc(
		middleware.Authenticate((c.createCluster)))
}

// createCluster
func (c *cluster) createCluster(w http.ResponseWriter, r *http.Request) {
	cluster := &ecsv1.KubernetesCluster{}
	err := json.NewDecoder(r.Body).Decode(cluster)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	region := mux.Vars(r)["region"]
	name := mux.Vars(r)["name"]

	err = c.opt.Service.CreateCluster(r.Context(), region, name, cluster)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}
