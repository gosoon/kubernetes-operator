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

package protobuf

import (
	"reflect"
	"testing"

	installerv1 "github.com/gosoon/kubernetes-operator/pkg/apis/installer/v1"
)

func TestConvertClusterTypes(t *testing.T) {
	testCases := []*installerv1.KubernetesClusterRequest{
		{
			TypeMeta: installerv1.TypeMeta{
				Kind:       "KubernetesCluster",
				APIVersion: "ecs.yun.com/v1",
			},
			ObjectMeta: installerv1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
			},
			Spec: installerv1.KubernetesClusterSpec{
				Cluster: installerv1.Cluster{
					TimeoutMins:          "",
					ClusterType:          "kubernetes",
					PodCIDR:              "192.168.0.0/24",
					ServiceCIDR:          "172.10.0.0/24",
					MasterList:           []installerv1.Node{installerv1.Node{IP: "192.168.1.100"}},
					ExternalLoadBalancer: "192.168.1.100",
					NodeList:             []installerv1.Node{installerv1.Node{IP: "192.168.1.100"}},
					EtcdList:             []installerv1.Node{installerv1.Node{IP: "192.168.1.100"}},
					Region:               "default",
					KubeVersion:          "v1.15.3",
					AuthConfig: installerv1.AuthConfig{
						Username:      "root",
						Password:      "123",
						PrivateSSHKey: "asdasdas",
					},
				},
				Addons: installerv1.Addons{},
			},
		},
		{
			Spec: installerv1.KubernetesClusterSpec{
				Cluster: installerv1.Cluster{
					TimeoutMins:          "",
					ClusterType:          "kubernetes",
					PodCIDR:              "192.168.0.0/24",
					ServiceCIDR:          "172.10.0.0/24",
					MasterList:           []installerv1.Node{installerv1.Node{IP: "192.168.1.100"}},
					ExternalLoadBalancer: "192.168.1.100",
					NodeList:             []installerv1.Node{installerv1.Node{IP: "192.168.1.100"}},
					EtcdList:             []installerv1.Node{installerv1.Node{IP: "192.168.1.100"}},
					Region:               "default",
					KubeVersion:          "v1.15.3",
					AuthConfig: installerv1.AuthConfig{
						Username:      "root",
						Password:      "123",
						PrivateSSHKey: "asdasdas",
					},
				},
				Addons: installerv1.Addons{},
			},
		},
	}

	for _, test := range testCases {
		cluster, err := ClusterConvertToTypes(test)
		if err != nil {
			t.Fatalf("except: nil but get err:%v ", err)
		}

		clusterRequest, err := ClusterConvertToProtobuf(cluster)
		if err != nil {
			t.Fatalf("except: nil but get err:%v ", err)
		}

		if !reflect.DeepEqual(test, clusterRequest) {
			t.Fatalf("except: is equal but get not equal ")
		}
	}
}
