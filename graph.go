package goji

import (
	"github.com/andygrunwald/go-jira"
	"github.com/awalterschulze/gographviz"
)

func getStatusColor() string {
	//status := status_field["statusCategory"]["name"].upper()
	status := "IN PROGRESS"
	if status == "IN PROGRESS" {
		return "yellow"
	} else if status == "DONE" {
		return "green"
		return "white"
	}
	return ""
}

func walk(client *jira.Client, issue *jira.Issue, graph *gographviz.Graph) *gographviz.Graph {
	if graph == nil {
		graph = gographviz.NewGraph()
	}

	graph.AddNode("G", issue.Key, nil)

	for _, linkedIssue := range issue.Fields.IssueLinks {
		linkedIssueObj, _, err := client.Issue.Get(linkedIssue.ID, nil)
		if err != nil {
			// todo log err
		}
		walk(client, linkedIssueObj, graph)
	}

	//if fields.has_key("issuelinks"){
	//	for other_link in fields["issuelinks"]{
	//		result = process_link(fields, issue_key, other_link)
	//		if result is not None {
	//			children.append(result[0])
	//			if result[1] is not None {
	//				graph.append(result[1])
	//			}
	//		}
	//	}
	//}
	//// now construct graph data for all subtasks and links of this issue
	//for child in (x for x in children if x not in seen) {
	//	walk(child, graph)
	//	return graph
	//}

	return nil
}

func BuildGraph(client *jira.Client, issue *jira.Issue) *gographviz.Graph {
	return walk(client, issue, nil)
}
