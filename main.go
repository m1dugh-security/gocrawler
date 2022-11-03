package main

import (
    "net/http"
    "fmt"
)

func main() {

    scope := NewScope([]string{`([\w\-]\.)+com`})

    crawler := NewCrawler(scope)

    crawler.AddCallback(func(res *http.Response, _ string) {
        fmt.Println(res.Request.URL)
    })

    crawler.Crawl([]string{"https://www.google.com"})
}
