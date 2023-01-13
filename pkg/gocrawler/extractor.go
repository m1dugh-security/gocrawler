package gocrawler

import (
    "regexp"
    "fmt"
    "strings"
    "io"
    "golang.org/x/net/html"
)

var (
    param string = `(\?[\-\w=\.~\;\[\]&]+)?`
    anchor = `(#[\w\.\-]*)?`
    localUrl = `(/[\w\.~=\-]+)+/?` + anchor
    rootUrl = `https?://([\w\-]+\.)+[a-z]{2,7}`
)

var rootUrlRegex = regexp.MustCompile(rootUrl)
var urlRegex = regexp.MustCompile(rootUrl + localUrl + param)
var emailRegex = regexp.MustCompile(`[\w\-\.]+@([\w\-]+\.)+[a-z]{2,7}`)

var tagExtractor = regexp.MustCompile(`("|')`+ localUrl + param + `("|')`)

func parseHTMLPage(reader io.Reader, url string) []string {

    rootUrl := rootUrlRegex.FindString(url)
    tokenizer := html.NewTokenizer(reader)

    res := make ([]string, 0, 0)

    tokenType := tokenizer.Next()
    for ;tokenType != html.ErrorToken; tokenType = tokenizer.Next() {
        if tokenType != html.StartTagToken {
            continue;
        }
        token := tokenizer.Token()

        for _, attr := range token.Attr {
            if attr.Key == "src" || attr.Key == "href" || attr.Key == "action" {
                var link string = attr.Val
                if strings.HasPrefix(link, "//") {
                    link = "https:" + link
                } else if strings.HasPrefix(link, "/") {
                    link = rootUrl + link
                } else if strings.HasPrefix(link, "#") {
                    link = url + link
                }
                res = append(res, link)
            }
        }
    }

    return res
}

func ExtractUrls(page string, url string) []string {

    rootUrl := rootUrlRegex.FindString(url)
    
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
                res[elements] = fmt.Sprintf("%s%s", rootUrl, path)
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
