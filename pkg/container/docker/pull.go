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
