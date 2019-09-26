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
package controller

import (
	"github.com/gorilla/mux"
	"github.com/gosoon/code-generator/_examples/server/service"
	"k8s.io/client-go/kubernetes"
)

// Options contains the config by controller
type Options struct {
	KubeClientset kubernetes.Interface
	Service       service.Interface
}

// Controller helps register to router.
type Controller interface {
	Register(router *mux.Router)
}
