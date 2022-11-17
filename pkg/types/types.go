package types

import (
    "time"
    "sync"
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
    if t.threads < t.MaxThreads {
        t.threads++;
        t.wg.Add(1)
        t.mut.Unlock()
        return
    }
    t.mut.Unlock()
    for true {
        t.mut.Lock()
        if t.threads < t.MaxThreads {
            t.wg.Add(1)
            t.threads++
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

