package utils

import (
    "regexp"
    "fmt"
    "strings"
    "html"
)

var (
    param string = `(\?[\-\w=\.~\;\[\]&]+)?`
    localUrl = `(/[\w\.~=\-]+)+/?`
)

var urlRegex = regexp.MustCompile(`https?://([\w\-]+\.)+[a-z]{2,7}` + localUrl + param)
var emailRegex = regexp.MustCompile(`[\w\-\.]+@([\w\-]+\.)+[a-z]{2,7}`)

var tagExtractor = regexp.MustCompile(`(src|href)="`+ localUrl + param + `"`)

func ExtractUrls(page string, root_url string) []string {
    
    res := urlRegex.FindAllString(page, -1)
    startIndex := len(res)
    elements := startIndex
    res = append(res, tagExtractor.FindAllString(page, -1)...)

    for ; startIndex < len(res); startIndex++ {
        path := res[startIndex]
        i := strings.Index(path, "\"")
        if i > 0 {
            path = path[i+1:]
            i = strings.Index(path, "\"")
            if i > 0 {
                path = path[:i]
                res[elements] = fmt.Sprintf("%s%s", root_url, path)
                elements++
            }
        }
    }

    for i := 0; i < elements; i++ {
        res[i] = html.UnescapeString(res[i])
    }

    return res[:elements]
}

func ExtractEmails(page string) []string {
    return emailRegex.FindAllString(page, -1)
}
