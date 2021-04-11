package main

import (
	"encoding/json"
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
	Results []dhtag `json:"results"`
}

func main() {
	org := os.Getenv("INPUT_ORG")
	repo := os.Getenv("INPUT_REPO")

	if org == "" {
		org = "library"
	}

	if repo == "" {
		log.Fatal("Repo is required")
	}

	url := fmt.Sprintf(`https://hub.docker.com/v2/repositories/%s/%s/tags/?page_size=1000&ordering=last_updated`, org, repo)

	dhClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "get-dockerhub-version-tag-action")

	res, getErr := dhClient.Do(req)
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

	var tags []*semver.Version
	for _, tag := range dhrepo1.Results {
		matched, _ := regexp.MatchString(`^[vV]*[0-9]+\.[0-9]+\.[0-9]+$`, tag.Name)

		if matched {
			log.Printf("Matched %s", tag.Name)
			tags = append(tags, semver.New(strings.Trim(tag.Name, "vV")))
		}
	}

	if len(tags) == 0 {
		log.Fatal(fmt.Sprintf(`Unable to find tags for %s/%s`, org, repo))
	}

	semver.Sort(tags)
	fmt.Println(fmt.Sprintf(`::set-output name=tag::%s`, tags[len(tags)-1]))
}
