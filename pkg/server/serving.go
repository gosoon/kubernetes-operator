package server

import (
	"net/http"
	"time"

	ctrl "github.com/gosoon/kubernetes-operator/pkg/server/controller"
	"github.com/gosoon/kubernetes-operator/pkg/server/controller/cluster"
	"github.com/gosoon/kubernetes-operator/pkg/server/service"

	"github.com/gorilla/mux"
)

type Server interface {
	http.Handler
	ListenAndServe() error
}

type Options struct {
	CtrlOptions *ctrl.Options
	ListenAddr  string
}

type server struct {
	opt    Options
	router *mux.Router
}

func New(opt Options) Server {
	// init service
	opt.CtrlOptions.Service = service.New(&service.Options{KubernetesClusterClientset: opt.CtrlOptions.KubernetesClusterClientset})

	router := mux.NewRouter().StrictSlash(true)
	cluster.New(opt.CtrlOptions).Register(router)

	return &server{
		opt:    opt,
		router: router,
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

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
