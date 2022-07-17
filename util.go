package composedscrape

import (
	"encoding/json"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// writes object to file as json
func JsonToFile(filename, indent string, object interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	je := json.NewEncoder(f)
	je.SetIndent("", indent)

	return je.Encode(object)
}

// k8s.io/helm/pkg/urlutil
// mod: supports // and /
// in case of multiple absolute paths, last is used
func URLJoin(baseURL string, paths ...string) (string, error) {
	// mod:
	// base is replaced by first with //
	var newBase int
	for i, p := range paths {
		if strings.HasPrefix(p, "//") {
			newBase = i
		}
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	if newBase != 0 {
		paths = paths[newBase+1:]

		old := u
		u, err = url.Parse(paths[newBase])
		if err != nil {
			return "", err
		}

		u.Scheme = old.Scheme
		if u.User == nil {
			u.User = old.User
		}
	}

	// mod:
	// allow rooting to domain with /
	var absPath int
	for i, p := range paths {
		if strings.HasPrefix(p, "/") {
			absPath = i
		}
	}

	// We want path instead of filepath because path always uses /.
	if absPath != 0 {
		u.Path = path.Join(paths[absPath:]...)
	} else {
		all := []string{u.Path}
		all = append(all, paths...)
		u.Path = path.Join(all...)
	}

	return u.String(), nil
}

// if only we could have `type Node cdp.Node` to use `func (n *Node)`
//
// if https://github.com/chromedp/cdproto/issues/20 is implemented, this func is deprecated

// Converts cdp.Node to goquery and filters children.
//
// panics: valid HTML almost always remains valid HTML
// func CdpFilterChildren(node *cdp.Node, sel string) *goquery.Selection {
// 	doc, err := goquery.NewDocumentFromReader(strings.NewReader(
// 		node.Dump("", "", false)))
// 	if err != nil {
// 		panic(fmt.Errorf("converting cdp.Node (%v) to goquery Doc: %e", node, err))
// 	}

// 	return doc.Selection.ChildrenFiltered(sel)
// }

// let's play dumb (unavail/-exported newSingleSelection), redoing things for nth time already
func RawEach(s *goquery.Selection) (a []*goquery.Selection) {
	s.Each(func(_ int, s *goquery.Selection) {
		a = append(a, s)
	})
	return a
}
