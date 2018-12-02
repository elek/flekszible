package helm

type Values struct {
	Image     Image
	Component map[string]Component
}

type Image struct {
	Repository string
	Tag        string
	PullPolicy string
}

type Component struct {
	Enabled bool
	Replica int
}
