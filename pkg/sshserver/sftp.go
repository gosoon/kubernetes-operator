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
	"os"

	"github.com/gosoon/glog"
)

func (s *sshServer) CopyFile(localFilePath string, remoteFilePath string) error {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		glog.Errorf("scp to %v err: %v", s.IP, err)
		return err
	}
	defer srcFile.Close()

	dstFile, err := s.SftpClient.Create(remoteFilePath)
	if err != nil {
		glog.Errorf("scp to %v err: %v", s.IP, err)
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf[0:n])
	}
	return nil
}
