package manifest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	manifest "github.com/theplant/createreactappmanifest"
)

var cases = []struct {
	name           string
	cfg            *manifest.Config
	getURLNames    []string
	exposedURLs    []string
	notExposedURLs []string
}{
	{
		name: "expose all",
		cfg: &manifest.Config{
			PublicURL:   "/",
			ManifestDir: "./example",
			IsDev:       false,
		},
		getURLNames: []string{
			"main.css",
			"main.css.map",
			"main.js",
			"main.js.map",
			"static/media/logo.svg",
		},
		exposedURLs: []string{
			"/favicon.ico",
			"/service-worker.js",
			"/asset-manifest.json",
			"/static/css/main.c17080f1.css",
			"/static/css/main.c17080f1.css.map",
			"/static/js/main.33fb4ad2.js",
			"/static/js/main.33fb4ad2.js.map",
			"/static/media/logo.5d5d9eef.svg",
			"/demo.html",
			"/index.html",
			"/img/logo.jpg",
			"/javascripts/pace.js",
		},
	},
	{
		name: "with public url prefix /cms",
		cfg: &manifest.Config{
			PublicURL:   "/cms",
			ManifestDir: "./example",
			IsDev:       false,
		},
		getURLNames: []string{
			"main.css",
			"main.css.map",
			"main.js",
			"main.js.map",
			"static/media/logo.svg",
		},
		exposedURLs: []string{
			"/cms/favicon.ico",
			"/cms/service-worker.js",
			"/cms/asset-manifest.json",
			"/cms/static/css/main.c17080f1.css",
			"/cms/static/css/main.c17080f1.css.map",
			"/cms/static/js/main.33fb4ad2.js",
			"/cms/static/js/main.33fb4ad2.js.map",
			"/cms/static/media/logo.5d5d9eef.svg",
			"/cms/demo.html",
			"/cms/index.html",
			"/cms/img/logo.jpg",
			"/cms/javascripts/pace.js",
		},
	},
	{
		name: "exclude top level *.html",
		cfg: &manifest.Config{
			PublicURL:             "/cms",
			ManifestDir:           "./example",
			MountExcludeForPublic: "*.html",
			IsDev: false,
		},
		getURLNames: []string{
			"main.css",
			"main.css.map",
			"main.js",
			"main.js.map",
			"static/media/logo.svg",
		},
		exposedURLs: []string{
			"/cms/favicon.ico",
			"/cms/service-worker.js",
			"/cms/asset-manifest.json",
			"/cms/static/css/main.c17080f1.css",
			"/cms/static/css/main.c17080f1.css.map",
			"/cms/static/js/main.33fb4ad2.js",
			"/cms/static/js/main.33fb4ad2.js.map",
			"/cms/static/media/logo.5d5d9eef.svg",
			"/cms/img/logo.jpg",
			"/cms/javascripts/pace.js",
		},
		notExposedURLs: []string{
			"/cms/demo.html",
			"/cms/index.html",
		},
	},
}

func TestMount(t *testing.T) {
	for _, c := range cases {
		mni, err := manifest.New(c.cfg)
		if err != nil {
			t.Fatal(err)
		}

		for _, name := range c.getURLNames {
			url := mni.GetURL(name)
			if len(url) == 0 {
				t.Errorf("should get url for name %s, but didn't get it", name)
			}
		}

		mux := http.NewServeMux()
		mni.Mount(mux)

		for _, u := range c.exposedURLs {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", u, nil)
			mux.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Errorf("should get %s OK, but was %d", u, rr.Code)
			}
		}

		for _, u := range c.notExposedURLs {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", u, nil)
			mux.ServeHTTP(rr, req)
			if rr.Code != http.StatusNotFound {
				t.Errorf("should get %s NotFound, but was %d", u, rr.Code)
			}
		}
	}
}
