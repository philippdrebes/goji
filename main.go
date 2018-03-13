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

	tp := jira.CookieAuthTransport{
		Username: username,
		Password: password,
		AuthURL:  "https://servicedesk.softec.ch/rest/auth/1/session",
	}

	jiraClient, _ := jira.NewClient(tp.Client(), "https://servicedesk.softec.ch")

	authenticated := jiraClient.Authentication.Authenticated()
	session, _ := jiraClient.Authentication.GetCurrentUser()
	u, _, _ := jiraClient.User.Get("pd")

	fmt.Println(authenticated)
	fmt.Println(session.Name)

	if u != nil {
		fmt.Printf("\nEmail: %v\nSuccess!\n", u.EmailAddress)
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
