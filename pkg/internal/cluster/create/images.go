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

package create

import (
	"fmt"
	"strings"

	"github.com/gosoon/kubernetes-operator/pkg/container/docker"
	"github.com/gosoon/kubernetes-operator/pkg/internal/util/cli"
)

// ensureNodeImages ensures that the node images used by the create
// configuration are present
func ensureNodeImages(status *cli.Status, image string) {
	// prints user friendly message
	if strings.Contains(image, "@sha256:") {
		image = strings.Split(image, "@sha256:")[0]
	}
	status.Start(fmt.Sprintf("Ensuring node image (%s) ", image))

	fmt.Println("pull image ...")
	// attempt to explicitly pull the image if it doesn't exist locally
	// we don't care if this errors, we'll still try to run which also pulls
	_, _ = docker.PullIfNotPresent(image, 3)
}
