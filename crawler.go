package main

import (
    "net/http"
    "regexp"
    "io/ioutil"
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

