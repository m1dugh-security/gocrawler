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
    "errors"
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

type headersFlags struct {
    header http.Header
}

func (h *headersFlags) String() string {
    return "TODO: provide implementation for headersFlag.String in main package"
}

func (h *headersFlags) Set(value string) error {
    
    var err error = errors.New("headers should be provided as 'Key: value'")   
    
    var key, headerValue string
    var spaceGroupCount int = 0
    var i int
    for i = 0; i < len(value) && value[i] != ':'; i++ {
        c := value[i]
        if c == ' ' {
            if i == 0 || value[i - 1] != ' ' {
                spaceGroupCount++
            }
        } else if (c <= 'Z' && c >= 'A') || (c <= 'z' && c >= 'a') || (c <= '9' && c >= '0') || c == '-' {
            if spaceGroupCount > 1 {
                return err
            }
            key += string(c)
        } else {
            return err
        }
    }

    if i == 0 || i + 1 >= len(value) {
        return err
    }

    spaceGroupCount = 0

    for i++; i < len(value); i++ {
        c := value[i]
        if c == ' ' {
            if value[i - 1] != ' ' {
                spaceGroupCount++
            }
        } else if spaceGroupCount <= 1 {
            headerValue += string(c)
        } else {
            return err
        }
    }
    h.header.Set(key, value)

    return nil
}

func  (h headersFlags) ToHTTPHeader() http.Header {
    return h.header
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

    var headers headersFlags = headersFlags{make(map[string][]string)}
    flag.Var(&headers, "H", "additional headers to add to the crawler")

    flag.Parse()

    var urls []string
    var err error

    if !isInPipe() {
        if len(inputFile) == 0 {
            log.Fatal("Missing urls file or stdin input")
        }
        urls, err = DeserializeUrls(inputFile)
        if err != nil {
            log.Fatal(err)
        }
    } else {
        urls = readStdin(os.Stdin)
    }


    if len(scopeFile) == 0 {
        log.Fatal("Missing -scope")
    }

    body, err := os.ReadFile(scopeFile)
    if err != nil {
        log.Fatal(err)
    }
    scope, err := gocrawler.DeserializeScope(body)
    if err != nil {
        log.Fatal(err)
    }

    config := &gocrawler.Config{
        MaxThreads: threadCount,
        MaxRequests: requestThrottle,
        Headers: headers.ToHTTPHeader(),
    }
    cr:= gocrawler.New(scope, config)

    cr.AddCallback(func(res *http.Response, _ string) {
        fmt.Println(res.Request.URL)
    })

    cr.Crawl(urls)
}
