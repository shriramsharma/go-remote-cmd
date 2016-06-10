package main

import (
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func init() {

}

// SSHAgent Authenticate using ssh private key.
// Reads the private key cert from the ssh agent of the operating system
func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

// SSHConfig Confugure SSH
func SSHConfig(user string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: "app",
		Auth: []ssh.AuthMethod{SSHAgent()},
	}
}
