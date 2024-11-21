package r2w_test

import (
	"fmt"
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestRssLib(t *testing.T) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://www.reddit.com/r/golang.rss")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(feed.Title)
}
