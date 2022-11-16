package crawler

import (
    "net/http"
    "regexp"
    "io/ioutil"
    "sync"
    "github.com/m1dugh/crawler/pkg/types"
    "github.com/m1dugh/crawler/pkg/utils"
)

var rootUrlRegex = regexp.MustCompile(`https?://([\w\-]+\.)[a-z]{2,7}`)


type Config struct {
    MaxThreads uint
}

func DefaultConfig() *Config {
    return &Config{10}
}

type Callback func(*http.Response, string)

type Crawler struct {
    Scope *types.Scope
    urls *types.Queue
    Discovered *types.StringSet
    callbacks []Callback
    Config *Config
}

func New(scope *types.Scope, config *Config) *Crawler {
    if config == nil {
        config = DefaultConfig()
    }

    res := &Crawler{scope,
    types.CreateQueue(),
    types.NewStringSet(nil),
    nil,
    config}

    return res
}

func (cr *Crawler) runCallbacks(resp *http.Response, body string) {
    for _, f := range cr.callbacks {
        f(resp, body)
    }
}

func (cr *Crawler) ExtractPageInfo(url string) ([]string, []string) {
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
        urls, _ := cr.ExtractPageInfo(url)
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

