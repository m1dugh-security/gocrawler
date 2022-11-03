package main

import (
    "net/http"
    "regexp"
    "io/ioutil"
    "fmt"
)

var rootUrlRegex = regexp.MustCompile(`https?://([\w\-]+\.)[a-z]{2,7}`)

func ExtractPageInfo(url string) ([]string, []string) {
    rootUrl := rootUrlRegex.FindString(url)
    resp, err := http.Get(url)
    if err != nil {
        return nil, nil
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, nil
    }

    var content string = string(body)
    
    return ExtractUrls(content, rootUrl), ExtractEmails(content)
}

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
}

func NewCrawler(scope *Scope) *Crawler {
    res := &Crawler{scope, CreateQueue(), NewStringSet(nil)}

    return res
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
            fmt.Println(url)
            urls, _ := ExtractPageInfo(url)
            for _, u := range urls {
                cr.urls.Enqueue(u)
            }
        }
    }
}

