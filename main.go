package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/andygrunwald/go-jira"
	"github.com/bgentry/speakeasy"
)

func main() {
	fmt.Println("Hello Goji!")

	// Create new parser object
	parser := argparse.NewParser("print", "Prints provided string to stdout")
	// Create string flag
	user := parser.String("u", "user", &argparse.Options{Required: false, Help: "Username"})
	// Parse input
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
	juser, _, _ := jiraClient.User.Get("pd")

	if juser != nil {
		fmt.Printf("Version: %s\n", juser.Name)
	}
}

func getCredentials(user string) (string, string) {
	var username string
	var password string

	if len(user) == 0 {
		fmt.Print("Username: ")
		fmt.Scanln(&username)
	}

	password, err := speakeasy.Ask("Password: ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return username, password
}
