package utils

import (
    "io/ioutil"
    "os"
    "log"
    "encoding/json"
    "github.com/m1dugh/gocrawler/pkg/types"
)

type ScopeRepr struct {
    In []string `json:"include"`
    Ex []string `json:"exclude"`
}

func DeserializeScope(file string) *types.Scope {
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
    return types.NewScope(s.In, s.Ex)
}
