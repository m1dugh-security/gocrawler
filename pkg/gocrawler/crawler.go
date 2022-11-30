package gocrawler

import (
    "net/http"
    "regexp"
    "io/ioutil"
    "github.com/m1dugh/gocrawler/pkg/types"
)

var rootUrlRegex = regexp.MustCompile(`https?://([\w\-]+\.)[a-z]{2,7}`)


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

func (cr *Crawler) extractPageInfo(url string) ([]string, []string) {
    rootUrl := rootUrlRegex.FindString(url)
    resp, err := http.Get(url)
    if err != nil {
        return nil, nil
    }


    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, nil
    }
    var content string = string(body)
    resp.Body.Close()
    cr.runCallbacks(resp, content)

    // return ExtractUrls(content, rootUrl), ExtractEmails(content)
    return ExtractUrls(content, rootUrl), nil
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
    if !cr.Scope.InScope(url) {
        return
    }

    added := cr.Discovered.AddWord(url)
    if added {
        cr.throttler.AskRequest()
        urls, _ := cr.extractPageInfo(url)
        for _, u := range urls {
            cr.urls.Enqueue(u)
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

