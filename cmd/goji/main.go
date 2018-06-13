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
	"github.com/philippdrebes/goji"
	"io/ioutil"
	"path"
	"time"
	"github.com/skratchdot/open-golang/open"
	"path/filepath"
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
	err := parser.Parse(os.Args)

	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return
	}

	client, err := login()

	if client == nil || err != nil {
		fmt.Printf("\nError while trying to log in.\n%v\n", err)
		return
	}

	CallClear()

	var actions []Action
	actions = append(actions, Action{"assignedTasks", "Display assigned tasks", displayAssignedTasks})
	actions = append(actions, Action{"linkedIssueGraph", "Display graph of linked issues", createLinkedIssueGraph})
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

func createLinkedIssueGraph(client *goji.Client) {
	issue, _, _ := client.JiraClient.Issue.Get("SWECOWEB-14", nil)
	graph := goji.BuildGraph(client.JiraClient, issue)

	dotFile, err := ioutil.TempFile(os.TempDir(), "goji-deps")
	if err != nil {
		fmt.Printf("\nError while opening temp dotFile.\n%v\n", err)
		return
	}
	defer os.Remove(dotFile.Name())

	err = ioutil.WriteFile(dotFile.Name(), []byte(graph.String()), 0755)
	png, err := exec.Command("dot", "-Tpng", dotFile.Name()).Output()
	if err != nil {
		fmt.Printf("\nError while creating graph.\n%v\n", err)
		return
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	now := time.Now()
	pngFile := path.Join(exPath, fmt.Sprintf("%s_%d-%02d-%02d.png", issue.Key, now.Year(), now.Month(), now.Day()))

	err = ioutil.WriteFile(pngFile, png, 0755)
	//defer os.Remove(pngFile)

	if err != nil {
		fmt.Printf("\nError while saving graph.\n%v\n", err)
		return
	}

	open.Run(pngFile)
	return
}

func login() (*goji.Client, error) {
	url, username, password := getCredentials()
	client, err := goji.NewClient(url, username, password)

	if err != nil {
		return nil, err
	}

	CallClear()
	return client, nil
}

func getCredentials() (string, string, string) {
	var url string
	var username string
	var password string

	config := goji.GetConfig()

	url = config.Url
	username = config.Username

	r := bufio.NewReader(os.Stdin)

	if len(url) == 0 {
		fmt.Print("Jira Url: ")
		url, _ = r.ReadString('\n')
	}

	url = strings.TrimSpace(url)
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	if len(username) == 0 {
		fmt.Print("Jira Username: ")
		username, _ = r.ReadString('\n')
	}
	username = strings.TrimSpace(username)

	fmt.Print("Jira Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password = strings.TrimSpace(string(bytePassword))

	if len(config.Url) == 0 {
		fmt.Print("\nSave as default? [Y/n]: ")
		save, _ := r.ReadString('\n')

		if strings.ToLower(strings.TrimSpace(save)) != "n" {
			config.Url = url
			config.Username = username
			goji.SaveConfig(config)
		}
	}

	return url, username, password
}
