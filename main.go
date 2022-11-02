package main

import "fmt"

func main() {
    t := CreatePrefixTree()
    t.AddWord("test")
    t.AddWord("abcd")
    t.AddWord("tet")
    t.AddWord("https://www.google.com")
    fmt.Println(t.SearchWord("https://www.google.com"))
    fmt.Println(t.SearchWord("https://"))
    fmt.Println(t.SearchWord("https://www.google.com/"))
    for _,v := range t.ListWords() {
        fmt.Println(v)
    }
}
