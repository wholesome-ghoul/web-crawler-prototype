package content_parser

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Content struct {
	Url  string
	Path string
}

func Parse(crawled Content) error {
	counter := 0
	if crawled.Path == "" {
		return nil
	}

	fileReader, _ := os.Open(crawled.Path)
	defer fileReader.Close()

	z := html.NewTokenizer(fileReader)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			fmt.Println("A", crawled.Url, counter)
			return z.Err()
		}

		if tt == html.StartTagToken {
			t := z.Token()
			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" {
						if a.Val == "" || a.Val == "#" || a.Val == "/" {
							continue
						}

						if strings.HasPrefix(a.Val, "/") {
							a.Val = crawled.Url + a.Val
						}

						// fmt.Println(a.Val)
						counter++
					}
				}
			}
		}
	}
}
