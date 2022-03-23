package conf

var (
	// Application settings
	App struct {
		// ⚠️ WARNING: Should only be set by the main package (i.e. "gogs.go").
		Version string `ini:"-"`

		BrandName string
		RunUser   string
		RunMode   string
	}
)
