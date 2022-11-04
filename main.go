package main

import (
    "net/http"
    "fmt"
    "flag"
    "log"
)

func main() {

    var scopeFile string
    flag.StringVar(&scopeFile, "scope", "", "the file containing regexes for the scope")
    var threadCount uint
    flag.UintVar(&threadCount, "threads", 10, "Number of max concurrent threads")

    flag.Parse()

    if len(scopeFile) == 0 {
        log.Fatal("missing required argument -scope")
    }

    scope := DeserializeScope(scopeFile)

    crawler := NewCrawler(scope, nil)

    crawler.Config.MaxThreads = threadCount

    crawler.AddCallback(func(res *http.Response, _ string) {
        fmt.Println(res.Request.URL)
    })

    crawler.Crawl([]string{"https://www.google.com"})
}
