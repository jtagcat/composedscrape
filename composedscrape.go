package composedscrape

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type ScraperOpts struct {
	Executable string
}

// populates internal values on Scraper
func NewScraper(raw *Scraper) *Scraper {
	ctx := context.Background() // to be implemented
	opts := chromedp.DefaultExecAllocatorOptions[:]
	opts = append(opts, raw.InitExtraAllocatorOpts...)

	allocCtx, _ := chromedp.NewExecAllocator(ctx, opts...)
	raw.ctx = allocCtx

	raw.downloadsLimit = make(chan bool, downloadsMaxActive)

	if raw.InitGlobalConcurrentLimit <= 0 {
		raw.InitGlobalConcurrentLimit = 32
	}
	raw.globalLimit = make(chan bool, raw.InitGlobalConcurrentLimit)
	return raw
}

const downloadsMaxActive = 10

// use NewScraper to initialize internal values
type Scraper struct {
	Cookies []*network.CookieParam // required: Name, Value, Domain: ".ope.ee"
	Timeout time.Duration          // 0: disabled

	// Only works after NewScraper()
	InitExtraAllocatorOpts    []chromedp.ExecAllocatorOption
	InitGlobalConcurrentLimit int

	ctx            context.Context
	downloadsLimit chan bool
	globalLimit    chan bool
}

func (s *Scraper) limitLock() error {
	select {
	case <-s.ctx.Done():
		return context.Canceled
	case s.globalLimit <- false: // false is placeholder
		return nil
	}
}

func (s *Scraper) limitUnlock() {
	<-s.globalLimit
}
