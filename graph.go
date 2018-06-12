package goji

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

func createNodeText() string {
	// issue_key, fields, islink=True
	//summary := fields["summary"]
	//status := fields["status"]
	//
	//if word_wrap == True {
	//	if len(summary) > MAX_SUMMARY_LENGTH {
	//		// split the summary into multiple lines adding a \n to each line
	//		summary = textwrap.fill(fields["summary"], MAX_SUMMARY_LENGTH)
	//	}
	//} else {
	//	//truncate long labels with "...", but only if the three dots are replacing more than two characters
	//	// -- otherwise the truncated label would be taking more space than the original.
	//	if len(summary) > MAX_SUMMARY_LENGTH+2 {
	//		summary = fields["summary"][:MAX_SUMMARY_LENGTH] + "..."
	//		summary = summary.replace('"', '\\"')
	//	}
	//}

	//if islink {
	//	return fmt.Sprintf("\"%s\n(%s)\"", issue_key, summary.encode("utf-8"))
	//	return fmt.Sprintf("\"%s\n(%s)\" [href=\"%s\", fillcolor=\"%s\", style=filled]", issue_key, summary.encode("utf-8"), jira.get_issue_uri(issue_key), get_status_color(status))
	//}

	return ""
}

func walk() {
	// issue_key, graph
	//""" issue is the JSON representation of the issue """

	//issue := jira.get_issue(issue_key)
	//children := []
	//fields := issue["fields"]
	//seen.append(issue_key)

	//if ignore_closed and (fields["status"]["name"] in "Closed") {
	//	return graph
	//}

	//if not traverse and ((project_prefix + '-') not in issue_key) {
	//	return graph
	//}
	//graph.append(create_node_text(issue_key, fields, islink=False))

	//if not ignore_subtasks {
	//	if fields["issuetype"]["name"] == "Epic" and
	//	not
	//	ignore_epic{
	//		issues := jira.query("\"Epic Link\" = \"%s\"" % issue_key)
	//		for subtask in issues:
	//		subtask_key = get_key(subtask)
	//		log(subtask_key + '  = > references epic = > ' + issue_key)
	//		node = "{}->{}[color=orange]".format(
	//		create_node_text(issue_key, fields),
	//		create_node_text(subtask_key, subtask["fields"]))
	//		graph.append(node)
	//		children.append(subtask_key)
	//	}
	//
	//	if fields.has_key("subtasks") and not ignore_subtasks{
	//		for subtask in fields["subtasks"] {
	//			subtask_key = get_key(subtask)
	//			log(issue_key + '  = > has subtask = > ' + subtask_key)
	//			node = "{}->{}[color=blue][label=\"subtask\"]".format (
	//			create_node_text(issue_key, fields),
	//			create_node_text(subtask_key, subtask["fields"]))
	//			graph.append(node)
	//			children.append(subtask_key)
	//		}
	//	}
	//}

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
}

func BuildGraph() {
	//return walk(start_issue_key, [])
}
