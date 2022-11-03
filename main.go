package main

import "fmt"

func main() {
    q := CreateQueue()
    q.Enqueue("test")
    fmt.Println(q.Dequeue())
}
