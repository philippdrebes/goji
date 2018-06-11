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

func BuildGraph() {

}
