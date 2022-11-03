package main

// import "fmt"

func main() {

    scope := NewScope([]string{`([\w\-]\.)+com`})

    crawler := NewCrawler(scope)

    crawler.Crawl([]string{"https://www.google.com"})
}
