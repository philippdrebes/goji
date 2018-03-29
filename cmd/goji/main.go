package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/akamensky/argparse"
	"golang.org/x/crypto/ssh/terminal"
	"os/exec"
	"runtime"
	"github.com/atotto/clipboard"
	"github.com/philippdrebes/goji/pkg/goji"
)

var clear map[string]func() //create a map for storing clear funcs

type fv func(client *goji.Client)

type Action struct {
	Key         string
	Description string
	Function    fv
}

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok { //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func main() {
	CallClear()
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
	actions = append(actions, Action{"assignedTasks", "Display assigned tasks", displayAssignedTasks})
	actions = append(actions, Action{"quit", "Quit", nil})

	for {
		selectedAction := promptForAction(actions)
		if selectedAction != nil {
			if selectedAction.Description == "Quit" {
				os.Exit(2)
			} else {
				selectedAction.Function(client)
			}
		}
		CallClear()
	}
}

func promptForAction(actions []Action) *Action {
	fmt.Println()
	for index, element := range actions {
		fmt.Printf("%d) %s\n", index+1, element.Description)
	}

	var input int
	n, err := fmt.Scanln(&input)
	if n < 1 || err != nil || n > len(actions) {
		fmt.Println("invalid input")
		return nil
	}

	return &actions[input-1]
}

func displayAssignedTasks(client *goji.Client) {
	var actions []Action
	clipboardAction := Action{"clipboard", "Copy to clipboard", nil}
	refreshAction := Action{"refresh", "Refresh", nil}
	backAction := Action{"back", "Back", nil}
	actions = append(actions, clipboardAction)
	actions = append(actions, refreshAction)
	actions = append(actions, backAction)

	for {
		issues, err := client.GetAssignedTasks(client.CurrentUser.Name)
		if err != nil {
			fmt.Printf("\nError while trying to get assigned issues.\n%v\n", err)
			return
		}

		for _, element := range issues {
			fmt.Printf("\n%s - %s", element.Key, element.Fields.Summary)
		}
		fmt.Println()

		selectedAction := promptForAction(actions)
		if selectedAction.Key == clipboardAction.Key {
			clipboard.WriteAll("asdf") // todo :shipit:
		} else if selectedAction.Key == backAction.Key {
			fmt.Println()
			return
		}
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

	CallClear()
	fmt.Printf("Logged in as %v\n", u.EmailAddress)
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
