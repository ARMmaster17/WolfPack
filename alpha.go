package main

import "time"

type node struct {
	identifier string
	uri        string
	lastseenat time.Time
}
