package composedscrape_test

import (
	"testing"

	cs "github.com/jtagcat/composedscrape"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	s := cs.NewScraper()
	// s.Executable = "chromium"
	nodes, _, err := s.Get("https://www.c7.ee/", "document")
	assert.Nil(t, err)
	print(nodes)
	panic("test not implemented")
}
