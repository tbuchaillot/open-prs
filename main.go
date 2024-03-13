package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
)

type PullRequest struct {
	URL   string `json:"html_url"`
	Title string `json:"title"`
	User  struct {
		Login string `json:"login"`
	} `json:"user"`
}

var (
	repository = flag.String("repository", "", "github team repository")
	org        = flag.String("org", "", "github organization name")
	output     = flag.String("output", "stdout", "ouput type (stdout,csv) ")
)

func main() {
	flag.Parse()
	if *repository == "" || *org == "" {
		flag.Usage()
		os.Exit(1)
	}

	getPullRequests(*org, *repository, *output)
}

func getPullRequests(owner, repo, output string) {
	baseURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?state=open", owner, repo)

	prDetailsPerContributor := make(map[string]*ContributorPRCount)

	for {
		pullRequests, nextURL, err := fetchPullRequests(baseURL)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		for _, pr := range pullRequests {
			if _, exists := prDetailsPerContributor[pr.User.Login]; !exists {
				prDetailsPerContributor[pr.User.Login] = &ContributorPRCount{Login: pr.User.Login, URLs: []string{}}
			}
			contributor := prDetailsPerContributor[pr.User.Login]
			contributor.Count++
			contributor.URLs = append(contributor.URLs, pr.URL)
		}

		if nextURL == "" {
			break
		}
		baseURL = nextURL
	}

	var contributors Contributors
	for _, details := range prDetailsPerContributor {
		contributors = append(contributors, *details)
	}

	// Sort the slice by PR count in descending order
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Count > contributors[j].Count
	})

	contributors.Output(output)
}

func fetchPullRequests(url string) ([]PullRequest, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching pull requests: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error reading response body: %w", err)
	}

	var pullRequests []PullRequest
	if err := json.Unmarshal(body, &pullRequests); err != nil {
		return nil, "", fmt.Errorf("error parsing JSON: %w", err)
	}

	nextURL := getNextPageURL(resp.Header.Get("Link"))

	return pullRequests, nextURL, nil
}

func getNextPageURL(linkHeader string) string {
	if linkHeader == "" {
		return ""
	}

	links := strings.Split(linkHeader, ",")
	for _, link := range links {
		segments := strings.Split(link, ";")
		if len(segments) == 2 && strings.TrimSpace(segments[1]) == `rel="next"` {
			nextURL := strings.Trim(segments[0], " <>")
			return nextURL
		}
	}
	return ""
}
