package content_parser

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

func Parse(r io.Reader, results chan int) error {
	fmt.Println("STARTING TO PARSE THE CONTENT")
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return z.Err()
		}
	}
}
