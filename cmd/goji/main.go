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
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = clear["linux"]
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
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
	fmt.Println("\nGoji")

	parser := argparse.NewParser("print", "Prints provided string to stdout")
	user := parser.String("u", "user", &argparse.Options{Required: false, Help: "Username"})
	err := parser.Parse(os.Args)

	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return
	}

	goji.GetConfig()
	client, err := login(user)

	if client == nil || err != nil {
		fmt.Printf("\nError while trying to log in.\n%v\n", err)
		return
	}

	CallClear()

	var actions []Action
	actions = append(actions, Action{"assignedTasks", "Display assigned tasks", displayAssignedTasks})
	actions = append(actions, Action{"quit", "Quit", nil})

	for {
		fmt.Println("\nHello Goji!")
		fmt.Printf("Logged in as %v\n", client.CurrentUser.EmailAddress)
		selectedAction := promptForAction(actions)
		if selectedAction != nil {
			if selectedAction.Description == "Quit" {
				os.Exit(0)
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
	CallClear()

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

		issueSummary := ""
		for _, element := range issues {
			issue := fmt.Sprintf("\n%s: %s", element.Key, element.Fields.Summary)
			issueSummary += issue
		}
		fmt.Printf("%s\n", issueSummary)
		fmt.Println()

		selectedAction := promptForAction(actions)
		if selectedAction.Key == clipboardAction.Key {
			clipboard.WriteAll(issueSummary)
		} else if selectedAction.Key == backAction.Key {
			fmt.Println()
			return
		}
	}
}

func login(user *string) (*goji.Client, error) {
	r := bufio.NewReader(os.Stdin)

	config := goji.GetConfig()
	uname := config.Username

	if len(*user) > 0 {
		uname = *user
	}

	url := config.Url
	if len(url) == 0 {
		fmt.Print("Jira Url: ")
		url, _ = r.ReadString('\n')
	}

	url = strings.TrimSpace(url)
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	username, password := getCredentials(uname)
	client, err := goji.NewClient(url, username, password)

	if err != nil {
		return nil, err
	}

	if len(config.Url) == 0 {
		fmt.Print("\nSave as default? [Y/n]: ")
		save, _ := r.ReadString('\n')

		if  strings.ToLower(strings.TrimSpace(save)) != "n" {
			config.Url = url
			config.Username = username
			goji.SaveConfig(config)
		}
	}

	CallClear()
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
