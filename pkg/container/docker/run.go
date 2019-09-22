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
	args = append(args, image)
	args = append(args, o.ContainerArgs...)
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
