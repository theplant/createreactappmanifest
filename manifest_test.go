package manifest_test

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/theplant/testingutils"

	manifest "github.com/theplant/createreactappmanifest"
)

var cases = []struct {
	name           string
	cfg            *manifest.Config
	getURLNames    [][]string
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
		getURLNames: [][]string{
			[]string{"main.css", "/static/css/main.c17080f1.css"},
			[]string{"main.css.map", "/static/css/main.c17080f1.css.map"},
			[]string{"main.js", "/static/js/main.33fb4ad2.js"},
			[]string{"main.js.map", "/static/js/main.33fb4ad2.js.map"},
			[]string{"static/media/logo.svg", "/static/media/logo.5d5d9eef.svg"},
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
			"/img/logo.jpg",
			"/javascripts/pace.js",
		},
	},
	{
		name: "with public url prefix /cms",
		cfg: &manifest.Config{
			PublicURL:   "/cms",
			ManifestDir: "./example/build",
			IsDev:       false,
		},
		getURLNames: [][]string{
			[]string{"main.css", "/cms/static/css/main.c17080f1.css"},
			[]string{"main.css.map", "/cms/static/css/main.c17080f1.css.map"},
			[]string{"main.js", "/cms/static/js/main.33fb4ad2.js"},
			[]string{"main.js.map", "/cms/static/js/main.33fb4ad2.js.map"},
			[]string{"static/media/logo.svg", "/cms/static/media/logo.5d5d9eef.svg"},
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
			"/cms/img/logo.jpg",
			"/cms/javascripts/pace.js",
		},
	},
	{
		name: "exclude top level *.html",
		cfg: &manifest.Config{
			PublicURL:             "/cms",
			ManifestDir:           "./example/build",
			MountExcludeForPublic: "*.html",
			IsDev: false,
		},
		getURLNames: [][]string{
			[]string{"demo.html", "/cms/demo.html"},
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
			url := mni.GetURL(name[0])
			diff := testingutils.PrettyJsonDiff(name[1], url)
			if len(diff) > 0 {
				t.Error(diff)
			}
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
				t.Errorf("should get %s OK, but was %d, response is: \n%s", u, rr.Code, string(res))
			}
		}
	}
}
