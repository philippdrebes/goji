package goji

import (
	"github.com/andygrunwald/go-jira"
)

type Client struct {
	JiraClient *jira.Client
}

func NewClient(baseUrl string, username string, password string) *Client {
	client := &Client{}
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	client.JiraClient, _ = jira.NewClient(tp.Client(), baseUrl)
	return client
}

func (c Client) GetAssignedTasks(user string) {

}
