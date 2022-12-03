package main

import (
	"encoding/xml"
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

type AllReleases struct {
	XMLName  xml.Name  `xml:"a"`
	Releases []Release `xml:"r"`
}

type Release struct {
	Version   string `xml:"v"`
	Stability string `xml:"s"`
}

func main() {
	pkg := os.Getenv("INPUT_PACKAGE")
	stbl := os.Getenv("INPUT_STABILITY")
	filter := os.Getenv("INPUT_FILTER")

	if pkg == "" {
		log.Fatal("Package name is required")
	}

	if stbl == "" {
		stbl = "stable"
	}

	url := fmt.Sprintf(`https://pecl.php.net/rest/r/%s/allreleases.xml`, pkg)

	log.Print(url)

	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "get-pecl-package-version-action")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var allReleases AllReleases
	xml.Unmarshal(body, &allReleases)

	var versions []*semver.Version
	for _, release := range allReleases.Releases {
		if release.Stability != stbl {
			continue
		}

		version := strings.Trim(release.Version, "vV")

		if filter != "" && !strings.HasPrefix(version, filter) {
			continue
		}

		matched, _ := regexp.MatchString(`^[vV]*[0-9]+\.[0-9]+\.[0-9]+$`, version)

		if matched {
			versions = append(versions, semver.New(release.Version))
		}
	}

	if len(versions) == 0 {
		log.Fatal(fmt.Sprintf(`Unable to find versions for package %s`, pkg))
	}

	semver.Sort(versions)

	setOutput("version", versions[len(versions)-1].String())
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

