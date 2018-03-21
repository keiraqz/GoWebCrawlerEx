package main

import (
	"sync"
	"fmt"
)

type SafeMap struct {
	set map[string]bool
	sync.Mutex
}

func (cachedUrl *SafeMap) Check(url string) bool {
	cachedUrl.Lock()
	_, ok := cachedUrl.set[url]
	if ok {
		fmt.Printf("exist URL: %s.\n", url)
	}
	defer cachedUrl.Unlock()
	return ok
}

func (cachedUrl *SafeMap) Insert(url string) {
	cachedUrl.Lock()
	cachedUrl.set[url] = true
	cachedUrl.Unlock()
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, cachedUrl SafeMap, wg *sync.WaitGroup) {
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
			go Crawl(u, depth-1, fetcher, cachedUrl, wg)
		}
	}
	// the following won't be necessary when you use defer
	//wg.Done()
	return
}

func main() {
	var wg sync.WaitGroup // make sure to pass its pointers rather than just value
	cachedUrl := SafeMap{set: make(map[string]bool)}
	wg.Add(1)
	go Crawl("https://golang.org/", 4, fetcher, cachedUrl, &wg)
	wg.Wait()
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
