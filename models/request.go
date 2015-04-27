package models

import (
	"sort"
	"strings"
)

type Configuration struct {
	Ranks []string //Ranks to display on signature
}

type RequestParameters struct {
	LayoutName string
	PlayerID   int
	Layout     Layout
	Ranks      *Configuration
}

// Combine ranks to form ID
func (conf Configuration) CombineRanks() string {
	sort.Strings(conf.Ranks)
	return strings.Join(conf.Ranks, "")
}
