package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	repo_url *string
	title    *string
	desc     *string
	cmd      *string
)

// facilitates CLI arguments
func init() {
	// TODO Check for preset env variables like GITHUB_TOKEN, terminate if absent

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

	if *title == "" {
		fmt.Print("Issue title: ")
		*title, _ = reader.ReadString('\n')
	}

	if *desc == "" {
		fmt.Println("Issue description:")
		*desc, _ = reader.ReadString('\n')
	}
	// Creates an issue
	j := github.IssueRequest{
		Title: title,
		Body:  desc}
	issue_struct, _, err := client.Issues.Create(ctx, org, project, &j)
	if err == nil {
		fmt.Println("Issue#", *issue_struct.ID, " created")
	}

	// Capture some output
	if len(flag.Args()) != 0 {
		fmt.Println(flag.Args(), " will be executed")
		// TODO get rid of arg
		arg := flag.Args()[:]
		cmdObj := exec.Command(string(flag.Args()[0]), arg[1:]...)
		logs, err := cmdObj.CombinedOutput()
		stringlogs := string(logs)
		if err != nil {
			panic(err)
		}
		content := github.GistFile{
			Filename: &flag.Args()[0],
			// TODO Fix types from array of bytes to pointer to string
			Content: &stringlogs,
		}
		m := make(map[github.GistFilename]github.GistFile)
		m["logs"] = content
		de, _ := reader.ReadString('\n')
		g := github.Gist{
			Description: &de,
			Files:       m}
		client.Gists.Create(ctx, &g)
	}
}
