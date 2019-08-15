/*
createreactappmanifest integrate Create React App with Go web application
*/
package manifest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/theplant/appkit/kerrs"
)

const manifestFileName = "asset-manifest.json"

const runtimeJSAssetKey = "runtime~main.js"
const mainJSAssetKey = "main.js"
const mainCSSAssetKey = "main.css"

var jsChunkRegexp = regexp.MustCompile(`^static/js/\d+\.[a-z0-9]+\.chunk\.js$`)
var cssChunkRegexp = regexp.MustCompile(`^static/css/\d+\.[a-z0-9]+\.chunk\.css$`)

type Config struct {
	// the PUBLIC_URL environment variable when do yarn build in react app
	PublicURL   string
	ManifestDir string

	// sometime you want to excludes the files in react app public folder.
	MountExcludeForPublic string
	IsDev                 bool

	/*
		Disable the code splitting in development mode, so you can get only one bundle at `http://localhost:3000/static/js/bundle.js`
		1. Install `@rescripts/cli` as a devDependency. `yarn add -D @rescripts/cli`
		2. Change the start script in package.json from "start": "react-scripts start" to "start": "rescripts start"
		3. Create a `.rescriptsrc.js` file at your project root with the following contents:
		```
		module.exports = config => {
		if (process.env.NODE_ENV === "development") {
			config.optimization.runtimeChunk = false;
			config.optimization.splitChunks = {
			cacheGroups: {
				default: false
			}
			};
		}

		return config;
		};
		```
	*/
	DevBundleURL string
}

type MData struct {
	JS  []string
	CSS []string
}

type Manifest struct {
	cfg        *Config
	mdata      MData
	mountFiles []os.FileInfo
	prefix     string
}

type Data struct {
	Files map[string]string
}

func New(cfg *Config) (m *Manifest, err error) {
	var data Data
	var mdata MData

	if cfg.IsDev {
		mdata = MData{
			JS:  []string{cfg.DevBundleURL},
			CSS: []string{},
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
		mdata = processAssetsURLs(data)
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
		mdata:      mdata,
		mountFiles: mfiles,
		prefix:     prefix,
	}
	return
}

func (mdata *MData) appendJS(url string) {
	mdata.JS = append(mdata.JS, url)
}

func (mdata *MData) appendCSS(url string) {
	mdata.CSS = append(mdata.CSS, url)
}

func processAssetsURLs(data Data) MData {
	files := data.Files
	mdata := MData{
		JS:  []string{},
		CSS: []string{},
	}

	// runtime js need to be the first one
	if runTimeJS, ok := files[runtimeJSAssetKey]; ok {
		mdata.appendJS(runTimeJS)
	}

	for key := range files {
		if jsChunkRegexp.MatchString(key) {
			mdata.appendJS(files[key])
			continue
		}
		if cssChunkRegexp.MatchString(key) {
			mdata.appendCSS(files[key])
		}
	}

	// main.js need to be the last one
	if mainJS, ok := files[mainJSAssetKey]; ok {
		mdata.appendJS(mainJS)
	}

	// main.css need to be the last one
	if mainCSS, ok := files[mainCSSAssetKey]; ok {
		mdata.appendCSS(mainCSS)
	}

	return mdata
}

func prefixURLs(prefix string, urls []string) []string {
	urlsWithPrefix := make([]string, len(urls))
	for i, url := range urls {
		urlsWithPrefix[i] = withPrefix(prefix, url)
	}
	return urlsWithPrefix
}

/*
GetJSURLs get all the js urls
*/
func (m *Manifest) GetJSURLs() (urls []string) {
	urls = m.mdata.JS
	if m.cfg.IsDev {
		return
	}

	return prefixURLs(m.prefix, urls)
}

/*
GetCSSURLs get all the css urls
*/
func (m *Manifest) GetCSSURLs() (urls []string) {
	urls = m.mdata.CSS
	if m.cfg.IsDev {
		return
	}

	return prefixURLs(m.prefix, urls)
}

func withPrefix(prefix, p string) string {
	// in case p included the prefix, fix /cms/cms/static/css/main.c17080f1.css error
	if strings.Index(filepath.Join("/", p), prefix) == 0 {
		prefix = "/"
	}
	return filepath.Join(prefix, p)
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
