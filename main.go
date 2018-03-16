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

type fv func(client *goji.Client, user string)

type Action struct {
	description string
	function    fv
}

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

	var actions []Action
	actions = append(actions, Action{"Display assigned tasks", displayAssignedTasks})
	actions = append(actions, Action{"Quit", nil})

	for {
		selectedAction := promptForAction(actions)
		if selectedAction != nil {
			if selectedAction.description == "Quit" {
				os.Exit(2)
			} else {
				selectedAction.function(client, *user)
			}
		}
	}
}

func promptForAction(actions []Action) *Action {
	fmt.Println()
	for index, element := range actions {
		fmt.Printf("%d) %s\n", index+1, element.description)
	}

	var input int
	n, err := fmt.Scanln(&input)
	if n < 1 || err != nil || n > len(actions) {
		fmt.Println("invalid input")
		return nil
	}

	return &actions[input - 1]
}

func displayAssignedTasks(client *goji.Client, user string) {
	issues, err := client.GetAssignedTasks(user)
	if err != nil {
		fmt.Printf("\nError while trying to get assigned issues.\n%v\n", err)
		return
	}

	for _, element := range issues {
		fmt.Printf("\n%s - %s", element.Key, element.Fields.Summary)
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

	if len(user) == 0 {
		fmt.Print("Jira Username: ")
		username, _ = r.ReadString('\n')
	}

	fmt.Print("Jira Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password = string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password)
}
