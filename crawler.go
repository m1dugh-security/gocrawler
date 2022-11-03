package main

import (
    "net/http"
    "regexp"
    "io/ioutil"
    "sync"
)

var rootUrlRegex = regexp.MustCompile(`https?://([\w\-]+\.)[a-z]{2,7}`)


type Scope []*regexp.Regexp

func NewScope(regexes []string) *Scope {
    res := &Scope{}
    *res = make([]*regexp.Regexp, 0, len(regexes))
    for _, exp := range regexes {
        res.AddRule(exp)
    }

    return res
}

func (s *Scope) AddRule(v string) {
    
    re, err := regexp.Compile(v)
    if err != nil {
        // log.Warn(fmt.Sprintf(`could not compile "%s"`), v)
    }
    *s = append(*s, re)
}

func (s *Scope) InScope(url string) bool {
    for _, re := range *s {
        if re.MatchString(url) {
            return true
        }
    }

    return false
}

type Config struct {
    MaxThreads uint
}

func DefaultConfig() *Config {
    return &Config{10}
}

type Callback func(*http.Response, string)

type Crawler struct {
    Scope *Scope
    urls *Queue
    Discovered *StringSet
    callbacks []Callback
    Config *Config
}

func NewCrawler(scope *Scope, config *Config) *Crawler {
    if config == nil {
        config = DefaultConfig()
    }

    res := &Crawler{scope,
    CreateQueue(),
    NewStringSet(nil),
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

    return ExtractUrls(content, rootUrl), ExtractEmails(content)
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

