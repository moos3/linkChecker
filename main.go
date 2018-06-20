package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/cheggaaa/pb.v1"
)

// Node - every link found stored as a node
type Node struct {
	link        string
	redirectURL string
	statusCode  int
}

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
	return lines, scanner.Err()
}

func urlCheck(url string) (Node, error) {

	var data Node
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	data.link = url
	data.statusCode = resp.StatusCode
	if resp.Request.URL.String() != url {
		data.redirectURL = resp.Request.URL.String()
	}

	return data, err

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

	checked := []Node{}
	bar := pb.StartNew(len(links))
	fmt.Println("Checking Urls")
	// crawl each website in input file one consecutively
	for i := 0; i < len(links); i++ {
		bar.Increment()
		l, e := urlCheck(links[i])
		if e != nil {
			fmt.Printf("Got Error: %s when fetching: %s", e, links[i])
		} else {
			checked = append(checked, l)
		}

	}
	fmt.Println("Reported 404 urls")
	for _, link := range checked {
		if link.statusCode == 404 {
			fmt.Printf("Url: %s, reported: %d\n", link.link, link.statusCode)
			//fmt.Println(link)
		}
	}

}
