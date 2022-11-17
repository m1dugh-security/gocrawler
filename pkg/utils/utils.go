package utils

import (
    "fmt"
    "io/ioutil"
    "os"
    "bufio"
    "encoding/json"
    "errors"
    "github.com/m1dugh/gocrawler/pkg/types"
)

type ScopeRepr struct {
    In []string `json:"include"`
    Ex []string `json:"exclude"`
}

func DeserializeScope(path string) (*types.Scope, error) {
    f, err := os.Open(path)

    if err != nil {
        return nil, errors.New(fmt.Sprintf("Could not open file %s", path))
    }

    defer f.Close()

    bytes, err := ioutil.ReadAll(f)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("Could not read file %s", path))
    }

    var s ScopeRepr
    err = json.Unmarshal(bytes, &s)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("Could not parse json file %s", path))
    }
    return types.NewScope(s.In, s.Ex), nil
}

func DeserializeUrls(path string) ([]string, error) {
    f, err := os.Open(path)

    if err != nil {
        return nil, errors.New(fmt.Sprintf("Could not open file %s", path))
    }

    defer f.Close()

    scanner := bufio.NewScanner(f)
    var res []string
    for scanner.Scan() {
        res = append(res, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        return nil, errors.New(fmt.Sprintf("Error while scanning file %s", path))
    }

    return res, nil
}
