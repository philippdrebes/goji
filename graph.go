package goji

import (
	"github.com/andygrunwald/go-jira"
	"github.com/awalterschulze/gographviz"
	"fmt"
)

const (
	inward  int = 0
	outward int = 1
)

func wrapInQuotes(value string) string {
	return "\"" + value + "\""
}

func getStatusColor(issue *jira.Issue) string {
	status := issue.Fields.Status.Name
	if status == "IN PROGRESS" {
		return "yellow"
	} else if status == "DONE" {
		return "green"
	}
	return "white"
}

func addNode(graph *gographviz.Graph, issue *jira.Issue) {
	if graph.IsNode(issue.Key) == false {
		attrs := make(map[string]string)
		attrs[string(gographviz.Shape)] = "box"
		//attrs[string(gographviz.HREF)] = ""
		//attrs[string(gographviz.FillColor)] = ""
		//attrs[string(gographviz.Style)] = "filled"

		graph.AddNode("G", wrapInQuotes(issue.Key), attrs)
	}
}

func walk(client *jira.Client, issue *jira.Issue, graph *gographviz.Graph) *gographviz.Graph {
	if graph == nil {
		graph = gographviz.NewGraph()
		if err := graph.SetDir(true); err != nil {
			panic(err)
		}
	}

	addNode(graph, issue)

	for _, li := range issue.Fields.IssueLinks {
		fmt.Printf("issue: %s\n", issue.Key)

		var direction int
		var linkedIssue *jira.Issue

		if li.InwardIssue != nil {
			direction = inward
			linkedIssue = li.InwardIssue
		} else if li.OutwardIssue != nil {
			direction = outward
			linkedIssue = li.OutwardIssue
		}

		linkedIssue, _, err := client.Issue.Get(linkedIssue.ID, nil)
		if err != nil {
			panic(err)
		}

		fmt.Printf("direction: %d\n", direction)
		fmt.Printf("linked issue: %s\n", linkedIssue.Key)

		if graph.IsNode(wrapInQuotes(linkedIssue.Key)) == false {
			walk(client, linkedIssue, graph)
		}
		graph.AddEdge(wrapInQuotes(issue.Key), wrapInQuotes(linkedIssue.Key), true, nil)
	}

	return graph
}

func BuildGraph(client *jira.Client, issue *jira.Issue) *gographviz.Graph {
	return walk(client, issue, nil)
}
