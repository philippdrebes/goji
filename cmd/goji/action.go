package main

import "github.com/philippdrebes/goji/pkg/goji"

type fv func(client *goji.Client)

type Action struct {
	Key         string
	Description string
	Function    fv
}