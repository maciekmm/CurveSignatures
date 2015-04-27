package managers

import (
	"errors"
	"github.com/maciekmm/curvesignatures/models"
	"log"
	"os"
	"regexp"
	"strings"
)

var regions *[3]models.Region
var regionRegex *regexp.Regexp = regexp.MustCompile("([a-z0-9]+)_([a-z]+)")

func init() {
	checkErrors := func(ex error) {
		if ex != nil {
			log.Fatalln("Could not load file.", ex)
			os.Exit(1)
		}
	}
	region_asia, err := LoadImageFromFile("./assets/asia.png")
	checkErrors(err)
	region_europe, err := LoadImageFromFile("./assets/eu.png")
	checkErrors(err)
	region_murica, err := LoadImageFromFile("./assets/usa.png")
	checkErrors(err)

	regions = &[...]models.Region{
		models.Region{"eu", region_europe, "Europe"},
		models.Region{"asia", region_asia, "Asia"},
		models.Region{"usa", region_murica, "America"},
	}
}

// Gets region based on id
func GetRegion(region string) (models.Region, error) {
	for _, val := range regions {
		if val.Suffix == region {
			return val, nil
		}
	}
	return models.Region{}, errors.New("Could not find region " + region)
}

// Extracts region from the rank id
func GetRegionFromRank(rank string) (models.Region, error) {
	if !regionRegex.MatchString(rank) {
		return models.Region{}, errors.New("Invalid syntax " + rank)
	}
	reg := regionRegex.FindStringSubmatch(rank)
	region, err := GetRegion(reg[2])
	if err != nil {
		return models.Region{}, err
	}
	return region, nil
}

// Converts region to readable format
// TODO Move this to relevant place
func ConvertToReadableName(rank string) (string, error) {
	if !regionRegex.MatchString(rank) {
		return "", errors.New("Invalid syntax " + rank)
	}
	reg := regionRegex.FindStringSubmatch(rank)
	return strings.ToUpper(reg[1]), nil
}
