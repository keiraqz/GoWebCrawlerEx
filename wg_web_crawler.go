package main

import (
	"sync"
	"fmt"
)

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func WgCrawl(url string, depth int, fetcher Fetcher, cachedUrl SafeMap, wg *sync.WaitGroup) {
	// defer allows you to defer a function call to the end of the currently
	// executing function no matter HOW and WHERE it returns
	// good read: https://kylewbanks.com/blog/when-to-use-defer-in-go
	defer wg.Done()
	if depth <= 0 {
		// the following won't be necessary when you use defer
		//wg.Done()
		return
	}
	exist := cachedUrl.Check(url)
	if !exist {
		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			fmt.Println(err)
			// the following won't be necessary when you use defer
			//wg.Done()
			return
		}
		fmt.Printf("found: %s %q\n", url, body)
		// insert into cache so it won't crawl twice
		cachedUrl.Insert(url)
		for _, u := range urls {
			wg.Add(1)
			go WgCrawl(u, depth-1, fetcher, cachedUrl, wg)
		}
	}
	// the following won't be necessary when you use defer
	//wg.Done()
	return
}

//func main() {
//	var wg sync.WaitGroup // make sure to pass its pointers rather than just value
//	cachedUrl := SafeMap{set: make(map[string]bool)}
//	wg.Add(1)
//	go WgCrawl("https://golang.org/", 4, fetcher, cachedUrl, &wg)
//	wg.Wait()
//}