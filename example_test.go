package manifest_test

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/theplant/createreactappmanifest"
)

/*
Mount expose assets with ServeMux, and GetURL get correct assets path for you.
*/
func ExampleNew() {
	mux := http.DefaultServeMux

	m, _ := manifest.New(&manifest.Config{
		ManifestDir: "./example/build",
		PublicURL:   "/cms",
	})

	renderTemplate := func() string {
		buf := bytes.NewBuffer(nil)
		fmt.Fprintf(buf, `<link href="%s" rel="stylesheet">`, m.GetURL("main.css"))
		fmt.Fprintf(buf, `<script type="text/javascript" src="%s"></script>`, m.GetURL("main.js"))
		return buf.String()
	}

	m.Mount(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, renderTemplate())
	})

}
