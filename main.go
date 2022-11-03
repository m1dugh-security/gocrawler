package main

import "fmt"

func main() {

    set := NewStringSet(nil)

    scope := NewScope([]string{`([\w\-]\.)+com`})

    q := CreateQueue()
    q.Enqueue("https://www.google.com/search?q=test")

    for q.Length > 0 {
        elem, err := q.Dequeue()
        if err != nil {
            break
        }

        url := elem.(string)
        if !scope.InScope(url) {
            continue
        }

        if set.AddWord(url) {
            fmt.Println(url)
            urls, _ := ExtractPageInfo(url)
            for _, u := range urls {
                q.Enqueue(u)
            }
        }
    }

}
