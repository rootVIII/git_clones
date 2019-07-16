package main

// rootVIII
// GET all of a user's repositories
// USAGE:
// ./git_clones -u <username>

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var user string
	flag.StringVar(&user, "u", "", "Github Username")
	flag.Parse()
	if len(user) < 1 {
		fmt.Println("USAGE: ./git_clones -u <Github Username>")
		os.Exit(1)
	}
	clone(makeRequest(user), user)
}

func clone(html string, username string) {
	lines := strings.Fields(html)
	git := "git"
	clone := "clone"
	for i := range lines {
		if strings.Contains(lines[i], "href") && strings.Contains(lines[i], username) {
			elements := strings.Split(lines[i], "/")
			if len(elements) != 3 || elements[1] != username {
				continue
			}
			url := "https://github.com/" + username + "/"
			url += elements[2][:len(elements[2])-1] + ".git"
			_, err := exec.Command(git, clone, url).Output()
			if err != nil {
				fmt.Printf("Unable to download %s", url)
			} else {
				fmt.Printf("Cloning:  %s\n", url)
			}
		}
	}
}

func makeRequest(username string) string {
	url := "https://github.com/" + username + "?&tab=repositories&q=&type=source"
	response, err := http.Get(url)
	if err != nil || response.StatusCode != 200 {
		fmt.Printf("unable to retrieve Github page for " + username)
		os.Exit(1)
	}
	body, _ := ioutil.ReadAll(response.Body)
	return string(body)
}
