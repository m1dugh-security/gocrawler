package main

import (
    "net/http"
    "fmt"
    "flag"
    "log"
    "os"
    "bufio"
    "io"
    "github.com/m1dugh/crawler/pkg/crawler"
    "github.com/m1dugh/crawler/pkg/utils"
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

    var urls []string

    if !isInPipe() {
        log.Fatal("expected input from stdin")
    }

    urls = readStdin(os.Stdin)

    var scopeFile string
    flag.StringVar(&scopeFile, "scope", "", "the file containing regexes for the scope")
    var threadCount uint
    flag.UintVar(&threadCount, "threads", 10, "Number of max concurrent threads")

    flag.Parse()

    if len(scopeFile) == 0 {
        log.Fatal("missing required argument -scope")
    }

    scope := utils.DeserializeScope(scopeFile)

    cr:= crawler.New(scope, nil)

    cr.Config.MaxThreads = threadCount

    cr.AddCallback(func(res *http.Response, _ string) {
        fmt.Println(res.Request.URL)
    })

    cr.Crawl(urls)
}
