package composedscrape

import (
	"context"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type ScraperOpts struct {
	Executable string
}

func NewScraper() *Scraper {
	ctx := context.Background() // to be implemented
	allocCtx, _ := chromedp.NewExecAllocator(ctx, chromedp.DefaultExecAllocatorOptions[:]...)
	return &Scraper{
		ctx: allocCtx,
	}
}

type Scraper struct {
	Cookies []*network.CookieParam

	ctx context.Context
}

// startFunc should return one (non-slice) object
func (s *Scraper) Get(url, sel string, by func(*chromedp.Selector)) (nodes []*cdp.Node, _ error) {
	ctx, cancel := chromedp.NewContext(s.ctx)
	defer cancel()

	actions := []chromedp.Action{
		chromedp.Navigate(url),
		chromedp.WaitReady(":root"),
		chromedp.Nodes(sel, &nodes, by),
	}

	if len(s.Cookies) > 0 {
		actions = append([]chromedp.Action{network.SetCookies(s.Cookies)}, actions...)
	}

	if err := chromedp.Run(ctx, actions...); err != nil {
		return nil, err
	}
	return nodes, nil
}
