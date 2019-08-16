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
	"bytes"
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gosoon/kubernetes-operator/pkg/types"
	"github.com/gosoon/kubernetes-operator/pkg/utils"

	"github.com/gosoon/glog"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// sshServer xxx
type sshServer struct {
	IP           string
	Port         int
	Username     string
	Password     string
	CmdFile      string
	Cmds         string
	CmdList      []string
	Key          string
	Timeout      time.Duration
	CipherList   []string
	ClientConfig *ssh.ClientConfig
	Client       *ssh.Client
	Session      *ssh.Session
	SftpClient   *sftp.Client
}

// NewSSHServer xxx
func NewSSHServer(host *types.SSHInfo) (Interface, error) {
	var (
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		err          error
	)
	// get auth method
	var auth []ssh.AuthMethod

	if host.Key == "" {
		auth = append(auth, ssh.Password(host.Password))
	} else {
		pemBytes := []byte(host.Key)
		if valided, key := utils.ValidBase64Str(host.Key); valided {
			pemBytes = []byte(key)
		}
		var signer ssh.Signer
		if host.Password == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(host.Password))
		}
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	if len(host.CipherList) == 0 {
		config = ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com",
				"arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
		}
	} else {
		config = ssh.Config{
			Ciphers: host.CipherList,
		}
	}

	// create client config
	clientConfig = &ssh.ClientConfig{
		User:    host.Username,
		Auth:    auth,
		Timeout: host.Timeout,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	// connet to ssh server
	addr := fmt.Sprintf("%s:%d", host.IP, host.Port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		glog.Errorf("unable to connect: %v", err)
		return nil, err
	}

	// create sftp client
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return nil, err
	}

	// create session
	session, err := client.NewSession()
	if err != nil {
		glog.Errorf("unable to create session: %v", err)
		return nil, err
	}

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		glog.Errorf("request for pseudo terminal failed: %v", err)
		session.Close()
		return nil, err
	}

	return &sshServer{
		IP:           host.IP,
		Port:         host.Port,
		Username:     host.Username,
		Password:     host.Password,
		CmdFile:      host.CmdFile,
		Cmds:         host.Cmds,
		CmdList:      host.CmdList,
		Key:          host.Key,
		Timeout:      host.Timeout,
		ClientConfig: clientConfig,
		Client:       client,
		Session:      session,
		SftpClient:   sftpClient}, nil
}

func (s *sshServer) Dossh(ch chan<- types.PrecheckResult) {
	result := make(chan types.PrecheckResult)
	go s.DoSSHSessionRun(result)

	ctx, cacnel := context.WithTimeout(context.Background(), s.Timeout)
	defer cacnel()

	select {
	case <-ctx.Done():
		precheckTimeout := types.PrecheckResult{
			Host:    s.IP,
			Success: false,
			Result:  fmt.Sprintf("sshServer exec check list timeout (>%v second)", strconv.Itoa(int(s.Timeout))),
		}
		ch <- precheckTimeout
	case res := <-result:
		ch <- res
	}
	return
}

func (s *sshServer) DoSSHSessionRun(ch chan<- types.PrecheckResult) {
	defer s.Client.Close()
	defer s.Session.Close()

	result := types.PrecheckResult{
		Host:    s.IP,
		Success: true,
		CmdList: s.CmdList,
	}

	s.CmdList = append(s.CmdList, "exit")
	newcmd := strings.Join(s.CmdList, "&&")

	var outbt, errbt bytes.Buffer
	s.Session.Stdout = &outbt
	s.Session.Stderr = &errbt

	// run cmd
	err := s.Session.Run(newcmd)
	if err != nil {
		result.Success = false
		result.Result = outbt.String() + fmt.Sprintf("\n<%v>", err.Error())
		ch <- result
		return
	}

	// get output
	result.Result = outbt.String()
	if len(errbt.String()) != 0 {
		result.Success = false
		result.Result = outbt.String() + "\n" + errbt.String()
	} else {
		result.Success = true
	}
	ch <- result
}
