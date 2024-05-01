package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllSitesAsStrings(t *testing.T) {
	sites := AllSitesAsStrings()
	assert.Len(t, sites, len(AllSites()))
}
