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
