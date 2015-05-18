package models

import (
	"sort"
	"strings"
)

// Configuration represents configuration for layout generation telling it what ranks it should use
type Configuration struct {
	Ranks []string //Ranks to display on signature
}

type RequestParameters struct {
	LayoutName string
	PlayerID   int
	Layout     Layout
	Ranks      *Configuration
}

// CombineRanks combine ranks into file name format
func (conf Configuration) CombineRanks() string {
	sort.Strings(conf.Ranks)
	return strings.Join(conf.Ranks, "")
}
