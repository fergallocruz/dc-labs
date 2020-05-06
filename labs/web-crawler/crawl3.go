// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/todostreaming/gopl.io/ch5/links"
)

type httpPkg struct{}

var http httpPkg
var depth int

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

func crawl(url string, count int) []string {
	if count == 0 {
		fmt.Printf("<- Done with %v, depth 0.\n", url)
		return []string{""}
	}
	fmt.Println(count, url)
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(url)
	fmt.Printf("Found: %s %q\n", url)
	done := make(chan bool)
	worklist := make(chan []string)
	for i, u := range list {
		fmt.Printf("-> Crawling child %v/%v of %v : %v.\n", i, len(list), url, u)
		go func(u string) {
			worklist <- crawl(u, count-1)
			done <- true
		}(u)
	}
	for i, u := range list {
		fmt.Printf("<- [%v] %v/%v Waiting for child %v.\n", url, i, len(list), u)
		<-done
	}

	fmt.Printf("<- Done with %v\n", url)
	<-tokens // release the token
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	worklist := make(chan []string)
	depth := flag.Int("depth", depth, "depth limit")
	d := *depth
	// Parse the flags.
	flag.Parse()
	fmt.Println("  depth:", *depth)
	var n int // number of pending sends to worklist

	// Start with the command-line arguments.
	n++
	go func() { worklist <- os.Args[1:] }()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist //gets each link
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- crawl(link, d) //receives new links
				}(link)
			}
		}
	}
}

//!-
