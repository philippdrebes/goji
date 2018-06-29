package goji

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
)

type Client struct {
	JiraClient  *jira.Client
	CurrentUser *jira.User
}

func NewClient(baseUrl string, username string, password string) (*Client, error) {
	client := &Client{}
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	jclient, err := jira.NewClient(tp.Client(), baseUrl)
	if err != nil {
		return nil, err
	}
	client.JiraClient = jclient

	usr, _, err := client.JiraClient.User.Get(username)
	if err != nil {
		return nil, err
	}
	client.CurrentUser = usr

	return client, nil
}

func (c Client) GetAssignedTasks(user string) ([]jira.Issue, error) {
	if len(user) == 0 {
		user = c.CurrentUser.Name
	}

	issues, _, err := c.JiraClient.Issue.Search(fmt.Sprintf("assignee in (%s)", user), nil)
	return issues, err
}

func (c Client) GetCurrentUserWorklog(daysInThePast int) ([]jira.Issue, error) {
	var jql string
	hours := 24 * daysInThePast
	jql = fmt.Sprintf("updatedDate >= startOfDay(\"-%dh\") and updatedDate < endOfDay(\"-%dh\") AND assignee was currentUser() ORDER BY updatedDate ASC", hours, hours)
	issues, _, err := c.JiraClient.Issue.Search(jql, nil)
	return issues, err
}
