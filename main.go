package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/cheggaaa/pb.v1"
)

// Node - every link found stored as a node
type Node struct {
	link        string
	redirectURL string
	statusCode  int
}

// readLines -
// This reads a file of urls and will remove ones that are dupicates.
// This func is very dirty and could be smartter.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	encountered := map[string]bool{}
	result := []string{}
	for v := range lines {
		if encountered[lines[v]] == true {
			// do not add dupe
		} else {
			encountered[lines[v]] = true
			// Append to result slice
			result = append(result, lines[v])
		}
	}
	return result, scanner.Err()
}

// urlCheck -
// This function is for URL checking and will return a single Node
func urlCheck(url string) (Node, error) {

	var data Node

	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Print(err)
	}
	data.link = url
	data.statusCode = resp.StatusCode
	if resp.Request.URL.String() != url {
		data.redirectURL = resp.Request.URL.String()
	}

	return data, err

}

// runChecker - This just loops over a slice of strings
func runChecker(l []string) []Node {
	checked := []Node{}
	fmt.Println("Checking Urls")
	bar := pb.StartNew(len(l))
	// crawl each website in input file one consecutively
	for i := 0; i < len(l); i++ {
		bar.Increment()
		n, e := urlCheck(l[i])
		if e != nil {
			fmt.Printf("Got Error: %s when fetching: %s", e, l[i])
		} else {
			checked = append(checked, n)
		}

	}

	return checked
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s file path\n", os.Args[0])
		os.Exit(1)
	}
	links, err := readLines(os.Args[1])
	if err != nil {
		fmt.Println("Error reading file: ", err)
	}

	checked := runChecker(links)

	fmt.Println("Reported 404 urls")
	for _, link := range checked {
		if link.statusCode == 404 {
			fmt.Printf("Url: %s, reported: %d\n", link.link, link.statusCode)
		}
	}

}
