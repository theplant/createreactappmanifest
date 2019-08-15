package manifest_test

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/theplant/createreactappmanifest"
)

/*
Mount expose assets with ServeMux, and GetJSURLs, GetCSSURLs get correct assets path for you.
*/
func ExampleNew() {
	mux := http.DefaultServeMux

	m, _ := manifest.New(&manifest.Config{
		ManifestDir: "./example/build",
		PublicURL:   "/cms",
	})

	renderTemplate := func() string {
		buf := bytes.NewBuffer(nil)

		jsURLs := m.GetJSURLs()
		for _, url := range jsURLs {
			fmt.Fprintf(buf, `<script type="text/javascript" src="%s"></script>`, url)
		}

		cssURLs := m.GetCSSURLs()
		for _, url := range cssURLs {
			fmt.Fprintf(buf, `<link href="%s" rel="stylesheet">`, url)
		}

		return buf.String()
	}

	m.Mount(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, renderTemplate())
	})

}
