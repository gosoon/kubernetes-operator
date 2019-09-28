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

package kubeadmjoin

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cloudflare/cfssl/log"
	"github.com/gosoon/glog"
	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	installerv1 "github.com/gosoon/kubernetes-operator/pkg/apis/installer/v1"
	"github.com/gosoon/kubernetes-operator/pkg/exec"
	"github.com/gosoon/kubernetes-operator/pkg/fs"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/constants"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/nodes"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/actions"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Action implements action for creating the kubeadm join
// and deployng it on the bootrap control-plane node.
type Action struct{}

// NewAction returns a new action for creating the kubeadm jion
func NewAction() actions.Action {
	return &Action{}
}

// Execute runs the action
func (a *Action) Execute(ctx *actions.ActionContext) error {
	if ctx.Cluster.Role == ecsv1.SecondaryControlPlaneRole {
		//allNodes := ctx.Cluster.Config.Nodes
		if err := joinSecondaryControlPlanes(ctx); err != nil {
			return err
		}
	}

	// then join worker nodes if any
	if ctx.Cluster.Role == ecsv1.WorkerRole {
		if err := joinWorkers(ctx); err != nil {
			return err
		}
	}

	return nil
}

func joinSecondaryControlPlanes(ctx *actions.ActionContext) error {
	ctx.Status.Start("Joining more control-plane nodes")
	defer ctx.Status.End(false)

	if err := runKubeadmJoinControlPlane(ctx); err != nil {
		return err
	}

	ctx.Status.End(true)
	return nil
}

func joinWorkers(ctx *actions.ActionContext) error {
	ctx.Status.Start("Joining worker nodes ")
	defer ctx.Status.End(false)

	if err := runKubeadmJoin(); err != nil {
		glog.Error(err)
		return err
	}

	ctx.Status.End(true)
	return nil
}

// runKubeadmJoinControlPlane executes kubadm join --control-plane command
func runKubeadmJoinControlPlane(ctx *actions.ActionContext) error {
	// creates the folder tree for pre-loading necessary cluster certificates
	// on the joining node
	if err := exec.Command("mkdir", "-p", "/etc/kubernetes/pki/etcd").Run(); err != nil {
		return errors.Wrap(err, "failed to join node with kubeadm")
	}

	// define the list of necessary cluster certificates
	fileNames := []string{
		"ca.crt", "ca.key",
		"front-proxy-ca.crt", "front-proxy-ca.key",
		"sa.pub", "sa.key",
		// TODO(someone): if we gain external etcd support these will be
		// handled differently
		"etcd/ca.crt", "etcd/ca.key",
	}

	allNodes := ctx.Cluster.Config.Nodes
	// get the handle for the bootstrap control plane node (the source for necessary cluster certificates)
	controlPlaneNode, err := nodes.BootstrapControlPlaneNode(allNodes)
	if err != nil {
		return err
	}

	// start grpc client
	address := controlPlaneNode.IP + ":" + ctx.Port
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := installerv1.NewInstallerClient(conn)
	// copy ca files from controlPlane
	copyCAFromControlPlaneNode(client, fileNames)

	// creates a temporary folder on the host that should acts as a transit area
	// for moving necessary cluster certificates
	tmpDir, err := fs.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	err = os.MkdirAll(filepath.Join(tmpDir, "/etcd"), os.ModePerm)
	if err != nil {
		return err
	}

	// TODO(gosoon): use grpc transfer file
	// copies certificates from the bootstrap control plane node to the joining node
	for _, fileName := range fileNames {
		// sets the path of the certificate into a node
		containerPath := path.Join("/etc/kubernetes/pki", fileName)
		// set the path of the certificate into the tmp area on the host
		tmpPath := filepath.Join(tmpDir, fileName)
		// copies from bootstrap control plane node to tmp area
		if err := nodes.CopyFrom(containerPath, tmpPath); err != nil {
			return errors.Wrapf(err, "failed to copy certificate %s", fileName)
		}
		// copies from tmp area to joining node
		if err := nodes.CopyTo(tmpPath, containerPath); err != nil {
			return errors.Wrapf(err, "failed to copy certificate %s", fileName)
		}
	}
	return runKubeadmJoin()
}

// copyCAFromControlPlaneNode is copy all apiserver and etcd ca from controlPlaneNode to locally
// the source ca file in /etc/kubernetes/pki/ dir
func copyCAFromControlPlaneNode(client installerv1.InstallerClient, fileNames []string) error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	errs := make(chan string, len(fileNames))
	var wg sync.WaitGroup
	for _, name := range fileNames {
		// destFilePath is /tmp/install/pki/
		destFilePath := constants.InstallPath + "pki/" + name
		// open file use append mode
		destFile, err := os.OpenFile(destFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, acquireFileMode(name))
		if err != nil {
			return err
		}
		defer destFile.Close()

		srcFile := "/etc/kubernetes/pki/" + name
		wg.Add(1)

		// if any file sync failed and write err to errs channel
		go func(srcFile string, destFile *os.File, errc chan<- string) {
			defer wg.Done()
			stream, err := client.CopyFile(ctx, &installerv1.File{Name: srcFile})
			if err != nil {
				errc <- err.Error()
				return
			}
			defer stream.CloseSend()
			for {
				fileStream, err := stream.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					errc <- err.Error()
					return
				}
				_, err = destFile.Write(fileStream.Content)
				if err != nil {
					errc <- err.Error()
					return
				}
			}
		}(srcFile, destFile, errs)
	}
	wg.Done()
	close(errs)

	if len(errs) != 0 {
		errStrs := []string{}
		for err := range errs {
			errStrs = append(errStrs, err)
		}
		return errors.New(strings.Join(errStrs, "\n"))
	}
	return nil
}

// runKubeadmJoin executes kubadm join command
func runKubeadmJoin() error {
	// run kubeadm join
	// TODO(bentheelder): this should be using the config file
	cmd := exec.Command(
		"kubeadm", "join",
		// the join command uses the config file generated in a well known location
		"--config", "/tmp/install/kubeadm.conf",
		// preflight errors are expected, in particular for swap being enabled
		// TODO(bentheelder): limit the set of acceptable errors
		"--ignore-preflight-errors=all",
		// increase verbosity for debugging
		"--v=6",
	)
	lines, err := exec.CombinedOutputLines(cmd)
	log.Debug(strings.Join(lines, "\n"))
	if err != nil {
		return errors.Wrap(err, "failed to join node with kubeadm")
	}
	return nil
}

// key file mode is 0600 only read
// crt file mode is 0644
func acquireFileMode(fileName string) os.FileMode {
	fileMode := os.FileMode(0644)
	if strings.Contains(fileName, "key") {
		fileMode = os.FileMode(0600)
	}
	return fileMode
}
