package composedscrape_test

import (
	"testing"

	"github.com/chromedp/chromedp"
	cs "github.com/jtagcat/composedscrape"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	s := cs.NewScraper()
	// s.Executable = "chromium"
	nodes, err := s.Get("https://www.c7.ee/", "document", chromedp.ByJSPath)
	assert.Nil(t, err)
	print(nodes)
	panic("")
}
