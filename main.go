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
	repoURL *string
	title   *string
	desc    *string
	cmd     *string
)

// facilitates CLI arguments
func init() {
	if os.Getenv("GITHUB_USER") == "" {
		fmt.Fprintln(os.Stderr, "Get a Github personal access token and create an environment variable GITHUB_USER and try again.")
		fmt.Println("You can create one right here https://github.com/settings/tokens")
	}
	repoURL = flag.String("repo", "", "a string") // the project
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
	if *repoURL == "" {
		fmt.Print("Enter organization/repository: ")
		fmt.Scanf("%s", repoURL)
	}
	// get the org_name and the project name
	q := strings.Split(*repoURL, "/")
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
	issueStruct, _, err := client.Issues.Create(ctx, org, project, &j)
	if err == nil {
		fmt.Println("Issue#", *issueStruct.ID, " created")
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
		fmt.Print("Add some description: ")
		de, _ := reader.ReadString('\n')
		g := github.Gist{
			Description: &de,
			Files:       m}
		gistRet, _, errRet := client.Gists.Create(ctx, &g)
		fmt.Println(*gistRet.HTMLURL)
		if errRet != nil {
			panic(errRet)
		}
	}
}
