# gostruct

**gostruct** populates Go `struct`s from [goquery][]’s `Document`s.

[goquery]: https://github.com/PuerkitoBio/goquery

## Install

    go get github.com/bfontaine/gostruct

## Usage

```go
type MyStruct struct {
    Title string `goquery:"#a-selector"`
    Desc string `goquery:"h1 .another .one"`
}

var ms MyStruct

gostruct.Fetch(&ms, "http://www.example.com")
```

## Example

The example program below searchs for "golang" on Google and prints the top
results.

```go
package main

import (
    "fmt"
    "os"

    "github.com/bfontaine/gostruct"
)

type Result struct {
    Title       string `gostruct:"h3.r"`
    Website     string `gostruct:".kv cite"`
    Description string `gostruct:".st"`
}

type Search struct {
    Results []Result `gostruct:"#search li.g"`
}

func main() {
    var s Search

    err := gostruct.Fetch(&s, "https://www.google.com/search?q=golang&hl=en")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    for _, r := range s.Results {
        fmt.Printf("%s (%s) - %s\n", r.Title, r.Website, r.Description)
    }
}
```
