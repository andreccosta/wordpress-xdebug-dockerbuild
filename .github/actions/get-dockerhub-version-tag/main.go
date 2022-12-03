package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/coreos/go-semver/semver"
)

type dhtag struct {
	Name string `json:"name"`
}

type dhrepo struct {
	Count   int     `json:"count"`
	Next    string  `json:"next"`
	Results []dhtag `json:"results"`
}

var org string
var repo string

var httpClient = http.Client{
	Timeout: time.Second * 2,
}

func getTags(url string) (tags []*semver.Version, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "get-dockerhub-version-tag-action")

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	dhrepo1 := dhrepo{}
	unmarshalErr := json.Unmarshal(body, &dhrepo1)
	if unmarshalErr != nil {
		log.Fatal(unmarshalErr)
	}

	for _, tag := range dhrepo1.Results {
		matched, _ := regexp.MatchString(`^[vV]*[0-9]+\.[0-9]+\.[0-9]+$`, tag.Name)

		if matched {
			log.Printf("Matched %s", tag.Name)
			tags = append(tags, semver.New(strings.Trim(tag.Name, "vV")))
		}
	}

	if len(tags) > 0 {
		return tags, nil
	} else if dhrepo1.Next != "" {
		return getTags(dhrepo1.Next)
	} else {
		return nil, errors.New(fmt.Sprintf(`Unable to find tags for %s/%s`, org, repo))
	}
}

func main() {
	org = os.Getenv("INPUT_ORG")
	repo = os.Getenv("INPUT_REPO")
	page_size := 100

	if org == "" {
		org = "library"
	}

	if repo == "" {
		log.Fatal("Repo is required")
	}

	url := fmt.Sprintf(`https://hub.docker.com/v2/repositories/%s/%s/tags/?page=1&page_size=%d&ordering=last_updated`, org, repo, page_size)
	tags, err := getTags(url)

	if err != nil {
		log.Fatal(err)
	}

	semver.Sort(tags)

	setOutput("tag", tags[len(tags)-1].String())
}

func setOutput(name, value string) error {
	file := os.Getenv("GITHUB_OUTPUT")
	if file == "" {
		return errors.New("GITHUB_OUTPUT env variable not specified")
	}

	return appendToFile(file, fmt.Sprintf("%s=%s\n", name, value))
}

func appendToFile(file, content string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	_, err = f.WriteString(content)
	closeErr := f.Close()
	if err != nil {
		return err
	}

	return closeErr
}
