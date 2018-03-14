package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/akamensky/argparse"
	"golang.org/x/crypto/ssh/terminal"
	"./src"
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

	client, err := login(user)

	if client == nil || err != nil {
		fmt.Printf("\nError while trying to log in.\n%v\n", err)
		return
	}

	issues, err := client.GetAssignedTasks(*user)
	if err != nil {
		fmt.Printf("\nError while trying to get assigned issues.\n%v\n", err)
		return
	}

	for index, element := range issues {
		fmt.Printf("\n%d) %s", index + 1, element.Fields.Summary)
	}

}

func login(user *string) (*goji.Client, error) {
	username, password := getCredentials(*user)
	client, err := goji.NewClient("https://servicedesk.softec.ch", username, password)

	if err != nil {
		return nil, err
	}

	u, _, err := client.JiraClient.User.Get(username)

	if err != nil {
		return nil, err
	}

	fmt.Printf("\n\nLogged in as %v\n", u.EmailAddress)
	return client, nil
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
