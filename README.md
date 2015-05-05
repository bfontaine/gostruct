# gostruct

**gostruct** populates Go `struct`s from [goquery][]’s `Document`s.

[goquery]: https://github.com/PuerkitoBio/goquery

## Example

```go
import (
    "fmt"
    "log"

    "github.com/bfontaine/gostruct"
)


type Project struct {
    Title  string `gostruct:".js-current-repository"`
    Author string `gostruct:".entry-title [rel=author]"`
    Desc   string `gostruct:".repository-description"`
}


func main() {
    var p Project

    if err := gostruct.Fetch(&p, "https://github.com/bfontaine/gostruct"); err != nil {
        log.Fatalln(err)
    }

    fmt.Printf("%s/%s: %s\n", p.Author, p.Title, p.Desc)
}
```
