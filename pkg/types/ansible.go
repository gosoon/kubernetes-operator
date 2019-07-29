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

package types

// HostsYamlFormat xxx
type HostsYamlFormat struct {
	All AllHosts `json:"all"`
}

// AllHosts xxx
type AllHosts struct {
	Hosts    map[string]*Host `json:"hosts"`
	Children Children         `json:"children"`
}

// Children xxx
type Children struct {
	KubeMaster map[string]map[string]*Host `json:"kube-master"`
	KubeNode   map[string]map[string]*Host `json:"kube-node"`
	Etcd       map[string]map[string]*Host `json:"etcd"`
	Calico     map[string]map[string]*Host `json:"calico-rr"`
}

// Host xxx
type Host struct {
	AnsibleHost string `json:"ansible_host,omitempty"`
	IP          string `json:"ip,omitempty"`
	AccessIP    string `json:"access_ip,omitempty"`
}
