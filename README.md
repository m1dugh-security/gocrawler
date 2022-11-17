# Gocrawler

## Features

- Fast crawling over given inputs
- Stdin/Stdout support tooling integration
- Integrated request throttler
- Support for advanced scopes

## Installation

### Go Module

To install GoCrawler as a GO module:
```shell
$ go install -v github.com/m1dugh/gocrawler@latest
```

### Docker 
The docker hub image is coming soon...

## Usage

```shell
$ gocrawler -h
```

It displays the help menu.

```shell
Usage of gocrawler:
  -requests int
        Max requests per second (default -1)
  -scope string
        the file containing regexes for the scope
  -threads uint
        Number of max concurrent threads (default 10)
  -urls string
        The file containing all urls
```

> GoCrawler has to be used with a json scope file provided by the `-scope` file

```json
{
    "include": [
        "\\.hackerone\\.com$",
        "^api\\.hackerone\\.com$"
    ],

    "exclude": [
        "^docs\\.hackerone\\.com$",
        "^http://"
    ]
}
```

> This scope provides two arrays containing regexes.

- `include`: It contains all the regexes that has to be matched for the url
to be added.

- `exclude`: It contains a set of regexes to remove certain urls validated in
`include`.

```shell
$ echo "https://hackerone.com" | gocrawler -scope scope.json
```
or
```shell
$ gocrawler -scope scope.json -urls urls.txt
```

