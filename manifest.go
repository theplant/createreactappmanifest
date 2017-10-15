/*
createreactappmanifest integrate [Create React App](https://github.com/facebookincubator/create-react-app) with Go web application
*/
package manifest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/theplant/appkit/kerrs"
)

const manifestFileName = "asset-manifest.json"

type Config struct {
	// the PUBLIC_URL environment variable when do yarn build in react app
	PublicURL   string
	ManifestDir string

	// sometime you want to excludes the files in react app public folder.
	MountExcludeForPublic string
	IsDev                 bool
	DevBundleURL          string
}

type Manifest struct {
	cfg        *Config
	mdata      map[string]string
	mountFiles []os.FileInfo
	prefix     string
}

func New(cfg *Config) (m *Manifest, err error) {
	var data map[string]string

	if cfg.IsDev {
		data = map[string]string{
			"main.css":     "",
			"main.css.map": "",
			"main.js":      cfg.DevBundleURL,
			"main.js.map":  cfg.DevBundleURL + ".map",
		}
	} else {
		var f *os.File
		f, err = os.Open(filepath.Join(cfg.ManifestDir, manifestFileName))
		if err != nil {
			err = kerrs.Wrapv(err, "failed to open asset-manifest.json file")
			return
		}
		defer f.Close()
		err = json.NewDecoder(f).Decode(&data)
		if err != nil {
			err = kerrs.Wrapv(err, "failed to decode asset-manifest.json file")
			return
		}
	}

	var publicfiles, mfiles []os.FileInfo
	publicfiles, err = ioutil.ReadDir(cfg.ManifestDir)
	if err != nil {
		err = kerrs.Wrapv(err, "failed to glob files", "manifest_dir", cfg.ManifestDir)
		return
	}

	for _, f := range publicfiles {
		if matched, _ := filepath.Match(cfg.MountExcludeForPublic, f.Name()); !matched {
			mfiles = append(mfiles, f)
		}
	}
	prefix := filepath.Join("/", cfg.PublicURL)
	m = &Manifest{
		cfg:        cfg,
		mdata:      data,
		mountFiles: mfiles,
		prefix:     prefix,
	}
	return
}

/*
GetURL get dynamic compiled url like `/static/css/main.c17080f1.css` with simple name like `main.css` to be used in views
*/
func (m *Manifest) GetURL(name string) (url string) {
	p := name
	if urlpart, ok := m.mdata[name]; ok {
		p = urlpart
	}
	if m.cfg.IsDev {
		return p
	}
	url = filepath.Join(m.prefix, p)
	return
}

/*
Mount automatically mounts Create React App build directory into Go ServeMux
*/
func (m *Manifest) Mount(mux *http.ServeMux) {
	h := http.FileServer(http.Dir(m.cfg.ManifestDir))
	if m.prefix != "/" {
		h = http.StripPrefix(m.prefix, h)
	}

	for _, mf := range m.mountFiles {
		pt := filepath.Join(m.prefix, mf.Name())
		if mf.IsDir() {
			mux.Handle(pt+"/", h)
		} else {
			mux.Handle(pt, h)
		}
	}
	return
}
