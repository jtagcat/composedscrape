package composedscrape

import (
	"encoding/json"
	"errors"
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

var errorTooManyAbsolute = errors.New("cannot join more than 2 absolute paths")

// TODO: BUG:? is mode even needed? docs?
// k8s.io/helm/pkg/urlutil
func UrlJoin(baseURL string, paths ...string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// mod:
	if strings.HasPrefix(paths[0], "/") {
		for _, u := range paths[1:] {
			if strings.HasPrefix(u, "/") {
				return "", errorTooManyAbsolute
			}
		}

		u.Path = paths[0]
		return u.String(), nil
	}

	// We want path instead of filepath because path always uses /.
	all := []string{u.Path}
	all = append(all, paths...)
	u.Path = path.Join(all...)
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
