package composedscrape

import (
	"context"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

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

// based on https://github.com/chromedp/examples/blob/3384adb2158f6df7e6a48458875a3a5f24aea0c3/download_file/main.go
func (s *Scraper) DownloadFile(url, outdir string, timeout time.Duration) (suggested, filename, newURL string, _ error) {
	ctx, cancel := chromedp.NewContext(s.ctx)
	defer cancel()

	// create a timeout as a safety net to prevent any infinite wait loops
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	// handle download event
	done := make(chan string, 1)
	chromedp.ListenTarget(ctx, func(v interface{}) {
		if ev, ok := v.(*browser.EventDownloadWillBegin); ok {
			suggested = ev.SuggestedFilename
		}

		if ev, ok := v.(*browser.EventDownloadProgress); ok {
			if ev.State == browser.DownloadProgressStateCompleted {
				done <- ev.GUID
				close(done)
			}
		}
	})

	actions := []chromedp.Action{
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(outdir).
			WithEventsEnabled(true),
		chromedp.Navigate(url),
		chromedp.Location(&newURL),
	}
	if len(s.Cookies) > 0 {
		actions = append(
			[]chromedp.Action{network.SetCookies(s.Cookies)},
			actions...)
	}

	if err := chromedp.Run(ctx, actions...); err != nil && !strings.Contains(err.Error(), "net::ERR_ABORTED") {
		// Upstream note: Ignoring the net::ERR_ABORTED page error is essential here
		// since downloads will cause this error to be emitted, although the
		// download will still succeed.
		return "", "", "", err
	}

	guid := <-done // blocks

	return suggested, guid, newURL, nil
}
