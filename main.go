package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	logfile  *string
	repo_url *string
	title    *string
	desc     *string
)

// facilitates CLI arguments
func init() {
	// TODO Check for preset env variables like GITHUB_TOKEN, terminate if absent

	logfile = flag.String("send", "", "a string")  // a logfile
	repo_url = flag.String("repo", "", "a string") // the project
	title = flag.String("title", "", "a string")
	desc = flag.String("desc", "", "a string")
	flag.Parse()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var (
		org     string
		project string
	)
	// 1. Authentication
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_USER")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// 2. Add missing data
	if *repo_url == "" {
		fmt.Print("Enter organization/repository ")
		fmt.Scanf("%s", repo_url)
	}
	// get the org_name and the project name
	q := strings.Split(*repo_url, "/")
	org = q[0]
	project = q[1]

	if *logfile == "" {
		fmt.Print("Log-file: ")
		fmt.Scanf("%s", logfile)
	}

	if *title == "" {
		fmt.Print("Issue title: ")
		*title, _ = reader.ReadString('\n')
	}

	if *desc == "" {
		fmt.Println("Issue description:")
		*desc, _ = reader.ReadString('\n')
	}

	j := github.IssueRequest{
		Title: title,
		Body:  desc}
	issue_struct, _, err := client.Issues.Create(ctx, org, project, &j)
	if err == nil {
		fmt.Println("Issue#", issue_struct.ID, " created")
	}
}
