package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gosoon/kubernetes-operator/cmd/installer/kube-on-kube/server/app"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	command := app.NewServerCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
