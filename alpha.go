package main

import "time"

type node struct {
	Identifier string
	URI        string
	Lastseenat time.Time
}
