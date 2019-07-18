#!/bin/bash

mockgen github.com/gosoon/kubernetes-operator/pkg/server/service Interface >mock_service/service.go
mockgen github.com/gosoon/kubernetes-operator/pkg/client/clientset/versioned Interface >mock_versioned/service.go
mockgen k8s.io/client-go/kubernetes Interface >mock_kubernetes/service.go
