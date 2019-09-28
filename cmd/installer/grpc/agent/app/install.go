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

package app

import (
	"context"
	"fmt"
	"io"
	"os"

	installerv1 "github.com/gosoon/kubernetes-operator/pkg/apis/installer/v1"

	"google.golang.org/grpc"
)

func NewAgentClient(ip string, port string) (installerv1.InstallerClient, error) {
	server := ip + ":" + port
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := installerv1.NewInstallerClient(conn)
	//CopyFrom(client)
	return client, nil
}

func CopyFrom(client installerv1.InstallerClient, fileName string) error {
	stream, err := client.CopyFile(context.Background(), &installerv1.File{Name: "main.go"})
	if err != nil {
		return err
	}
	defer stream.CloseSend()

	destFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	for {
		file, err := stream.Recv()
		// receiver finished
		if err == io.EOF {
			//close(waitc)
			return nil
		}
		if err != nil {
			return err
		}
		_, err = destFile.Write(file.Content)
		if err != nil {
			return err
		}
	}

	return nil
}
