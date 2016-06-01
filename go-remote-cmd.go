package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type allSessions struct {
	session *ssh.Session
	host    string
}

// SSHAgent Authenticate using ssh private key.
// Reads the private key cert from the ssh agent of the operating system
func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func executeCommand(ip string, command string, sshConfig *ssh.ClientConfig, s *[]allSessions) {
	host := fmt.Sprintf("%s:%s", ip, "22")
	connection, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		log.Fatal(err)
	}

	session, err := connection.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	*s = append(*s, allSessions{session, ip})

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)

	go func() {
		for scanner.Scan() {
			fmt.Printf("%s %s\n", "\x1b[36m"+ip+"\x1b[0m", scanner.Text())
		}
	}()

	if err := session.Run(command); err != nil {
		log.Fatal(err)
	}

}

// HandleControlCGracefully This function would handle ctrl-C gracefully by closing all the remote sessions.
func HandleControlCGracefully(sessions *[]allSessions) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	<-sigs
	for _, session := range *sessions {
		fmt.Printf("Closing session on host: %s\n", session.host)
		session.session.Close()
	}
	os.Exit(0)
}

func main() {

	ipsFile := os.Args[1]
	command := os.Args[2]

	var wg sync.WaitGroup
	var sessions []allSessions

	sshConfig := &ssh.ClientConfig{
		User: "app",
		Auth: []ssh.AuthMethod{SSHAgent()},
	}

	i := 0

	file, err := os.Open(ipsFile)
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			i++
			wg.Add(i)
			ip := scanner.Text()
			go func(ip string) {
				defer wg.Done()
				executeCommand(ip, command, sshConfig, &sessions)
			}(ip)
		}
	}
	defer file.Close()

	go HandleControlCGracefully(&sessions)

	wg.Wait()

}
