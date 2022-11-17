package gocrawler

import (
    "net/http"
    "regexp"
    "io/ioutil"
    "sync"
    "github.com/m1dugh/gocrawler/pkg/types"
    "github.com/m1dugh/gocrawler/pkg/utils"
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

type Crawler struct {
    Scope *types.Scope
    urls *types.Queue
    Discovered *types.StringSet
    callbacks []Callback
    Config *Config
    throttler *types.RequestThrottler
}

func New(scope *types.Scope, config *Config) *Crawler {
    if config == nil {
        config = DefaultConfig()
    }

    res := &Crawler{
        Scope: scope,
        urls: types.CreateQueue(),
        Discovered: types.NewStringSet(nil),
        callbacks: nil,
        Config: config,
        throttler: types.NewRequestThrottler(config.MaxRequests),
    }

    return res
}

func (cr *Crawler) runCallbacks(resp *http.Response, body string) {
    for _, f := range cr.callbacks {
        f(resp, body)
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
        resp.Body.Close()
        return nil, nil
    }

    var content string = string(body)
    resp.Body.Close()
    cr.runCallbacks(resp, content)

    return utils.ExtractUrls(content, rootUrl), utils.ExtractEmails(content)
}

func (cr *Crawler) AddCallback(f Callback) {
    cr.callbacks = append(cr.callbacks, f)
}

const maxWorkers = 10

func (cr *Crawler) crawlPage(count *uint, mut *sync.Mutex, url string) {

    if !cr.Scope.InScope(url) {
        mut.Lock()
        *count = *count - 1
        mut.Unlock()
        return
    }

    mut.Lock()
    added := cr.Discovered.AddWord(url)
    mut.Unlock()

    if added {
        cr.throttler.AskRequest()
        urls, _ := cr.extractPageInfo(url)
        mut.Lock()
        for _, u := range urls {
            cr.urls.Enqueue(u)
        }
        mut.Unlock()
    }

    mut.Lock()
    *count = *count - 1
    mut.Unlock()
}

func (cr *Crawler) Crawl(endpoints []string) {
    for _, v := range endpoints {
        cr.urls.Enqueue(v)
    }


    var workerCount uint = 0
    var mut sync.Mutex

    for cr.urls.Length > 0 || workerCount > 0{
        for workerCount < cr.Config.MaxThreads {
            mut.Lock()
            elem, err := cr.urls.Dequeue()
            mut.Unlock()
            if err != nil {
                break
            }

            url := elem.(string)
            mut.Lock()
            workerCount++
            mut.Unlock()
            go cr.crawlPage(&workerCount, &mut, url)
        }
    }
}

