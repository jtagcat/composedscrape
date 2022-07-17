package composedscrape

import (
	"context"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type ScraperOpts struct {
	Executable string
}

// populates internal values on Scraper
func NewScraper(raw *Scraper, extraAllocatorOpts ...chromedp.ExecAllocatorOption) *Scraper {
	ctx := context.Background() // to be implemented
	opts := chromedp.DefaultExecAllocatorOptions[:]
	opts = append(opts, extraAllocatorOpts...)

	allocCtx, _ := chromedp.NewExecAllocator(ctx, opts...)
	raw.ctx = allocCtx

	raw.downloadsQueue = make(chan bool, downloadsMaxActive)
	return raw
}

// use NewScraper to initialize internal values
type Scraper struct {
	Cookies []*network.CookieParam // required: Name, Value, Domain: ".ope.ee"

	ctx            context.Context
	downloadsQueue chan bool
}

const downloadsMaxActive = 10
