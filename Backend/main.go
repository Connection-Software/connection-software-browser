package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var urlRegex string = "((http|https)://)(www.)?" +
	"[a-zA-Z0-9@:%._\\+~#?&//=]" +
	"{2,256}\\.[a-z]" +
	"{2,6}\\b([-a-zA-Z0-9@:%" +
	"._\\+~#?&//=]*)"

func main() {

	var (
		command string
		url     string
	)

	flag.StringVar(&command, "command", "geturl", "Command to execute")
	flag.StringVar(&url, "url", "google", "Url to get")

	flag.Parse()

	if command == "geturl" {
		fmt.Printf(geturl(url))
	}
}

func geturl(link string) string {

	url := getDomain(link)

	// Scan the cache to see if it already exists
	files, _ := ioutil.ReadDir("./cache")

	for _, file := range files {
		if url == file.Name() {
			res, _ := ioutil.ReadFile("./cache/" + url)
			return string(res)
		}
	}

	// Else, download the page, save it to cache, and return the html

	// check if it is a link, else treat it as a search query
	matches, _ := regexp.MatchString(urlRegex, link)
	if matches {
		res, err := http.Get(link) // get data

		return getRequest(res, err, url)
	}
	res, err := http.Get("https://www.google.com/search?q=" + link)
	return getRequest(res, err, url)
}

func getRequest(res *http.Response, err error, url string) string {
	if err != nil {
		return err.Error()
	}

	defer res.Body.Close()

	str, _ := ioutil.ReadAll(res.Body)

	// obtain html data
	cache, _ := os.Create("./cache/" + url) // create the file

	defer cache.Close()

	_, _ = cache.WriteString(string(str)) // save html to the file

	return string(str)
}

func getDomain(url string) string { // return a url that can be saved as file
	domain := strings.Replace(strings.Replace(url, "/", "_", -1), ":", ".", -1)
	return domain + ".html"
}
