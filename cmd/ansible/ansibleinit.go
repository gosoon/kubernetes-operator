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

package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gosoon/kubernetes-operator/pkg/client/clientset/versioned/scheme"
	"github.com/gosoon/kubernetes-operator/pkg/enum"
	"github.com/gosoon/kubernetes-operator/pkg/types"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	// cmd
	DeployEtcdCmd = `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--key-file ./private-key --become --become-user=root ansible/etcd.yml -vv`

	DeployMasterCmd = `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--key-file ./private-key --become --become-user=root ansible/master.yml -vv`

	DeployNodeCmd = `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--key-file ./private-key --become --become-user=root ansible/node.yml -vv`

	ScaleUpCmd = `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--key-file ./private-key --become --become-user=root ansible/scaleup-node.yml -vv`

	ScaleDownCmd = `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--key-file ./private-key --become --become-user=root ansible/scaledown-node.yml -vv`

	TerminatingCmd = `ansible-playbook -i ansible/inventory/production/hosts.yaml \
				--key-file ./private-key --become --become-user=root ansible/terminating.yml -vv`

	KubeconfigCmd = `ssh -i ./private-key root@%v cat ~/.kube/config`

	// env
	OperationEnv            = "OPERATION"
	ClusterNameEnv          = "CLUSTER_NAME"
	ClsuterNamespaceEnv     = "CLUSTER_NAMESPACE"
	MasterHostsEnv          = "MASTER_HOSTS"
	ExternalLoadBalancerEnv = "MASTER_VIP"
	EtcdHostsEnv            = "ETCD_HOSTS"
	NodeHostsEnv            = "NODE_HOSTS"
	HostsYAMLEnv            = "HOSTS_YAML"
	PrivateKeyEnv           = "PRIVATE_KEY"

	// config
	region                  = "config.region"
	server                  = "config.server"
	token                   = "config.token"
	timeout                 = "config.timeout"
	creatingCallbackPath    = "config.creatingCallbackPath"
	scalingUpCallbackPath   = "config.ccalingUpCallbackPath"
	scalingDownCallbackPath = "config.scalingDownCallbackPath"
	terminatingCallbackPath = "config.terminatingCallbackPath"

	// FileName is save env val file
	EnvFileName        = "./scripts/deploy/hosts_env"
	HostsYAMLFileName  = "./ansible/inventory/production/hosts.yaml"
	PrivateKeyFileName = "./private-key"
)

var (
	cfgFile string

	MasterHostsVal, ExternalLoadBalancerVal, NodeHostsVal, EtcdHostsVal, OperationVal string
	ClusterNameVal, ClusterNamespaceVal, HostsYAMLVal, PrivateKeyVal                  string
)

func init() {
	flag.StringVar(&cfgFile, "config", "", "config file")
	flag.Parse()
	initDefaultConfig()

	// get all env
	MasterHostsVal = os.Getenv(MasterHostsEnv)
	ExternalLoadBalancerVal = os.Getenv(ExternalLoadBalancerEnv)
	NodeHostsVal = os.Getenv(NodeHostsEnv)
	EtcdHostsVal = os.Getenv(EtcdHostsEnv)
	OperationVal = os.Getenv(OperationEnv)
	ClusterNameVal = os.Getenv(ClusterNameEnv)
	ClusterNamespaceVal = os.Getenv(ClsuterNamespaceEnv)
	HostsYAMLVal = os.Getenv(HostsYAMLEnv)
	PrivateKeyVal = os.Getenv(PrivateKeyEnv)
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
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	go func() {
		// listening OS shutdown singal
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan

		fmt.Println("Got OS shutdown signal, shutting down server gracefully...")
		ctx.Done()
	}()

	cmdStdout := make(chan string)
	cmdError := make(chan error, 1)

	// save env to file and use ansible copy to all hosts specified dir
	success := envSaveFile()
	if !success {
		fmt.Println("save env failed")
		os.Exit(1)
	}
	fmt.Println("save env success")

	//save hosts.yaml
	success = stringSaveFile(HostsYAMLVal, HostsYAMLFileName, os.ModePerm)
	if !success {
		fmt.Println("save hosts.yaml failed")
		os.Exit(1)
	}
	fmt.Println("save hosts.yaml success")

	//save private.key
	success = stringSaveFile(PrivateKeyVal, PrivateKeyFileName, os.FileMode(0600))
	if !success {
		fmt.Println("save private.key failed")
		os.Exit(1)
	}
	fmt.Println("save private.key success")

	switch OperationVal {
	case enum.KubeCreating:
		go func() { execKubeCreatingCmds(cmdStdout, cmdError) }()
	case enum.KubeScalingUp:
		//go func() {
		//deployEtcdCmd := exec.Command("/bin/bash", "-c", `df -lh`)
		//s := execCmd(deployEtcdCmd, cmdError)
		//fmt.Println("case:", s)
		//cmdStdout <- "asda"
		//fmt.Println("cmdStdout <- s ")
		//}()
		go func() { execKubeScalingUpCmds(cmdStdout, cmdError) }()
	case enum.KubeScalingDown:
		go func() { execKubeScalingDownCmds(cmdStdout, cmdError) }()
	case enum.KubeTerminating:
		go func() { execKubeTerminatingCmds(cmdStdout, cmdError) }()

	// do not know the callback path, exit
	default:
		fmt.Println("the OPERATION env not found,exit")
		os.Exit(1)
	}

	// default timeout is ten minutes
	timeout := viper.GetInt(timeout)
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		wrapErr := errors.Errorf("the operation is timeout(%v)", timeout)
		fmt.Println(wrapErr)
		callback("", wrapErr)
	case err := <-cmdError:
		fmt.Println(err)
		callback("", err)
	case stdout := <-cmdStdout:
		callback(stdout, nil)
	case <-ctx.Done():
		os.Exit(0)
	}
}

func callback(stdout string, err error) {
	fmt.Println("-------callback")

	resp := types.Callback{
		Name:       ClusterNameVal,
		Namespace:  ClusterNamespaceVal,
		Region:     viper.GetString(region),
		KubeConfig: "",
		Success:    true,
		Message:    stdout,
	}

	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
	}

	var path string
	switch OperationVal {
	case enum.KubeCreating:
		if err == nil {
			// get kubeconfig if deploy success
			cmdError := make(chan error, 1)
			kubeconfig := execGetKubeconfig(cmdError)
			if kubeconfig == "" {
				os.Exit(1)
			}
			resp.KubeConfig = kubeconfig
		}
		path = viper.GetString(creatingCallbackPath)
	case enum.KubeScalingUp:
		path = viper.GetString(scalingUpCallbackPath)
	case enum.KubeScalingDown:
		path = viper.GetString(scalingDownCallbackPath)
	case enum.KubeTerminating:
		path = viper.GetString(terminatingCallbackPath)
	}

	packPath := packURLPath(path, map[string]string{"region": "", "name": ClusterNameVal})
	sendRequest(packPath)
}

// sendRequest is send request to controller
func sendRequest(path string) {
	fmt.Println("----------sendRequest")

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
		fmt.Printf("new restclient failed with:%v \n", err)
		return
	}

	resp, err := c.Post().
		Body().
		Do().
		Raw()

	if err != nil {
		fmt.Printf("response failed with:%v \n", err)
		return
	}

	fmt.Printf("response result is:%v \n", string(resp))
}

func execKubeCreatingCmds(cmdStdout chan<- string, cmdError chan<- error) {
	deployEtcdCmd := exec.Command("sh", "-c", DeployEtcdCmd)
	execCmd(deployEtcdCmd, cmdError)
	if len(cmdError) != 0 {
		fmt.Println("exec deployEtcdCmd failed")
		return
	}
	fmt.Println("exec deployEtcdCmd success")

	deployMasterCmd := exec.Command("sh", "-c", DeployMasterCmd)
	execCmd(deployMasterCmd, cmdError)
	if len(cmdError) != 0 {
		fmt.Println("exec deployMasterCmd failed")
		return
	}
	fmt.Println("exec deployMasterCmd success")

	deployNodeCmd := exec.Command("sh", "-c", DeployNodeCmd)
	stdout := execCmd(deployNodeCmd, cmdError)
	if len(cmdError) != 0 {
		fmt.Println("exec deployNodeCmd failed")
		return
	}
	fmt.Println("exec deployNodeCmd success")
	// the job is finished
	cmdStdout <- stdout
}

func execKubeScalingUpCmds(cmdStdout chan<- string, cmdError chan<- error) {
	scaleUpCmd := exec.Command("sh", "-c", ScaleUpCmd)
	stdout := execCmd(scaleUpCmd, cmdError)
	cmdStdout <- stdout
}

func execKubeScalingDownCmds(cmdStdout chan<- string, cmdError chan<- error) {
	scaleDownCmd := exec.Command("sh", "-c", ScaleDownCmd)
	stdout := execCmd(scaleDownCmd, cmdError)
	cmdStdout <- stdout
}

func execKubeTerminatingCmds(cmdStdout chan<- string, cmdError chan<- error) {
	terminatingCmd := exec.Command("sh", "-c", TerminatingCmd)
	stdout := execCmd(terminatingCmd, cmdError)
	cmdStdout <- stdout
}

func execGetKubeconfig(cmdError chan<- error) string {
	fmt.Println("---------start get kubeconfig")

	var kubeconfig string
	for _, master := range strings.Split(MasterHostsVal, ",") {
		packKubeconfigCmd := fmt.Sprintf(KubeconfigCmd, master)
		getKubeconfigCmd := exec.Command("sh", "-c", packKubeconfigCmd)
		kubeconfig = execCmd(getKubeconfigCmd, cmdError)
		if kubeconfig != "" {
			break
		}
	}
	return kubeconfig
}

func execCmd(cmd *exec.Cmd, cmdError chan<- error) string {
	// create command pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		wrapErr := errors.Errorf("obtain %v stdout pipe for command failed with:%v", cmd.Args, err)
		cmdError <- wrapErr
		return ""
	}

	// exec command
	if err = cmd.Start(); err != nil {
		wrapErr := errors.Errorf("command %v start failed with:%v", cmd.Args, err)
		cmdError <- wrapErr
		return ""
	}

	// read all stdout
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		wrapErr := errors.Errorf("read %v stdout failed with:%v", cmd.Args, err)
		cmdError <- wrapErr
		return ""
	}

	// wait cmd exec finished
	if err = cmd.Wait(); err != nil {
		wrapErr := errors.Errorf("wait cmd %v exec finished failed with:%v", cmd.Args, err)
		cmdError <- wrapErr
		return ""
	}

	// print logs in job's pod
	fmt.Println(string(bytes))
	return string(bytes)
}

func stringSaveFile(env string, fileName string, perm os.FileMode) bool {
	// create hosts.yaml or overrite it
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, perm)
	if err != nil {
		fmt.Printf("open file %v failed with:%v \n", fileName, err)
		return false
	}
	defer f.Close()
	defer f.Sync()

	_, err = f.WriteString(env)
	if err != nil {
		fmt.Printf("write file %v failed with:%v \n", fileName, err)
		return false
	}
	return true
}

func envSaveFile() bool {
	f, err := os.OpenFile(EnvFileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("open file %v failed with:%v \n", EnvFileName, err)
		return false
	}

	defer f.Close()
	defer f.Sync()

	envMaps := map[string]string{
		MasterHostsEnv:          MasterHostsVal,
		ExternalLoadBalancerEnv: ExternalLoadBalancerVal,
		NodeHostsEnv:            NodeHostsVal,
		EtcdHostsEnv:            EtcdHostsVal,
		OperationEnv:            OperationVal,
		ClusterNameEnv:          ClusterNameVal,
		ClsuterNamespaceEnv:     ClusterNamespaceVal,
	}

	var envsStr string
	for k, v := range envMaps {
		envsStr += fmt.Sprintf("%v=\"%v\" \n", k, v)
	}

	_, err = f.WriteString(envsStr)
	if err != nil {
		fmt.Printf("write file %v failed with:%v \n", EnvFileName, err)
		return false
	}
	return true
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
