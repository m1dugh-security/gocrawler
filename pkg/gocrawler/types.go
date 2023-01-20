package gocrawler

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
    "errors"
)

type ScopeEntry struct {
    Advanced bool
    Enabled bool        `json:"enabled"`
    Host string         `json:"host"`
    Protocol string     `json:"protocol"`
    File string         `json:"file"`
    URL string          `json:"url"`
    hostReg *regexp.Regexp
    protocolReg *regexp.Regexp
    fileReg *regexp.Regexp
    urlReg *regexp.Regexp
}

func (s *ScopeEntry) IsEnabled() bool {
    return s.Enabled
}

func (s *ScopeEntry) Setup(advanced bool) error {
    s.Advanced = advanced
    if advanced {
        if len(s.Protocol) > 0 {
            if strings.ToLower(s.Protocol) == "any" {
                s.Protocol = `^[a-z]{2,7}$`
            }
            reg, err := regexp.Compile(s.Protocol)
            if err != nil {
                return err
            }
            s.protocolReg = reg
        }
        if len(s.Host) > 0 {
            reg, err := regexp.Compile(s.Host)
            if err != nil {
                return err
            }
            s.hostReg = reg
        }

        if len(s.File) > 0 {
            reg, err := regexp.Compile(s.File)
            if err != nil {
                return err
            }
            s.fileReg = reg
        }
    } else {
        if len(s.URL) > 0 {
            reg, err := regexp.Compile(s.URL)
            if err != nil {
                return err
            }
            s.urlReg = reg
        }
    }

    return nil
}

func (s *ScopeEntry) IsValid(host, protocol, file string) bool {
    if s.Advanced {
        if s.hostReg != nil && !s.hostReg.MatchString(host) {
            return false
        }

        if s.protocolReg != nil && !s.protocolReg.MatchString(protocol) {
            return false
        }

        if s.fileReg != nil && !s.fileReg.MatchString(file) {
            return false
        }
    } else {
        url := fmt.Sprintf("%s://%s%s", protocol, host, file)
        if s.urlReg != nil && !s.urlReg.MatchString(url) {
            return false
        }
    }

    return true
}

type Scope struct {
    Advanced bool           `json:"advanced"`
    Exclude []*ScopeEntry   `json:"exclude"`
    Include []*ScopeEntry   `json:"include"`
}

func NewSimpleScope(include []string, exclude []string) (*Scope, error) {
    scope, err := NewScope(make([]*ScopeEntry, 0, len(include)), make([]*ScopeEntry, 0, len(exclude)), false)
    if err != nil {
        return nil, err
    }

    for _, s := range include {
        entry := &ScopeEntry{
            Enabled: true,
            URL: s,
        }
        scope.AddRule(entry, true)
    }

    for _, s := range exclude {
        entry := &ScopeEntry{
            Enabled: true,
            URL: s,
        }
        scope.AddRule(entry, false)
    }

    return scope, err
}

type burpScope struct {
    Scope *Scope `json:"scope"`
    Target *struct{
        Scope *Scope `json:"scope"`
    } `json:"target"`
}

func (s *Scope) setup() error {
    for _, entry := range s.Include {
        err := entry.Setup(s.Advanced)
        if err != nil {
            return err
        }
    }

    for _, entry := range s.Exclude {
        err := entry.Setup(s.Advanced)
        if err != nil {
            return err
        }
    }

    return nil
}

func (s *burpScope) getScope() (*Scope, error) {
    var res *Scope
    if s.Target != nil && s.Target.Scope != nil {
        res = s.Target.Scope
    } else {
        res = s.Scope
    }

    err := res.setup()
    if err != nil {
        return nil, err
    }

    return res, nil

}

func DeserializeScope(body []byte) (*Scope, error) {
    var scope burpScope
    err := json.Unmarshal(body, &scope)
    if err != nil {
        return nil, err
    }

    res, err := scope.getScope()
    if err != nil {
        return nil, errors.New(fmt.Sprintf("Could not deserialize scope: %s", err))
    }
    return res, nil
}

func NewScope(include []*ScopeEntry, exclude []*ScopeEntry, advanced bool) (*Scope, error) {
    res := &Scope{
        Exclude: exclude,
        Include: include,
        Advanced: advanced,
    }

    err := res.setup()
    if err != nil {
        return nil, err
    }
    return res, nil
}

func (s *Scope) AddRule(entry *ScopeEntry, in bool) {
    if in {
        s.Include = append(s.Include, entry)
    } else {
        s.Exclude = append(s.Exclude, entry)
    }
    entry.Setup(s.Advanced)
}

func (s *Scope) InScope(url string) bool {
    
    var protocol, host, file string
    splits := strings.SplitN(url, "://", 2)
    if len(splits) <= 1 {
        return false
    }
    protocol = splits[0]
    splits = strings.SplitN(splits[1], "/", 2)
    host = splits[0]
    file = "/"
    if len(splits) > 1 {
        file += splits[1]
    }

    for _, entry := range s.Exclude {
        if entry.IsEnabled() && entry.IsValid(host, protocol, file) {
            return false
        }
    }

    for _, entry := range s.Include {
        if entry.IsEnabled() && entry.IsValid(host, protocol, file) {
            return true
        }
    }

    return false
}
