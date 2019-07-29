module github.com/gosoon/kubernetes-operator

go 1.12

require (
	github.com/deckarep/golang-set v1.7.1
	github.com/golang/mock v1.1.1
	github.com/gorilla/mux v1.7.3
	github.com/gosoon/glog v0.0.0-20180521124921-a5fbfb162a81
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/resouer/k8s-controller-custom-resource v0.0.0-20180915125134-dbc1c9320e34
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.3.0
	k8s.io/api v0.0.0-20190718062839-c8a0b81cb10e
	k8s.io/apimachinery v0.0.0-20190717022731-0bb8574e0887
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/yaml v1.1.0
)
