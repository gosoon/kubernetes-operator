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

package docker

import (
	"github.com/gosoon/glog"
	"github.com/gosoon/kubernetes-operator/pkg/exec"
)

// RunOpt is an option for Run
type RunOpt func(*runOpts) *runOpts

// actual options struct
type runOpts struct {
	RunArgs       []string
	ContainerArgs []string
}

// WithRunArgs sets the args for docker run
// as in the args portion of `docker run args... image containerArgs...`
func WithRunArgs(args ...string) RunOpt {
	return func(r *runOpts) *runOpts {
		r.RunArgs = args
		return r
	}
}

// WithContainerArgs sets the args to the container
// as in the containerArgs portion of `docker run args... image containerArgs...`
// NOTE: this is only the args portion before the image
func WithContainerArgs(args ...string) RunOpt {
	return func(r *runOpts) *runOpts {
		r.ContainerArgs = args
		return r
	}
}

// Run creates a container with "docker run", with some error handling
func Run(image string, opts ...RunOpt) error {
	o := &runOpts{}
	for _, opt := range opts {
		o = opt(o)
	}

	// construct the actual docker run argv
	args := []string{"run"}
	args = append(args, o.RunArgs...)
	args = append(args, image)
	cmd := exec.Command("docker", args...)

	output, err := exec.CombinedOutputLines(cmd)
	if err != nil {
		// log error output if there was any
		for _, line := range output {
			glog.Error(line)
		}
		return err
	}
	return nil
}
