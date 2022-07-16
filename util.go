package composedscrape

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path"
	"strings"
)

// writes object to file as json
func JsonToFile(filename string, object interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	je := json.NewEncoder(f)
	return je.Encode(object)
}

var errorTooManyAbsolute = errors.New("cannot join more than 2 absolute paths")

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
