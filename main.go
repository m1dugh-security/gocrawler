package main

import (
    "net/http"
    "fmt"
)

func main() {

    scope := NewScope([]string{`([\w\-]\.)+google\.com`})

    crawler := NewCrawler(scope, nil)

    crawler.AddCallback(func(res *http.Response, _ string) {
        fmt.Println(res.Request.URL)
    })

    crawler.Crawl([]string{"https://www.google.com"})
}
