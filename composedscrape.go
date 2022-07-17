package composedscrape

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
		ctx: allocCtx,
	}
}

type Scraper struct {
	Cookies []*network.CookieParam // required: Name, Value, Domain: ".ope.ee"

	ctx context.Context
}

// sel: goquery selector
func (s *Scraper) Get(url, sel string) (_ *goquery.Selection, newURL string, _ error) { // by func(*chromedp.Selector)
	ctx, cancel := chromedp.NewContext(s.ctx)
	defer cancel()

	var gotHtml string
	actions := []chromedp.Action{
		chromedp.Navigate(url),
		chromedp.WaitReady(":root"),
		chromedp.OuterHTML(":root", &gotHtml),
		// chromedp.InnerHTML("document", &gotHtml, chromedp.ByJSPath),
		// chromedp.Nodes(sel, &nodes, by),
		chromedp.Location(&newURL),
	}

	if len(s.Cookies) > 0 {
		actions = append(
			[]chromedp.Action{network.SetCookies(s.Cookies)},
			actions...)
	}

	if err := chromedp.Run(ctx, actions...); err != nil {
		return nil, "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(gotHtml))
	if err != nil {
		return nil, "", err
	}

	// return nodes, nil
	return doc.Find(sel), newURL, nil
}
