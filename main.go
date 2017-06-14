package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	logfile  *string
	repo_url *string
)

// facilitates CLI arguments
func init() {
	// TODO Check for preset env variables like GITHUB_TOKEN, terminate if absent
	logfile = flag.String("send", "", "a string")  // a logfile
	repo_url = flag.String("repo", "", "a string") // the project
	flag.Parse()
}

// authentication with a Github personal access token
func passwordlessAuth() string {
	key := os.Getenv("GITHUB_USER")
	if key == "" {
		fmt.Fprintf(os.Stderr, "Set GITHUB_USER env variable with a github personal token for faster access")

	}
	return key
}

func main() {
	var (
		org     string
		project string

		key string
	)
	// 1. Authentication
	key = passwordlessAuth()

	// 2. Add missing data
	if *repo_url == "" {
		fmt.Print("Repository URL: ")
		fmt.Scanf("%s", repo_url)
		// get the org_name and the project name
		l := strings.Split(*repo_url, "github.com")[1]
		q := strings.Split(l, "/")
		org = q[1]
		project = q[2]
	}

	if *logfile == "" {
		fmt.Print("Log-file: ")
		fmt.Scanf("%s", logfile)
	}
}
