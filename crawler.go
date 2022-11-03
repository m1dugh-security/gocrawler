package main

import (
    "net/http"
    "regexp"
    "io/ioutil"
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

type Crawler struct {
    Scope *Scope
    urls *Queue
    Discovered *StringSet
    callbacks []func(*http.Response, string)
}

func NewCrawler(scope *Scope) *Crawler {
    res := &Crawler{scope, CreateQueue(), NewStringSet(nil), nil}

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

func (cr *Crawler) AddCallback(f func(*http.Response, string)) {
    cr.callbacks = append(cr.callbacks, f)
}

func (cr *Crawler) Crawl(endpoints []string) {
    for _, v := range endpoints {
        cr.urls.Enqueue(v)
    }
    
    for cr.urls.Length > 0 {
        elem, err := cr.urls.Dequeue()
        if err != nil {
            break
        }

        url := elem.(string)
        if !cr.Scope.InScope(url) {
            continue
        }

        if cr.Discovered.AddWord(url) {
            urls, _ := cr.ExtractPageInfo(url)
            for _, u := range urls {
                cr.urls.Enqueue(u)
            }
        }
    }
}

