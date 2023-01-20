package types

import (
    "time"
    "sync"
)

type RequestThrottler struct {
    MaxRequests int
    requests int
    lastFlush int64
    mut *sync.Mutex
}

func NewRequestThrottler(maxRequests int) *RequestThrottler {
    res := &RequestThrottler{
        MaxRequests: maxRequests,
        requests: 0,
        lastFlush: time.Now().UnixMicro(),
        mut: &sync.Mutex{},
    }
    return res
}

func (r *RequestThrottler) AskRequest() {
    if r.MaxRequests < 0 {
        return
    }
    r.mut.Lock()
    defer r.mut.Unlock()

    timeStampMicro := time.Now().UnixMicro()
    delta := timeStampMicro - r.lastFlush
    if delta > 1000000 {
        r.requests = 1
        r.lastFlush = timeStampMicro
    } else {
        if r.requests < r.MaxRequests {
            r.requests++
        } else {
            for timeStampMicro - r.lastFlush < 1000000 {
                time.Sleep(time.Microsecond)
                timeStampMicro = time.Now().UnixMicro()
            }

            r.requests = 1
            r.lastFlush = timeStampMicro
        }
    }
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
    defer t.mut.Unlock()
    return t.threads
}

func (t *ThreadThrottler) RequestThread() {
    t.mut.Lock()
    // Thread count is incremented before being started avoiding it to be
    // fetched with the wrong value while a thread is being requested.
    t.threads++
    if t.threads <= t.MaxThreads {
        t.wg.Add(1)
        t.mut.Unlock()
        return
    }
    t.mut.Unlock()
    for true {
        t.mut.Lock()
        if t.threads <= t.MaxThreads {
            t.wg.Add(1)
            t.mut.Unlock()
            break
        }
        t.mut.Unlock()
    }
}

func (t *ThreadThrottler) Done() {
    t.mut.Lock()
    if t.threads > 0 {
        t.threads--
    }
    t.mut.Unlock()
    t.wg.Done()
}

func (t *ThreadThrottler) Wait() {
    t.wg.Wait()
}

