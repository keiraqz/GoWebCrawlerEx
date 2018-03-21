package main

import (
	"fmt"
	"sync"
	"flag"
)

func main() {
	algorithm := flag.String("algorithm", "channel", "'channel' (default) or 'wg'")
	flag.Parse()
	cachedUrl := SafeMap{set: make(map[string]bool)}

	switch *algorithm {
	case "channel":
		ChanCrawl("https://golang.org/", 4, fetcher, cachedUrl)
	case "wg":
		var wg sync.WaitGroup // make sure to pass its pointers rather than just value
		wg.Add(1)
		go WgCrawl("https://golang.org/", 4, fetcher, cachedUrl, &wg)
		wg.Wait()
	default:
		fmt.Println("Unknown algorithm. Valid values are: 'channel' or 'sync'")
	}
}

type SafeMap struct {
	set map[string]bool
	sync.Mutex
}

func (cachedUrl *SafeMap) Check(url string) bool {
	cachedUrl.Lock()
	defer cachedUrl.Unlock()
	_, ok := cachedUrl.set[url]
	if ok {
		fmt.Printf("exist URL: %s.\n", url)
	}
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
