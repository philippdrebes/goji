package goji

type fv func(client *Client)

type Action struct {
	Key         string
	Description string
	Function    fv
}