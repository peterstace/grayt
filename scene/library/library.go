package library

import (
	"sort"

	"github.com/peterstace/grayt/scene"
	"github.com/peterstace/grayt/scene/cornellbox"
)

var registry map[string]func() scene.Scene

func init() {
	registry = map[string]func() scene.Scene{
		"cornellbox_classic":    cornellbox.Classic,
		"cornellbox_splitbox":   cornellbox.Splitbox,
		"cornellbox_mirror":     cornellbox.Mirror,
		"cornellbox_spheretree": cornellbox.SphereTree,
	}
}

func Lookup(sceneName string) (func() scene.Scene, bool) {
	fn, ok := registry[sceneName]
	return fn, ok
}

func Listing() []string {
	var names []string
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
