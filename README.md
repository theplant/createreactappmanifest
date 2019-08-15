

createreactappmanifest integrate Create React App with Go web application




* [Type Config](#type-config)
* [Type Data](#type-data)
* [Type MData](#type-mdata)
* [Type Manifest](#type-manifest)
  * [New](#manifest-new)
  * [Get CS SU RLs](#manifest-get-cs-su-rls)
  * [Get JS UR Ls](#manifest-get-js-ur-ls)
  * [Mount](#manifest-mount)






## Type: Config
``` go
type Config struct {
    // the PUBLIC_URL environment variable when do yarn build in react app
    PublicURL   string
    ManifestDir string

    // sometime you want to excludes the files in react app public folder.
    MountExcludeForPublic string
    IsDev                 bool
    DevBundleURL          string
}
```









## Type: Data
``` go
type Data struct {
    Files map[string]string
}
```









## Type: MData
``` go
type MData struct {
    JS  []string
    CSS []string
}
```









## Type: Manifest
``` go
type Manifest struct {
    // contains filtered or unexported fields
}
```






### Manifest: New
``` go
func New(cfg *Config) (m *Manifest, err error)
```

Mount expose assets with ServeMux, and GetURL get correct assets path for you.
```go
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
```



### Manifest: Get CS SU RLs
``` go
func (m *Manifest) GetCSSURLs() (urls []string)
```
GetCSSURLs get all the css urls




### Manifest: Get JS UR Ls
``` go
func (m *Manifest) GetJSURLs() (urls []string)
```
GetJSURLs get all the js urls




### Manifest: Mount
``` go
func (m *Manifest) Mount(mux *http.ServeMux)
```
Mount automatically mounts Create React App build directory into Go ServeMux





