package main

// rootVIII gitclones.go

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type clone interface {
	httpReq() string
	findRepositories(requestData string)
	download(repository string)
	setURL()
}

type gitClone struct {
	baseURL  string
	userURL  string
	username string
}

func (gc gitClone) httpReq() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", gc.userURL, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("request failed: %s\n", gc.userURL)
		os.Exit(1)
	}
	rText, err := ioutil.ReadAll(resp.Body)
	return string(rText)
}

func (gc *gitClone) setURL() {
	gc.userURL = gc.baseURL + gc.username
	gc.userURL += "?&tab=repositories&q=&type=source"
}

func (gc *gitClone) download(repository string) {
	projectURL := gc.baseURL + gc.username
	projectURL += fmt.Sprintf("/%s.git", repository)
	fmt.Printf("Cloning %s\n", repository)
	exec.Command("git", "clone", projectURL).Output()
}

func (gc *gitClone) findRepositories(requestData string) {
	pattern := fmt.Sprintf(`<a\s?href\W+%s\S+`, gc.username)
	result := regexp.MustCompile(pattern)
	links := result.FindAllString(requestData, -1)
	for _, tag := range links {
		gc.download(strings.Split(tag[:len(tag)-1], "/")[2])
	}
}

func main() {
	user := flag.String("u", "", "Provide a valid Github username")
	flag.Parse()
	if len(*user) < 1 {
		fmt.Println("required argument: -u <Github username>")
		os.Exit(1)
	}
	var gitclone clone
	gitclone = &gitClone{username: *user, baseURL: "https://github.com/"}
	gitclone.setURL()
	pageData := gitclone.httpReq()
	gitclone.findRepositories(pageData)
}
