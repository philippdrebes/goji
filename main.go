package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/akamensky/argparse"
	"github.com/andygrunwald/go-jira"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	fmt.Println("Hello Goji!")

	parser := argparse.NewParser("print", "Prints provided string to stdout")
	user := parser.String("u", "user", &argparse.Options{Required: false, Help: "Username"})
	err := parser.Parse(os.Args)

	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return
	}

	username, password := getCredentials(*user)

	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	jiraClient, _ := jira.NewClient(tp.Client(), "https://servicedesk.softec.ch")

	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	u, _, err := jiraClient.User.Get(username)

	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	fmt.Printf("\n\nLogged in as %v\n", u.EmailAddress)
}

func getCredentials(user string) (string, string) {
	var username string
	var password string

	r := bufio.NewReader(os.Stdin)

	//fmt.Print("Jira URL: ")
	//jiraURL, _ := r.ReadString('\n')

	if len(user) == 0 {
		fmt.Print("Jira Username: ")
		username, _ = r.ReadString('\n')
	}

	fmt.Print("Jira Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password = string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password)
}
