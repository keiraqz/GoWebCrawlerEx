# A Tour of Go Exercise: Web Crawler

https://tour.golang.org/concurrency/10

To run waiting group: `go run *.go -algorithm=wg`

To run channel: `go run *.go -algorithm=channel`

### Waiting Group
    
A sync.WaitGroup waits for a group of goroutines to finish. http://yourbasic.org/golang/wait-for-goroutines-waitgroup/

### Channel

A channel is a mechanism for goroutines to synchronize execution and communicate by passing values. Receivers always block until there is data to receive.

