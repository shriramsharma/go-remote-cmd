package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSHAgent Authenticate using ssh private key.
// Reads the private key cert from the ssh agent of the operating system
func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func executeCommand(ip string, command string, sshConfig *ssh.ClientConfig) {
	host := fmt.Sprintf("%s:%s", ip, "22")
	connection, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		log.Fatal(err)
	}

	session, err := connection.NewSession()
	fmt.Println("New session")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(os.Stdout, stdout)

	if err := session.Run(command); err != nil {
		log.Fatal(err)
	}

}

func main() {

	ipsFile := os.Args[1]
	command := os.Args[2]

	sshConfig := &ssh.ClientConfig{
		User: "app",
		Auth: []ssh.AuthMethod{SSHAgent()},
	}

	file, err := os.Open(ipsFile)
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			ip := scanner.Text()
			fmt.Println(ip)
			go func(ip string) {
				executeCommand(ip, command, sshConfig)
			}(ip)
		}
	}
	defer file.Close()

	for {
		time.Sleep(time.Second)
	}

}
