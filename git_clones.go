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
	baseURL, username string
}

func (gc gitClone) httpReq() *[]byte {
	userURL := gc.baseURL + gc.username + "?&tab=repositories&q=&type=source"
	client := &http.Client{}
	req, err := http.NewRequest("GET", userURL, nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	rBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return &rBytes
}

func (gc gitClone) download(repository []byte) {
	projectURL := gc.baseURL + gc.username
	projectURL += fmt.Sprintf("/%s.git", string(repository))
	fmt.Printf("Cloning %s\n", repository)
	_, err := exec.Command("git", "clone", projectURL).Output()
	if err != nil {
		fmt.Printf("Failed to download repository: %s\n", repository)
	}
}

func (gc gitClone) findRepositories(requestData []byte) {
	pattern := fmt.Sprintf(`<a\s?href\W+%s\S+`, gc.username)
	result := regexp.MustCompile(pattern)
	links := result.FindAll(requestData, -1)
	for _, tag := range links {
		gc.download(bytes.Split(tag[:len(tag)-1], []byte{0x2f})[2])
	}
}

func (gc gitClone) runGitClones() {
	pageData := gc.httpReq()
	gc.findRepositories(*pageData)
}

func main() {
	user := flag.String("u", "", "Provide a valid Github username")
	flag.Parse()
	if len(*user) < 1 {
		fmt.Println("required argument: -u <Github username>")
		os.Exit(1)
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
	}()

	var gitclone = &gitClone{username: *user, baseURL: "https://github.com/"}
	gitclone.runGitClones()
}
