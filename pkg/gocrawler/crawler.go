package gocrawler

import (
    "net/http"
    "github.com/m1dugh/gocrawler/pkg/types"
    "bytes"
    "io"
    "strings"
)

type Config struct {
    MaxThreads uint
    MaxRequests int
}

func DefaultConfig() *Config {
    return &Config{
        MaxThreads: 10,
        MaxRequests: -1,
    }
}

type Callback func(*http.Response, string)

type cb_holder struct {
    Callback
    id int
}

type Crawler struct {
    Scope *Scope
    urls *types.Queue[string]
    Discovered *types.StringSet
    callbacks []cb_holder
    Config *Config
    throttler *types.RequestThrottler
}

func New(scope *Scope, config *Config) *Crawler {
    if config == nil {
        config = DefaultConfig()
    }

    res := &Crawler{
        Scope: scope,
        urls: types.NewQueue[string](),
        Discovered: types.NewStringSet(nil),
        callbacks: nil,
        Config: config,
        throttler: types.NewRequestThrottler(config.MaxRequests),
    }

    return res
}

func (cr *Crawler) runCallbacks(resp *http.Response, body string) {
    for _, holder := range cr.callbacks {
        holder.Callback(resp, body)
    }
}

func (cr *Crawler) extractPageInfo(url string) []string {
    resp, err := http.Get(url)
    if err != nil {
        return nil
    }

    var res []string

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil
    }
    resp.Body.Close()

    content := string(body)

    for _, contentType := range resp.Header["Content-Type"] {
        if strings.Contains(contentType, "text/html") {
            reader := bytes.NewReader(body)
            res = parseHTMLPage(reader, url)
        } else {
            res = ExtractUrls(content, url)
        }
    }

    cr.runCallbacks(resp, content)

    return res
}

/// Adds a callback to the crawler
/// Returns a handler to remove the callback
func (cr *Crawler) AddCallback(f Callback) int {
    var id int = 0
    if len(cr.callbacks) > 0 {
        id = cr.callbacks[len(cr.callbacks) - 1].id + 1
    }
    cr.callbacks = append(cr.callbacks, cb_holder{
        id: id,
        Callback: f,
    })
    return id
}

func (cr *Crawler) RemoveCallback(handler int) {
    var pos int = 0
    for pos < len(cr.callbacks) && cr.callbacks[pos].id != handler {
        pos++
    }
    if pos < len(cr.callbacks) {
        res := cr.callbacks[:pos]
        cr.callbacks = append(res, cr.callbacks[pos + 1:]...)
    }
}

func (cr *Crawler) crawlPage(threads *types.ThreadThrottler, url string) {
    defer threads.Done()
    added := cr.Discovered.AddWord(url)
    if added {
        cr.throttler.AskRequest()
        urls := cr.extractPageInfo(url)
        for _, u := range urls {
            if cr.Scope.InScope(u) {
                cr.urls.Enqueue(u)
            }
        }
    }
}

func (cr *Crawler) Crawl(endpoints []string) {
    for _, v := range endpoints {
        cr.urls.Enqueue(v)
    }

    threads := types.NewThreadThrottler(cr.Config.MaxThreads)

    for cr.urls.Length() > 0 {
        for cr.urls.Length() > 0 || threads.Threads() > 0 {
            url, err := cr.urls.Dequeue()
            if err != nil {
                continue
            }
            threads.RequestThread()
            go cr.crawlPage(threads, url)
        }
        threads.Wait()
    }
}

