package goji

import (
	"github.com/andygrunwald/go-jira"
	"github.com/awalterschulze/gographviz"
	"fmt"
	"strings"
)

const (
	inward  int = 0
	outward int = 1
)

func wrapInQuotes(value string) string {
	return "\"" + value + "\""
}

func getStatusColor(issue *jira.Issue) string {
	status := strings.ToUpper(issue.Fields.Status.Name)
	if status == "IN PROGRESS" {
		return "yellow"
	} else if status == "DONE" {
		return "green"
	}
	return "white"
}

func addNode(graph *gographviz.Graph, issue *jira.Issue) {
	if graph.IsNode(issue.Key) == true {
		return
	}

	attrs := make(map[string]string)
	attrs[string(gographviz.Shape)] = "box"
	attrs[string(gographviz.FillColor)] = getStatusColor(issue)
	attrs[string(gographviz.Style)] = "filled"
	attrs[string(gographviz.Label)] = wrapInQuotes(fmt.Sprintf("%s\n%s", issue.Key, issue.Fields.Summary))

	graph.AddNode("G", wrapInQuotes(issue.Key), attrs)
}

func walk(client *jira.Client, issue *jira.Issue, graph *gographviz.Graph) *gographviz.Graph {
	if graph == nil {
		graph = gographviz.NewGraph()
		if err := graph.SetDir(true); err != nil {
			panic(err)
		}
		graph.Attrs.Add(string(gographviz.DPI), "300")
	}

	addNode(graph, issue)

	for _, li := range issue.Fields.IssueLinks {
		var direction int
		var linkedIssue *jira.Issue
		if li.InwardIssue != nil {
			linkedIssue = li.InwardIssue
			direction = inward
		} else if li.OutwardIssue != nil {
			linkedIssue = li.OutwardIssue
			direction = outward
		}

		// create node
		linkedIssue, _, err := client.Issue.Get(linkedIssue.ID, nil)
		if err != nil {
			panic(err)
		}

		if graph.IsNode(wrapInQuotes(linkedIssue.Key)) == false {
			walk(client, linkedIssue, graph)
		}

		// create edge
		attrs := make(map[string]string)
		if direction == inward {
			attrs[string(gographviz.Label)] = wrapInQuotes(li.Type.Inward)
		} else if direction == outward {
			attrs[string(gographviz.Label)] = wrapInQuotes(li.Type.Outward)
		}
		graph.AddEdge(wrapInQuotes(issue.Key), wrapInQuotes(linkedIssue.Key), true, attrs)
	}

	return graph
}

func BuildGraph(client *jira.Client, issue *jira.Issue) *gographviz.Graph {
	return walk(client, issue, nil)
}
