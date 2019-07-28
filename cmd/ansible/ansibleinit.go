package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gosoon/kubernetes-operator/pkg/client/clientset/versioned/scheme"
	"github.com/gosoon/kubernetes-operator/pkg/enum"
	"github.com/gosoon/kubernetes-operator/pkg/types"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

var cfgFile string

const (
	DeployEtcdCmd = `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--key-file ./private-key --become --become-user=root ansible/etcd.yml -vvvv`

	DeployMasterCmd = `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--key-file ./private-key --become --become-user=root ansible/master.yml -vvvv`

	DeployNodeCmd = `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--key-file ./private-key --become --become-user=root ansible/node.yml -vvvv`

	ScaleupNodeCmd = `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--key-file ./private-key --become --become-user=root ansible/scaleup-node.yml -vvvv`

	ScaledownNodeCmd = `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--key-file ./private-key --become --become-user=root ansible/scaledown-node.yml -vvvv`

	TerminatingCmd = `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--key-file ./private-key --become --become-user=root ansible/terminating.yml -vvvv`

	OperationEnv        = "OPERATION"
	ClusterNameEnv      = "CLUSTER_NAME"
	ClsuterNamespaceEnv = "CLUSTER_NAMESPACE"

	region                  = "config.region"
	server                  = "config.server"
	token                   = "config.token"
	timeout                 = "config.timeout"
	creatingCallbackPath    = "config.creatingCallbackPath"
	scalingUpCallbackPath   = "config.ccalingUpCallbackPath"
	scalingDownCallbackPath = "config.scalingDownCallbackPath"
	terminatingCallbackPath = "config.terminatingCallbackPath"
)

func init() {
	flag.StringVar(&cfgFile, "config", "", "config file")
	flag.Parse()
	initDefaultConfig()
}

// set default value
func initDefaultConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		fmt.Println("config is not found,exit.")
		os.Exit(1)
	}
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.SetDefault(region, "default")
	viper.SetDefault(server, "http://127.0.0.1:8000")
	viper.SetDefault(token, "")
	viper.SetDefault(timeout, 10*60) // 10 minutes
	viper.SetDefault(creatingCallbackPath, "/api/v1/region/{region}/cluster/{name}/create/callback")
	viper.SetDefault(scalingUpCallbackPath, "/api/v1/region/{region}/cluster/{name}/scaleup/callback")
	viper.SetDefault(scalingDownCallbackPath, "/api/v1/region/{region}/cluster/{name}/scaledown/callback")
	viper.SetDefault(terminatingCallbackPath, "/api/v1/region/{region}/cluster/{name}/delete/callback")
}

func main() {
	cmdStdout := make(chan string)
	cmdError := make(chan error, 1)

	operation := os.Getenv(OperationEnv)
	switch operation {
	case enum.KubeCreating:
		deployEtcdCmd := exec.Command("/bin/bash", "-c", `df -lh`)
		go execCmd(deployEtcdCmd, cmdStdout, cmdError)
		//go func() { packKubeCreatingCmd(cmdStdout, cmdError) }()
	case enum.KubeScalingUp:
		go func() { packKubeScalingUpCmd(cmdStdout, cmdError) }()
	case enum.KubeScalingDown:
		go func() { packKubeScalingDownCmd(cmdStdout, cmdError) }()
	case enum.KubeTerminating:
		go func() { packKubeTerminatingCmd(cmdStdout, cmdError) }()

	// do not know the callback path, exit
	default:
		fmt.Println("the OPERATION env not found,exit")
		os.Exit(1)
	}

	timeout := viper.GetInt(timeout)
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		callback(operation, "", errors.Errorf("the operation is timeout(%v)", timeout))
	case err := <-cmdError:
		callback(operation, "", err)
	case stdout := <-cmdStdout:
		callback(operation, stdout, nil)
	}
}

func execCmd(cmd *exec.Cmd, cmdStdout chan<- string, cmdError chan<- error) {
	// create command pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cmdError <- errors.Errorf("obtain stdout pipe for command failed with:%v\n", err)
		return
	}

	// exec command
	if err := cmd.Start(); err != nil {
		cmdError <- errors.Errorf("command start failed with:%v", err)
		return
	}

	// read all stdout
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		cmdError <- errors.Errorf("read stdout failed with:%v", err)
		return
	}

	// wait cmd exec finished
	if err := cmd.Wait(); err != nil {
		cmdError <- errors.Errorf("wait cmd exec finished failed with:%v", err)
		return
	}
	cmdStdout <- string(bytes)
}

func callback(operation string, stdout string, err error) {
	clusterName := os.Getenv(ClusterNameEnv)

	fmt.Println("stdout:", stdout)

	resp := types.Callback{
		Name:       os.Getenv(ClusterNameEnv),
		Namespace:  os.Getenv(ClsuterNamespaceEnv),
		Region:     viper.GetString(region),
		KubeConfig: "",
		Success:    true,
		Message:    stdout,
	}

	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
	}

	switch operation {
	case enum.KubeCreating:
		path := viper.GetString(creatingCallbackPath)
		packPath := packURLPath(path, map[string]string{"region": "", "name": clusterName})
		sendRequest(packPath)
	case enum.KubeScalingUp:
		path := viper.GetString(scalingUpCallbackPath)
		packPath := packURLPath(path, map[string]string{"region": "", "name": clusterName})
		sendRequest(packPath)
	case enum.KubeScalingDown:
		path := viper.GetString(scalingDownCallbackPath)
		packPath := packURLPath(path, map[string]string{"region": "", "name": clusterName})
		sendRequest(packPath)
	case enum.KubeTerminating:
		path := viper.GetString(terminatingCallbackPath)
		packPath := packURLPath(path, map[string]string{"region": "", "name": clusterName})
		sendRequest(packPath)
	}
}

// sendRequest is send request to controller
func sendRequest(path string) {
	c, err := rest.RESTClientFor(&rest.Config{
		Host: viper.GetString(server),
		ContentConfig: rest.ContentConfig{
			GroupVersion:         &schema.GroupVersion{Group: "", Version: ""},
			NegotiatedSerializer: scheme.Codecs.WithoutConversion(),
		},
		APIPath:     path,
		BearerToken: viper.GetString(token),
	})

	if err != nil {
		fmt.Println("new restclient failed with:", err)
		return
	}

	resp, err := c.Post().
		Do().
		Raw()

	if err != nil {
		fmt.Println("response failed with:", err)
		return
	}

	fmt.Println("response result is:", string(resp))
}

func packKubeCreatingCmd(cmdStdout chan<- string, cmdError chan<- error) {
	deployEtcdCmd := exec.Command("/bin/bash", "-c", `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--become --become-user=root etcd.yml`)
	execCmd(deployEtcdCmd, cmdStdout, cmdError)
	if len(cmdError) != 0 {
		return
	}

	deployMasterCmd := exec.Command("/bin/bash", "-c", `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--become --become-user=root master.yml`)
	execCmd(deployMasterCmd, cmdStdout, cmdError)
	if len(cmdError) != 0 {
		return
	}

	deployNodeCmd := exec.Command("/bin/bash", "-c", `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--become --become-user=root node.yml`)
	execCmd(deployNodeCmd, cmdStdout, cmdError)
}

func packKubeScalingUpCmd(cmdStdout chan<- string, cmdError chan<- error) {
	scalingUpCmd := exec.Command("/bin/bash", "-c", `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--become --become-user=root scaleup-node.yml`)

	execCmd(scalingUpCmd, cmdStdout, cmdError)
}

func packKubeScalingDownCmd(cmdStdout chan<- string, cmdError chan<- error) {
	scalingDownCmd := exec.Command("/bin/bash", "-c", `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--become --become-user=root scaledown-node.yml`)

	execCmd(scalingDownCmd, cmdStdout, cmdError)
}

func packKubeTerminatingCmd(cmdStdout chan<- string, cmdError chan<- error) {
	terminatingCmd := exec.Command("/bin/bash", "-c", `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--become --become-user=root terminating.yml`)

	execCmd(terminatingCmd, cmdStdout, cmdError)
}

func packURLPath(tpl string, args map[string]string) string {
	if args == nil {
		return tpl
	}
	if args["region"] == "" {
		args["region"] = viper.GetString(region)
	}
	for k, v := range args {
		tpl = strings.Replace(tpl, "{"+k+"}", url.QueryEscape(v), 1)
	}
	return tpl
}
