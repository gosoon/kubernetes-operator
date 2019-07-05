package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type certKeyConfig struct {
	Algo string `json:"algo"`
	Size int    `json:"size"`
}
type nameConfig struct {
	C  string `json:"C"`
	ST string `json:"ST"`
	L  string `json:"L"`
	O  string `json:"O"`
	OU string `json:"OU"`
}

type certConfig struct {
	CN    string        `json:"CN"`
	Hosts []string      `json:"hosts"`
	Key   certKeyConfig `json:"key"`
	Names []nameConfig  `json:"names"`
}

const (
	Etcd       string = "etcd"
	Kubernetes string = "kubernetes"

	EtcdListENV   string = "CLUSTER_ETCD_LIST"
	MasterListENV string = "CLUSTER_MASTER_LIST"
)

var component string
var csrFileName string

func init() {
	flag.StringVar(&component, "component", "", "generate kubernetes|etcd cert")
	flag.StringVar(&csrFileName, "csrfile", "", "cert csr file name")
}

func main() {
	flag.Parse()
	if len(os.Args) < 3 {
		log.Fatal("program needs component and csrfile")
	}

	switch component {
	case Etcd:
		generateCerts(EtcdListENV)
	case Kubernetes:
		generateCerts(MasterListENV)
	default:
		fmt.Println("please set component, 'etcd' or 'kubernetes'")
	}
}

func generateCerts(env string) {
	var nodeList []string
	nodeListEnv := os.Getenv(env)
	if len(nodeListEnv) != 0 {
		hosts := strings.Split(nodeListEnv, " ")

		for i := range hosts {
			if len(hosts[i]) != 0 {
				nodeList = append(nodeList, hosts[i])
			}
		}
	}
	configFile, err := ioutil.ReadFile(csrFileName)
	if err != nil {
		log.Fatalf("open config file failed: %s, %s", csrFileName, err)
	}

	certCfg := &certConfig{}
	err = json.Unmarshal(configFile, certCfg)
	if err != nil {
		log.Fatalf("parse cert config failed: %s, %s", csrFileName, err)
	}

	if len(nodeList) > 0 {
		certCfg.Hosts = append(certCfg.Hosts, nodeList...)
	}

	newCertConfig, err := json.Marshal(certCfg)
	fmt.Print(string(newCertConfig))
}
