package main

import (
    "fmt"
    "os"
    "bufio"
    "errors"
)


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
