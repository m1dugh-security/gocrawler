package main

import (
    "bufio"
    "os"
    "log"
)

func DeserializeScope(file string) *Scope {
    f, err := os.Open(file)

    if err != nil {
        log.Fatal("could not read scope")
    }

    defer f.Close()

    scanner := bufio.NewScanner(f)
    var res []string = nil

    for scanner.Scan() {
        res = append(res, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal("Could not properly read scope file.")
    }

    return NewScope(res)
}
