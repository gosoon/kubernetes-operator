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

package precheck

import (
	"fmt"
	"sync"
	"time"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/sshserver"
	"github.com/gosoon/kubernetes-operator/pkg/types"

	"github.com/gosoon/glog"
)

const (
	sshPort       = int(22)
	sshConcurrent = int(100)
	sshUsername   = "root"

	k8sCheckScriptFile    = "./configs/check_k8s.sh"
	k3sCheckScriptFile    = "./configs/check_k3s.sh"
	remoteCheckScriptFile = "/tmp/check.sh"

	sshTimeout time.Duration = 35 * time.Second
)

var CheckCmd = []string{"chmod +x /tmp/check.sh", "/tmp/check.sh"}

// Options xxx
type Options struct {
	Concurrent            int
	Timeout               time.Duration
	LocalCheckScriptFile  string
	RemoteCheckScriptFile string
}

// precheck xxx
type precheck struct {
	opt *Options
}

// New is create a precheck client.
func New(opt *Options) Interface {
	// ssh concurrent, default 100
	if opt.Concurrent == 0 {
		opt.Concurrent = sshConcurrent
	}

	// ssh timeout
	if int64(opt.Timeout) == 0 {
		opt.Timeout = sshTimeout
	}

	// set local check.sh path
	if opt.LocalCheckScriptFile == "" {
		opt.LocalCheckScriptFile = k8sCheckScriptFile
	}

	// set remote check dir
	if opt.RemoteCheckScriptFile == "" {
		opt.RemoteCheckScriptFile = remoteCheckScriptFile
	}

	return &precheck{opt: opt}
}

func (p *precheck) HostEnv(cluster *ecsv1.KubernetesCluster, results []chan types.PrecheckResult, finish chan bool) {
	// dispatch check.sh, TODO:add other type
	if cluster.Spec.Cluster.ClusterType == ecsv1.K3sClusterType {
		p.opt.LocalCheckScriptFile = k3sCheckScriptFile
	}

	chanLimits := make(chan bool, p.opt.Concurrent)
	nodeList := cluster.Spec.Cluster.NodeList

	var wg sync.WaitGroup
	for idx, node := range nodeList {
		sshInfo := types.SSHInfo{
			IP:       node.IP,
			Username: sshUsername,
			Password: cluster.Spec.Cluster.AuthConfig.Password,
			Port:     sshPort,
			CmdList:  CheckCmd,
			Key:      cluster.Spec.Cluster.AuthConfig.PrivateSSHKey,
			Timeout:  p.opt.Timeout,
		}

		results[idx] = make(chan types.PrecheckResult, 1)
		chanLimits <- true

		wg.Add(1)
		go p.execSSHManager(&wg, chanLimits, results[idx], &sshInfo)

	}
	wg.Wait()

	// when precheck scripts exec finished and write true to finish channel
	finish <- true
}

func (p *precheck) execSSHManager(wg *sync.WaitGroup, chanLimits <-chan bool, ch chan<- types.PrecheckResult,
	sshInfo *types.SSHInfo) {

	defer wg.Done()

	// new sshServer and exec check script
	sshServer, err := sshserver.NewSSHServer(sshInfo)
	if err != nil {
		glog.Errorf("new ssh connect failed with:%v", err)

		result := types.PrecheckResult{
			Host:    sshInfo.IP,
			CmdList: sshInfo.CmdList,
			Success: false,
			Result:  fmt.Sprintf("<%v>", err),
		}
		ch <- result
		return
	}

	err = sshServer.CopyFile(p.opt.LocalCheckScriptFile, p.opt.RemoteCheckScriptFile)
	if err != nil {
		glog.Errorf("scp check script failed with:%v", err)
	}

	// exec check script
	sshServer.Dossh(ch)
	<-chanLimits
}
