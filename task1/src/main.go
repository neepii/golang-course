package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type githubRepoUrl struct {
	owner string
	name  string
}

type githubRepoInfo struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Stargazers int `json:"stargazers_count"`
	Forkers int `json:"fork_count"`
	Creationdate string  `json:"created_at"`
}


func prettyPrintGithubRepo(info githubRepoInfo ) {
	format :=
		"Info about repo:" +
			"Name: %s\n" +
			"Description: %s\n" +
			"Star count: %d\n" +
			"Fork count: %d\n" +
			"Date: %s\n"
	fmt.Printf(format,
		info.Name, info.Description, info.Stargazers, info.Forkers, info.Creationdate)
}

func isValidUrl(str string) bool {
	match, _ := regexp.MatchString("https://github.com/[A-Za-z-]+/[A-Za-z-]+", str)
	return match
}

func parseUrl(url string) (githubRepoUrl, error) {
	if !isValidUrl(url) {
		return githubRepoUrl{}, errors.New("Url with github repo expected")
	}
	split := strings.Split(url, "/")
	// "https:", "", "github.com", "OWNER", "NAME"
	return githubRepoUrl{owner: split[3], name: split[4]}, nil
}

func usage() {
	usage_message :=
		"usage: %s [url_of_github_repo]\n" +
			"The flags are:\n" +
			"\t--help - print this message\n"

	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, usage_message, os.Args[0])
	os.Exit(2)
}

func getGithubRepoInfo(url string) (githubRepoInfo, error) {
	var dat githubRepoInfo
	repo, err := parseUrl(url)
	if err != nil {
		return dat, err
	}
	res, err := http.Get("https://api.github.com/repos/" + repo.owner + "/" + repo.name)
	if err != nil {
		return dat, errors.New("Error making http request")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return dat, errors.New("Bad request")
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return dat, errors.New("Can't read response body")
	}

	if err := json.Unmarshal(bodyBytes, &dat); err != nil {
		panic(err)
	}

	return dat, nil
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Input file is missing")
		os.Exit(1)
	}
	args := flag.Args()
	url := args[0]

	info, err := getGithubRepoInfo(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	prettyPrintGithubRepo(info)
}
