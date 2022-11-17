package gocrawler

import (
    "regexp"
)

type Scope struct {
    Exclude []*regexp.Regexp
    Include []*regexp.Regexp
}

func NewScope(include []string, exclude []string) *Scope {
    res := &Scope{}
    res.Include = make([]*regexp.Regexp, 0, len(include))
    for _, exp := range include {
        res.AddRule(exp, true)
    }

    res.Exclude = make([]*regexp.Regexp, 0, len(exclude))
    for _, exp := range exclude {
        res.AddRule(exp, false)
    }

    return res
}

func (s *Scope) AddRule(v string, in bool) {
    
    re, err := regexp.Compile(v)
    if err != nil {
        return
    }
    if in {
        s.Include = append(s.Include, re)
    } else {
        s.Exclude = append(s.Exclude, re)
    }
}

func (s *Scope) InScope(url string) bool {
    valid := false
    for _, re := range s.Include {
        if re.MatchString(url) {
            valid = true
            break
        }
    }

    for i := 0; i < len(s.Exclude) && valid; i++ {
        if s.Exclude[i].MatchString(url) {
            valid = false
        }
    }

    return valid
}
