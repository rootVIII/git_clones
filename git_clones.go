package main

// rootVIII gitclones.go

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
)

type gitClone struct {
	baseURL  string
	userURL  string
	username string
}

func (gc gitClone) httpReq() *[]byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", gc.userURL, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("request failed: %s\n", gc.userURL)
		os.Exit(1)
	}
	rBytes, err := ioutil.ReadAll(resp.Body)
	return &rBytes
}

func (gc *gitClone) setURL() {
	gc.userURL = gc.baseURL + gc.username
	gc.userURL += "?&tab=repositories&q=&type=source"
}

func (gc *gitClone) download(repository []byte) {
	projectURL := gc.baseURL + gc.username
	projectURL += fmt.Sprintf("/%s.git", string(repository))
	fmt.Printf("Cloning %s\n", repository)
	exec.Command("git", "clone", projectURL).Output()
}

func (gc *gitClone) findRepositories(requestData []byte) {
	pattern := fmt.Sprintf(`<a\s?href\W+%s\S+`, gc.username)
	result := regexp.MustCompile(pattern)
	links := result.FindAll(requestData, -1)
	for _, tag := range links {
		gc.download(bytes.Split(tag[:len(tag)-1], []byte{0x2f})[2])
	}
}

func main() {
	user := flag.String("u", "", "Provide a valid Github username")
	flag.Parse()
	if len(*user) < 1 {
		fmt.Println("required argument: -u <Github username>")
		os.Exit(1)
	}

	var gitclone = &gitClone{username: *user, baseURL: "https://github.com/"}
	gitclone.setURL()
	pageData := gitclone.httpReq()
	gitclone.findRepositories(*pageData)
}
