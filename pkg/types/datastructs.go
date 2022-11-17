package types

import (
    "errors"
    "sync"
)

func _compareStrings(a string, b string) int {
    la := len(a)
    lb := len(b)
    var min int

    if la > lb {
        min = lb
    } else {
        min = la
    }
    var res int = 0

    for i := 0; i < min && res == 0; i++ {
        res = int(a[i]) - int(b[i])
    }

    if res == 0 {
        if la < lb {
            return -1
        } else if la == lb {
            return 0
        } else {
            return 1
        }
    }

    return res
}


type StringSet struct {
    values []string
    mut *sync.Mutex
}

func NewStringSet(values []string) *StringSet {
    
    var res *StringSet = &StringSet{
        values: make([]string, len(values)),
        mut: &sync.Mutex{},
    }

    for _, v := range values {
        res.AddWord(v)
    }

    return res
}

func (set *StringSet) elemAt(i int) string {
    return set.values[i]
}

func (set *StringSet) _binsearch(value string) (int, bool) {

    start := 0
    end := len(set.values)

    for start < end {
        middle := start + (end - start) / 2
        s := set.elemAt(middle)
        res := _compareStrings(value, s)
        if res == 0 {
            return middle, true
        } else if res < 0 {
            end = middle
        } else {
            start = middle + 1
        }
    }

    return start, false
}

func (set *StringSet) _insertAt(value string, pos int) {
    set.values = append(set.values, value)
    for i := len(set.values) - 1; i > pos; i-- {
        set.values[i] = set.elemAt(i - 1)
    }

    set.values[pos] = value
}

func (set *StringSet) AddWord(value string) bool {
    set.mut.Lock()
    defer set.mut.Unlock()
    pos, found := set._binsearch(value)
    if found {
        return false
    }

    set._insertAt(value, pos)
    return true
}

func (set *StringSet) ContainsWord(value string) bool {
    set.mut.Lock()
    _, found := set._binsearch(value)
    set.mut.Unlock()
    return found
}

func (set *StringSet) ToArray() []string {
    set.mut.Lock()
    defer set.mut.Unlock()
    dest := make([]string, len(set.values))
    copy(dest, set.values)
    return dest
}

func (set *StringSet) Length() int {
    set.mut.Lock()
    l := len(set.values)
    set.mut.Unlock()
    return l
}

const extensionSize int = 10

type Queue[T any] struct {
    enqueueIndex int
    dequeueIndex int
    values []T
    _arraylen int
    length int
    mut     *sync.Mutex
}

func NewQueue[T any]() *Queue[T] {
    queue := &Queue[T]{
        enqueueIndex: 0,
        dequeueIndex: 0,
        values: make([]T, extensionSize),
        _arraylen: extensionSize,
        length: 0,
        mut: &sync.Mutex{},
    }
    return queue
}

func (q *Queue[T]) _getElements() []T {
    res := make([]T, q.length)
    if q.dequeueIndex < q.enqueueIndex {
        for i := q.dequeueIndex; i < q.enqueueIndex;i++ {
            res[i] = q.values[i]
        }
    } else {
        index := 0
        for i := q.dequeueIndex; i < q._arraylen;i++ {
            res[index] = q.values[i]
            index++
        }

        for i := 0; i < q.enqueueIndex; i++ {
            res[index] = q.values[i]
            index++
        }
    }

    return res
}

func (q *Queue[T]) _shrink() {
    elements := q._getElements()
    q.values = elements
    q._arraylen = len(elements)
    q.length = q._arraylen
    q.dequeueIndex = 0
    q.enqueueIndex = 0
}

func (q *Queue[T]) _flatten() {
    elements := q._getElements()
    copy(q.values, elements)
    q.dequeueIndex = 0
    q.enqueueIndex = len(elements) % q._arraylen
}

func (q *Queue[T]) _extend(deltasize int) {
    freespace := q._arraylen - q.length
    required := deltasize - freespace
    if required <= 0 {
        return
    }

    q._flatten()
    q.values = append(q.values, make([]T, required)...)
    q.enqueueIndex = q._arraylen
    q._arraylen += required
}

func (q *Queue[T]) Enqueue(x T) {
    q.mut.Lock()
    defer q.mut.Unlock()
    if q.length == q._arraylen {
        q._extend(extensionSize)
    }
    q.length++
    q.values[q.enqueueIndex] = x
    q.enqueueIndex = (q.enqueueIndex + 1) % q._arraylen
}

func (q *Queue[T]) Dequeue() (T, error) {
    q.mut.Lock()
    defer q.mut.Unlock()
    var res T
    if q.length == 0 {
        return res, errors.New("Could not dequeue empty queue")
    }
    res = q.values[q.dequeueIndex]
    q.length--
    q.dequeueIndex = (q.dequeueIndex + 1) % q._arraylen
    return res, nil
}

func (q *Queue[T]) Length() int {
    q.mut.Lock()
    defer q.mut.Unlock()
    res := q.length
    return res
}
