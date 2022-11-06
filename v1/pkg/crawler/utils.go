package crawler

import (
    "io/ioutil"
    "os"
    "log"
    "encoding/json"
)

type ScopeRepr struct {
    In []string `json:"include"`
    Ex []string `json:"exclude"`
}

func DeserializeScope(file string) *Scope {
    f, err := os.Open(file)

    if err != nil {
        log.Fatal("could not read scope")
    }

    defer f.Close()

    bytes, err := ioutil.ReadAll(f)


    if err != nil {
        log.Fatal(err)
    }

    var s ScopeRepr
    err = json.Unmarshal(bytes, &s)
    if err != nil {
        log.Fatal(err)
    }
    return NewScope(s.In, s.Ex)
}
