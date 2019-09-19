package docker

import (
	"os/exec"
	"time"

	"github.com/gosoon/glog"
)

// PullIfNotPresent will pull an image if it is not present locally
// retrying up to retries times
// it returns true if it attempted to pull, and any errors from pulling
func PullIfNotPresent(image string, retries int) (bool, error) {
	cmd := exec.Command("docker", "inspect", image)
	if err := cmd.Run(); err != nil {
		glog.Infof("Image: %s present locally", image)
		return false, nil
	}
	return pull(image, retries)

}

// Pull pulls an image, retrying up to retries times
func pull(image string, retries int) (bool, error) {
	glog.Infof("Pulling image: %s ...", image)
	err := exec.Command("docker", "pull", image).Run()
	// retry pulling up to retries times if necessary
	if err != nil {
		for i := 0; i < retries; i++ {
			time.Sleep(time.Second * time.Duration(i))
			glog.Infof("Trying again to pull image: %s ...", image)
			err = exec.Command("docker", "pull", image).Run()
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		glog.Infof("Failed to pull image: %s", image)
		return false, err
	}
	return true, nil
}
