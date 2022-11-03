package main

import (
    "errors"
)

const extensionSize int = 10

type Queue struct {
    enqueueIndex int
    dequeueIndex int
    values []interface{}
    _arraylen int
    Length int
}

func CreateQueue() *Queue {
    queue := &Queue{0, 0, make([]interface{}, extensionSize), extensionSize, 0}
    return queue
}

func (q *Queue) _getElements() []interface{} {
    res := make([]interface{}, q.Length)
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

func (q *Queue) _shrink() {
    elements := q._getElements()
    q.values = elements
    q._arraylen = len(elements)
    q.Length = q._arraylen
    q.dequeueIndex = 0
    q.enqueueIndex = 0
}

func (q *Queue) _flatten() {
    elements := q._getElements()
    copy(q.values, elements)
    q.dequeueIndex = 0
    q.enqueueIndex = len(elements) % q._arraylen
}

func (q *Queue) _extend(deltasize int) {
    freespace := q._arraylen - q.Length
    required := deltasize - freespace
    if required <= 0 {
        return
    }

    q._flatten()
    q.values = append(q.values, make([]interface{}, required)...)
    q.enqueueIndex = q._arraylen
    q._arraylen += required
}

func (q *Queue) Enqueue(x interface{}) {
    if q.Length == q._arraylen {
        q._extend(extensionSize)
    }
    q.Length++
    q.values[q.enqueueIndex] = x
    q.enqueueIndex = (q.enqueueIndex + 1) % q._arraylen
}

func (q *Queue) Dequeue() (interface{}, error) {
    if q.Length == 0 {
        return nil, errors.New("Could not dequeue empty queue")
    }
    res := q.values[q.dequeueIndex]
    q.Length--
    q.dequeueIndex = (q.dequeueIndex + 1) % q._arraylen
    return res, nil
}
