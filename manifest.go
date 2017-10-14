package manifest

import (
	"net/http"
)

const manifestFileName = "asset-manifest.json"

type Config struct {
	PublicURL   string
	ManifestDir string

	// sometime you want to excludes the files in react app public folder.
	MountExcludeForPublic string
	IsDev                 bool
	DevBundleURL          string
}

type Manifest struct {
	cfg *Config
}

func New(cfg *Config) (m *Manifest, err error) {
	err = validateConfig(cfg)
	if err != nil {
		return
	}

	m = &Manifest{
		cfg: cfg,
	}
	return
}

func validateConfig(cfg *Config) (err error) {
	return
}

func (m *Manifest) GetURL(name string) (url string) {
	return
}

func (m *Manifest) Mount(mux *http.ServeMux) {
	return
}
