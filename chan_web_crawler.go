package main

import (
	"fmt"
)

func _crawl(url string, depth int, fetcher Fetcher, cachedUrl SafeMap, chanResults chan links) {
	body, urls, err := fetcher.Fetch(url)
	// insert into cache so it won't crawl twice
	cachedUrl.Insert(url)
	// create a links
	urlResult := links{urls:urls, body:body, depth:depth-1, err:err}
	chanResults <- urlResult
	return
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func ChanCrawl(url string, depth int, fetcher Fetcher, cachedUrl SafeMap) {
	chanResults := make(chan links, 0)
	go _crawl(url, depth, fetcher, cachedUrl, chanResults)
	for counter := 1; counter > 0; counter-- {
		result := <- chanResults
		depth := result.depth
		body := result.body
		err := result.err
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("found: %s %q\n", url, body)
		if depth > 0 {
			for _, url := range result.urls {
				exist := cachedUrl.Check(url)
				if !exist {
					counter++
					go _crawl(url, depth, fetcher, cachedUrl, chanResults)
				}
			}
		}
	}
	close(chanResults)
}

type links struct {
	urls []string
	body string
	depth int
	err error
}