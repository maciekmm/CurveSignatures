package managers

import (
	"errors"
	"github.com/maciekmm/curvesignatures/models"
)

var layouts map[string]models.Layout

func init() {
	layouts = make(map[string]models.Layout)
}

// Not thread safe, but it will be called only on startup
func RegisterLayout(id string, layout models.Layout) error {
	_, present := layouts[id]
	if present {
		return errors.New("Could not register layout - layout with that name already exists")
	}
	layouts[id] = layout
	return nil
}

// Gets layout based on id
func GetLayoutById(id string) models.Layout {
	layout, present := layouts[id]
	if !present {
		return nil
	}
	return layout
}

// Get all registered layouts
func GetRegisteredLayoutNames() map[string]models.Layout {
	/*lays := make([]string, len(layouts))
	i := 0
	for k, _ := range layouts {
		lays[i] = k
		i++
	}*/
	return layouts
}
