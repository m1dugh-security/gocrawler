package types

import (
    "time"
    "sync"
    "regexp"
)

type RequestThrottler struct {
    MaxRequests int
    _requests int
    _lastFlush int64
    mut *sync.Mutex
}

func NewRequestThrottler(maxRequests int) *RequestThrottler {
    res := &RequestThrottler{
        maxRequests,
        0,
        time.Now().UnixMicro(),
        nil,
    }
    mut := &sync.Mutex{}
    res.mut = mut
    return res
}

func (r *RequestThrottler) AskRequest() {
    if r.MaxRequests < 0 {
        return
    }
    r.mut.Lock()
    defer r.mut.Unlock()

    timeStampMicro := time.Now().UnixMicro()
    delta := timeStampMicro - r._lastFlush
    if delta > 1000000 {
        r._requests = 1
        r._lastFlush = timeStampMicro
    } else {
        if r._requests < r.MaxRequests {
            r._requests++
        } else {
            for timeStampMicro - r._lastFlush < 1000000 {
                time.Sleep(time.Microsecond)
                timeStampMicro = time.Now().UnixMicro()
            }

            r._requests = 1
            r._lastFlush = timeStampMicro
        }
    }
}

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

type ThreadThrottler struct {
    MaxThreads uint
    threads uint
    mut *sync.Mutex
    wg *sync.WaitGroup
}

func NewThreadThrottler(maxThreads uint) *ThreadThrottler {
    return &ThreadThrottler{
        maxThreads,
        0,
        &sync.Mutex{},
        &sync.WaitGroup{},
    }
}

func (t* ThreadThrottler) Threads() uint {
    t.mut.Lock()
    value := t.threads
    t.mut.Unlock()
    return value
}

func (t *ThreadThrottler) RequestThread() {
    t.mut.Lock()
    defer t.mut.Unlock()
    if t.threads < t.MaxThreads {
        t.threads++;
        t.wg.Add(1)
        return
    }
    t.mut.Unlock()
    for true {
        t.mut.Lock()
        if t.threads < t.MaxThreads {
            t.wg.Add(1)
            t.threads++
            break
        }
        t.mut.Unlock()
    }
}

func (t *ThreadThrottler) Done() {
    t.wg.Done()
    t.mut.Lock()
    if t.threads > 0 {
        t.threads--
    }
    t.mut.Unlock()
}

func (t *ThreadThrottler) Wait() {
    t.wg.Wait()
}

