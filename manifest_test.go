package manifest_test

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/theplant/testingutils"

	manifest "github.com/theplant/createreactappmanifest"
)

var exposedURLs = []string{
	"/asset-manifest.json",
	"/favicon.ico",
	"/logo192.png",
	"/logo512.png",
	"/manifest.json",
	"/robots.txt",
	"/service-worker.js",
	"/static/css/main.b100e6da.chunk.css",
	"/static/js/main.f2081a32.chunk.js",
	"/static/js/main.f2081a32.chunk.js.map",
	"/static/js/runtime-main.62087466.js",
	"/static/js/runtime-main.62087466.js.map",
	"/static/css/2.764ccc25.chunk.css",
	"/static/js/2.76f9f456.chunk.js",
	"/static/js/2.76f9f456.chunk.js.map",
	"/precache-manifest.3a2c3995787ad241ddb8cff5b35d6418.js",
	"/static/css/2.764ccc25.chunk.css.map",
	"/static/css/main.b100e6da.chunk.css.map",
	"/static/media/logo.25bf045c.svg",
}

var exposedURLsWithPrefix = []string{
	"/cms/asset-manifest.json",
	"/cms/favicon.ico",
	"/cms/logo192.png",
	"/cms/logo512.png",
	"/cms/manifest.json",
	"/cms/robots.txt",
	"/cms/service-worker.js",
	"/cms/static/css/main.b100e6da.chunk.css",
	"/cms/static/js/main.f2081a32.chunk.js",
	"/cms/static/js/main.f2081a32.chunk.js.map",
	"/cms/static/js/runtime-main.62087466.js",
	"/cms/static/js/runtime-main.62087466.js.map",
	"/cms/static/css/2.764ccc25.chunk.css",
	"/cms/static/js/2.76f9f456.chunk.js",
	"/cms/static/js/2.76f9f456.chunk.js.map",
	"/cms/precache-manifest.3a2c3995787ad241ddb8cff5b35d6418.js",
	"/cms/static/css/2.764ccc25.chunk.css.map",
	"/cms/static/css/main.b100e6da.chunk.css.map",
	"/cms/static/media/logo.25bf045c.svg",
}

var jsURLs = []string{"/static/js/runtime-main.62087466.js", "/static/js/2.76f9f456.chunk.js", "/static/js/main.f2081a32.chunk.js"}
var jsURLsWithPrefix = []string{"/cms/static/js/runtime-main.62087466.js", "/cms/static/js/2.76f9f456.chunk.js", "/cms/static/js/main.f2081a32.chunk.js"}

var cssURLs = []string{"/static/css/2.764ccc25.chunk.css", "/static/css/main.b100e6da.chunk.css"}
var cssURLsWithPrefix = []string{"/cms/static/css/2.764ccc25.chunk.css", "/cms/static/css/main.b100e6da.chunk.css"}

var cases = []struct {
	name           string
	cfg            *manifest.Config
	jsURLs         []string
	cssURLs        []string
	exposedURLs    []string
	notExposedURLs []string
}{
	{
		name: "expose all",
		cfg: &manifest.Config{
			PublicURL:   "/",
			ManifestDir: "./example/build",
			IsDev:       false,
		},
		jsURLs:      jsURLs,
		cssURLs:     cssURLs,
		exposedURLs: exposedURLs,
	},
	{
		name: "expose all without PublicURL",
		cfg: &manifest.Config{
			ManifestDir: "./example/build",
			IsDev:       false,
		},
		jsURLs:      jsURLs,
		cssURLs:     cssURLs,
		exposedURLs: exposedURLs,
	},
	{
		name: "with public url prefix /cms",
		cfg: &manifest.Config{
			PublicURL:   "/cms",
			ManifestDir: "./example/build",
			IsDev:       false,
		},
		jsURLs:      jsURLsWithPrefix,
		cssURLs:     cssURLsWithPrefix,
		exposedURLs: exposedURLsWithPrefix,
	},
	{
		name: "exclude top level *.html",
		cfg: &manifest.Config{
			PublicURL:             "cms",
			ManifestDir:           "./example/build",
			MountExcludeForPublic: "*.html",
			IsDev:                 false,
		},
		jsURLs:      jsURLsWithPrefix,
		cssURLs:     cssURLsWithPrefix,
		exposedURLs: exposedURLsWithPrefix,
		notExposedURLs: []string{
			"/cms/index.html",
		},
	},
	{
		name: "dev with PublicURL",
		cfg: &manifest.Config{
			PublicURL:    "/cms",
			ManifestDir:  "./example/build",
			IsDev:        true,
			DevBundleURL: "http://localhost:3000/static/js/bundle.js",
		},
		jsURLs:  []string{"http://localhost:3000/static/js/bundle.js"},
		cssURLs: []string{},
	},
	{
		name: "dev without PublicURL",
		cfg: &manifest.Config{
			PublicURL:    "",
			ManifestDir:  "./example/build",
			IsDev:        true,
			DevBundleURL: "http://localhost:3000/static/js/bundle.js",
		},
		jsURLs:  []string{"http://localhost:3000/static/js/bundle.js"},
		cssURLs: []string{},
	},
	{
		name: "asset-manifest.json already have prefix /cms",
		cfg: &manifest.Config{
			PublicURL:   "cms",
			ManifestDir: "./example2/build",
			IsDev:       false,
		},
		jsURLs:  []string{"/cms/static/js/runtime-main.124c78f9.js", "/cms/static/js/2.cc7f12e5.chunk.js", "/cms/static/js/main.ca0dff28.chunk.js"},
		cssURLs: []string{"/cms/static/css/main.b100e6da.chunk.css"},
		exposedURLs: []string{
			"/cms/static/js/runtime-main.124c78f9.js",
			"/cms/static/js/2.cc7f12e5.chunk.js",
			"/cms/static/js/main.ca0dff28.chunk.js",
			"/cms/static/css/main.b100e6da.chunk.css",
		},
	},
}

func TestGetURL_Mount(t *testing.T) {
	for _, c := range cases {
		mni, err := manifest.New(c.cfg)
		if err != nil {
			t.Fatal(err)
		}

		jsURLs := mni.GetJSURLs()
		jsURLsDiff := testingutils.PrettyJsonDiff(jsURLs, c.jsURLs)
		if len(jsURLsDiff) > 0 {
			t.Error(jsURLsDiff)
		}

		cssURLs := mni.GetCSSURLs()
		cssURLsDiff := testingutils.PrettyJsonDiff(cssURLs, c.cssURLs)
		if len(cssURLsDiff) > 0 {
			t.Error(cssURLsDiff)
		}

		mux := http.NewServeMux()
		mni.Mount(mux)

		for _, u := range c.exposedURLs {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", u, nil)
			mux.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				res, _ := httputil.DumpResponse(rr.Result(), false)
				t.Errorf("should get %s OK, but was %d, response is: \n%s", u, rr.Code, string(res))
			}
		}

		for _, u := range c.notExposedURLs {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", u, nil)
			mux.ServeHTTP(rr, req)
			if rr.Code != http.StatusNotFound {
				res, _ := httputil.DumpResponse(rr.Result(), false)
				t.Errorf("should get %s StatusNotFound, but was %d, response is: \n%s", u, rr.Code, string(res))
			}
		}
	}
}
