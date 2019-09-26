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
package server

import (
	"net/http"
	"time"

	ctrl "github.com/gosoon/kubernetes-operator/pkg/installer/server/controller"
	"github.com/gosoon/kubernetes-operator/pkg/installer/server/controller/cluster"
	"github.com/gosoon/kubernetes-operator/pkg/installer/server/service"

	"github.com/gorilla/mux"
)

// Server helps start a http server.
type Server interface {
	http.Handler
	ListenAndServe() error
}

// Options contains the config required by server
type Options struct {
	CtrlOptions *ctrl.Options
	ListenAddr  string
}

// server implements the Server interface.
type server struct {
	opt    Options
	router *mux.Router
}

// New is create a server object.
func New(opt Options) Server {
	// init service
	options := &service.Options{
		KubeClientset: opt.CtrlOptions.KubeClientset,
	}

	opt.CtrlOptions.Service = service.New(options)

	router := mux.NewRouter().StrictSlash(true)
	cluster.New(opt.CtrlOptions).Register(router)

	return &server{
		opt:    opt,
		router: router,
	}
}

// ServeHTTP dispatches the handler registered in the matched route.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// ListenAndServe start a http server.
func (s *server) ListenAndServe() error {
	server := &http.Server{
		Handler: s.router,
		Addr:    s.opt.ListenAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
