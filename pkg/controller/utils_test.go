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

package controller

import (
	"reflect"
	"testing"

	ecsv1 "github.com/gosoon/kubernetes-operator/pkg/apis/ecs/v1"
	"github.com/gosoon/kubernetes-operator/pkg/enum"
	"github.com/stretchr/testify/assert"
)

func TestConvertNodesToString(t *testing.T) {
	testCases := []struct {
		nodeList []ecsv1.Node
		expect   string
	}{
		{
			nodeList: []ecsv1.Node{
				{IP: "192.168.1.2"},
				{IP: "192.168.1.3"},
			},
			expect: "192.168.1.2,192.168.1.3",
		},
		{
			nodeList: []ecsv1.Node{
				{IP: "192.168.1.2"},
			},
			expect: "192.168.1.2",
		},
		{
			nodeList: []ecsv1.Node{},
			expect:   "",
		},
	}

	for _, test := range testCases {
		res := convertNodesToString(test.nodeList)
		if !reflect.DeepEqual(res, test.expect) {
			t.Fatalf("expected: %v but get %v", test.expect, res)
		}
	}
}

func TestDiffNodeList(t *testing.T) {
	testCases := []struct {
		oldNodeList []ecsv1.Node
		newNodeList []ecsv1.Node
		operation   string
		expect      []ecsv1.Node
	}{
		// KubeScalingUp
		{
			oldNodeList: []ecsv1.Node{
				{IP: "192.168.1.2"},
				{IP: "192.168.1.3"},
			},
			newNodeList: []ecsv1.Node{
				{IP: "192.168.1.2"},
				{IP: "192.168.1.3"},
			},
			operation: enum.KubeScalingUp,
			//expect:    []ecsv1.Node(nil),
		},
		{
			oldNodeList: []ecsv1.Node{
				{IP: "192.168.1.2"},
			},
			newNodeList: []ecsv1.Node{
				{IP: "192.168.1.2"},
				{IP: "192.168.1.3"},
			},
			operation: enum.KubeScalingUp,
			expect: []ecsv1.Node{
				{"192.168.1.3"},
			},
		},
		{
			oldNodeList: []ecsv1.Node{
				{IP: "192.168.1.2"},
				{IP: "192.168.1.3"},
			},
			newNodeList: []ecsv1.Node{
				{IP: "192.168.1.2"},
			},
			operation: enum.KubeScalingUp,
			//expect:    []ecsv1.Node(nil),
		},
		{
			oldNodeList: []ecsv1.Node{},
			newNodeList: []ecsv1.Node{},
			operation:   enum.KubeScalingUp,
			//expect:      []ecsv1.Node(nil),
		},
		// KubeScalingDown
		{
			oldNodeList: []ecsv1.Node{
				{IP: "192.168.1.2"},
				{IP: "192.168.1.3"},
			},
			newNodeList: []ecsv1.Node{
				{IP: "192.168.1.2"},
				{IP: "192.168.1.3"},
			},
			operation: enum.KubeScalingDown,
			//expect:    []ecsv1.Node(nil),
		},
		{
			oldNodeList: []ecsv1.Node{
				{IP: "192.168.1.2"},
				{IP: "192.168.1.3"},
			},
			newNodeList: []ecsv1.Node{
				{IP: "192.168.1.2"},
			},
			operation: enum.KubeScalingDown,
			expect: []ecsv1.Node{
				{"192.168.1.3"},
			},
		},
		{
			oldNodeList: []ecsv1.Node{
				{IP: "192.168.1.2"},
			},
			newNodeList: []ecsv1.Node{
				{IP: "192.168.1.2"},
				{IP: "192.168.1.3"},
			},
			operation: enum.KubeScalingDown,
			//expect:    []ecsv1.Node(nil),
		},
		{
			oldNodeList: []ecsv1.Node{},
			newNodeList: []ecsv1.Node{},
			operation:   enum.KubeScalingDown,
			//expect:    []ecsv1.Node(nil),
		},
	}

	for _, test := range testCases {
		res := diffNodeList(test.oldNodeList, test.newNodeList, test.operation)
		if !assert.Equal(t, test.expect, res) {
			t.Fatalf("expected: %v but get %v", test.expect, res)
		}
	}
}
