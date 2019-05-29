package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
)

type ViaSSHDialer struct {
	client *ssh.Client
}

func (self *ViaSSHDialer) Dial(addr string) (net.Conn, error) {
	return self.client.Dial("tcp", addr)
}

func main() {

	sshHost := "example.com"   // SSH Server Hostname/IP
	sshPort := 22              // SSH Port
	sshUser := "ssh-user"      // SSH Username
	sshPass := "ssh-pass"      // Empty string for no password
	dbUser := "dbuser"         // DB username
	dbPass := "dbpass"         // DB Password
	dbHost := "localhost:3306" // DB Hostname/IP
	dbName := "database"       // Database name

	var agentClient agent.Agent
	// Establish a connection to the local ssh-agent
	if conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		defer conn.Close()

		// Create a new instance of the ssh agent
		agentClient = agent.NewClient(conn)
	}

	// The client configuration with configuration option to use the ssh-agent
	sshConfig := &ssh.ClientConfig{
		User:            sshUser,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// When the agentClient connection succeeded, add them as AuthMethod
	if agentClient != nil {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeysCallback(agentClient.Signers))
	}
	// When there's a non empty password add the password AuthMethod
	if sshPass != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PasswordCallback(func() (string, error) {
			return sshPass, nil
		}))
	}

	// Connect to the SSH Server
	sshcon, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshHost, sshPort), sshConfig)
	if err == nil {
		defer sshcon.Close()

		// Now we register the ViaSSHDialer with the ssh connection as a parameter
		mysql.RegisterDial("mysql+tcp", (&ViaSSHDialer{sshcon}).Dial)

		// And now we can use our new driver with the regular mysql connection string tunneled through the SSH connection
		if db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@mysql+tcp(%s)/%s", dbUser, dbPass, dbHost, dbName)); err == nil {

			fmt.Printf("Successfully connected to the db\n")

			if rows, err := db.Query("select sys_user_id,mobile from sys_user"); err == nil {
				for rows.Next() {
					var sys_user_id int64
					var mobile string
					rows.Scan(&sys_user_id, &mobile)
					fmt.Printf("sys_user_id: %d  mobile: %s\n", sys_user_id, mobile)
				}
				rows.Close()
			} else {
				fmt.Printf("Failure: %s", err.Error())
			}

			db.Close()

		} else {
			fmt.Printf("Failed to connect to the db: %s\n", err.Error())
		}
	} else {
		fmt.Println(err.Error())
	}
}
