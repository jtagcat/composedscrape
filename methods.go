package composedscrape

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
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
func (s *Scraper) DownloadFile(urlstr, outdir string, timeout time.Duration) (suggested, filename, newURL string, _ error) {
	//## block until we take a spot in the queue or parent ctx cancelled
	select {
	case <-s.ctx.Done():
		return "", "", "", context.Canceled
	case s.downloadsQueue <- false: // false is placeholder

	}
	defer func() {
		<-s.downloadsQueue // leave queue
	}()

	ctx, cancel := chromedp.NewContext(s.ctx)
	defer cancel()

	// create a timeout as a safety net to prevent any infinite wait loops
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan string, 1)

	//## handle download event
	var requestID network.RequestID
	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		// opt a: browser renders (png)
		case *network.EventRequestWillBeSent:
			if ev.Request.URL == urlstr {
				requestID = ev.RequestID
			}
		case *network.EventLoadingFinished:
			if ev.RequestID == requestID {
				close(done)
			}

		// opt b: direct download
		case *browser.EventDownloadWillBegin:
			suggested = ev.SuggestedFilename

		case *browser.EventDownloadProgress:
			if ev.State == browser.DownloadProgressStateCompleted {
				done <- ev.GUID
				close(done)
			}

		}
	})

	//## direct chrome interaction
	actions := []chromedp.Action{
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(outdir).
			WithEventsEnabled(true),
		chromedp.Navigate(urlstr),
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

	//## opt a: if browser rendered and is not direct donwload (eg. png)
	if guid == "" {
		// emulate a guid manually
		var guidPath string
		for {
			guid = uuid.New().String()
			guidPath = path.Join(outdir, guid)
			if _, err := os.Stat(guidPath); err != nil && !errors.Is(err, os.ErrNotExist) {
				return "", "", "", err
			} else if errors.Is(err, os.ErrNotExist) {
				break
			}
		}

		// get the downloaded bytes by request id
		var buf []byte
		if err := chromedp.Run(ctx,
			chromedp.ActionFunc(func(ctx context.Context) (err error) {
				buf, err = network.GetResponseBody(requestID).Do(ctx)
				return err
			})); err != nil {
			return "", "", "", err
		}

		if err := ioutil.WriteFile(guidPath, buf, os.ModePerm); err != nil {
			return "", "", "", err
		}
	}

	return suggested, guid, newURL, nil
}
