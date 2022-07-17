package composedscrape

import (
	"context"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type ScraperOpts struct {
	Executable string
}

func NewScraper(extraAllocatorOpts ...chromedp.ExecAllocatorOption) *Scraper {
	ctx := context.Background() // to be implemented
	opts := chromedp.DefaultExecAllocatorOptions[:]
	opts = append(opts, extraAllocatorOpts...)

	allocCtx, _ := chromedp.NewExecAllocator(ctx, opts...)
	return &Scraper{
		ctx:            allocCtx,
		downloadsQueue: make(chan bool, downloadsMaxActive),
	}
}

type Scraper struct {
	Cookies []*network.CookieParam // required: Name, Value, Domain: ".ope.ee"

	ctx            context.Context
	downloadsQueue chan bool
}

const downloadsMaxActive = 10
