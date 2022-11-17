package main

import (
    "net/http"
    "fmt"
    "flag"
    "log"
    "os"
    "bufio"
    "io"
    "github.com/m1dugh/gocrawler/pkg/gocrawler"
    "github.com/m1dugh/gocrawler/pkg/utils"
)

func isInPipe() bool {
    fileinfo, _ := os.Stdin.Stat()
    return fileinfo.Mode() & os.ModeCharDevice == 0
}

func readStdin(r io.Reader) []string {
    scanner := bufio.NewScanner(bufio.NewReader(r))
    var res []string
    for scanner.Scan() {
        res = append(res, scanner.Text())
    }

    if err := scanner.Err();err != nil {
        log.Fatal("failed reading from stdin")
    }

    return res
}

func main() {

    var scopeFile string
    flag.StringVar(&scopeFile, "scope", "", "the file containing regexes for the scope")
    var threadCount uint
    flag.UintVar(&threadCount, "threads", 10, "Number of max concurrent threads")
    var requestThrottle int
    flag.IntVar(&requestThrottle, "requests", -1, "Max requests per second")
    var inputFile string
    flag.StringVar(&inputFile, "urls", "", "The file containing all urls")

    flag.Parse()

    var urls []string
    var err error

    if !isInPipe() {
        urls, err = utils.DeserializeUrls(inputFile)
        if err != nil {
            log.Fatal(err)
        }
    } else {
        urls = readStdin(os.Stdin)
    }

    scope, err := utils.DeserializeScope(scopeFile)
    if err != nil {
        log.Fatal(err)
    }

    config := &gocrawler.Config{
        MaxThreads: threadCount,
        MaxRequests: requestThrottle,
    }
    cr:= gocrawler.New(scope, config)

    cr.AddCallback(func(res *http.Response, _ string) {
        fmt.Println(res.Request.URL)
    })

    cr.Crawl(urls)
}
