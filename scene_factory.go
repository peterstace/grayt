package grayt

// SceneFactory is implemented by scene writers.
type SceneFactory interface {
	NewScene(t float64) // 0 <= t < 1
}
