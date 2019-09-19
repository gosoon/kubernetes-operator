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

package sshserver

import (
	"testing"
	"time"

	"github.com/gosoon/kubernetes-operator/pkg/types"

	"github.com/stretchr/testify/assert"
)

var CheckCmd = []string{"chmod +x /tmp/check.sh", "/tmp/check.sh"}

func TestCopyFile(t *testing.T) {
	testCases := []types.SSHInfo{
		{
			IP:       "192.168.75.178",
			Username: "root",
			Port:     22,
			CmdList:  CheckCmd,
			Key:      "",
			Timeout:  35 * time.Second,
		},
		{
			IP:       "192.168.75.178",
			Username: "root",
			Port:     22,
			CmdList:  CheckCmd,
			Key:      "asdas",
			Timeout:  35 * time.Second,
		},
		{
			IP:       "",
			Username: "root",
			Port:     22,
			CmdList:  CheckCmd,
			Key:      "asdas",
			Timeout:  35 * time.Second,
		},
	}

	for _, test := range testCases {
		t.Log(test)
		sshServer, err := NewSSHServer(&test)
		if err != nil {
			t.Log(err)
		}
		sshServer.Dossh(ch)
		if !assert.Equal(t, nil, err) {
			t.Fatalf("expected: %v but get %v", nil, err)
		}
	}
}
